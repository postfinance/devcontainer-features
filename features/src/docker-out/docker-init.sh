#!/usr/bin/env bash
set -e

# Wrapper function to only use sudo if not already root
sudoIf() {
    if [ "$(id -u)" -ne 0 ]; then
        sudo "$@"
    else
        "$@"
    fi
}

# Source socket (socket on the host)
SOURCE_SOCKET="${SOURCE_SOCKET:-"/var/run/docker-host.sock"}"
# Target socket (socket in the container)
TARGET_SOCKET="${TARGET_SOCKET:-"/var/run/docker.sock"}"

# Determine the correct user
USERNAME="${USERNAME:-"${_REMOTE_USER:-"auto"}"}"
if [ "${USERNAME}" = "auto" ]; then
  USERNAME=$(awk -v val=1000 -F ":" '$3==val{print $1}' /etc/passwd)
fi

# Check if the docker group exists
if ! getent group docker > /dev/null 2>&1; then
  # Create the docker group
  echo "Adding missing docker group"
  sudoIf groupadd --system docker
fi

# Add the user to the docker group
sudoIf usermod -a -G docker $USERNAME

# By default, make the source and target sockets the same
if [ "${SOURCE_SOCKET}" != "${TARGET_SOCKET}" ]; then
    sudoIf touch "${SOURCE_SOCKET}"
    sudoIf ln -s "${SOURCE_SOCKET}" "${TARGET_SOCKET}"
fi

# Get the id of the docker group
DOCKER_GID=$(getent group docker | cut -d: -f3)

# Find the GID of the source socket
SOCKET_GID=$(stat -c '%g' ${SOURCE_SOCKET})

# If the group ids don't match, we need to adjust
if [ "$DOCKER_GID" = "$SOCKET_GID" ]; then
  echo "Docker group GID matches socket GID ($DOCKER_GID), no changes needed."
else
  # Find an existing group with the same GID as the source socket
  EXISTING_GROUP=$(getent group $SOCKET_GID | cut -d: -f1)
  # If no group is found, just adjust the docker group to match the socket GID
  if [ -z "$EXISTING_GROUP" ]; then
    echo "Adjusting docker group GID ($DOCKER_GID) to match socket GID ($SOCKET_GID)."
    sudoIf groupmod -g $SOCKET_GID docker
  else
    # Use socat
    echo "Using socat to bridge socket from GID $SOCKET_GID to docker GID $DOCKER_GID."
    sudoIf rm -rf ${TARGET_SOCKET}
    # Check if socat is installed
    if command -v socat > /dev/null 2>&1; then
      echo "socat is already installed."
    else
      echo "socat is not installed. Installing socat."
      if command -v apt-get > /dev/null 2>&1; then
        # Apt-based
        sudoIf apt-get update && sudoIf apt-get install -y socat
      elif command -v apk > /dev/null 2>&1; then
        # Apk-based
        sudoIf apk add socat
      else
        echo "Error: socat is required but could not be installed. Please install socat manually."
        exit 1
      fi
    fi
    # Create the socat bridge
    sudoIf socat UNIX-CONNECT:${SOURCE_SOCKET} UNIX-LISTEN:${TARGET_SOCKET},fork,user=${USERNAME},group=docker,mode=660
  fi
fi

# Execute whatever commands were passed in (if any). This allows us
# to set this script to ENTRYPOINT while still executing the default CMD.
set +e
exec "$@"

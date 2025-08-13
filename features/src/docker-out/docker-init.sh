#!/usr/bin/env bash
set -e

# Find the GID of the socket
SOCKET_GID=$(stat -c '%g' /var/run/docker.sock)

# Find an existing group with the same GID
# This is needed eg. for alpine which already has GID 999 (ping)
EXISTING_GROUP=$(getent group $SOCKET_GID | cut -d: -f1)
if [ -n "$EXISTING_GROUP" ]; then
  # Delete this group
  sudo groupdel $EXISTING_GROUP
fi

# Determine the correct user
USERNAME="${USERNAME:-"${_REMOTE_USER:-"auto"}"}"
if [ "${USERNAME}" = "auto" ]; then
  USERNAME=$(awk -v val=1000 -F ":" '$3==val{print $1}' /etc/passwd)
fi

# Create a matching docker group in the container and add the current user to it
sudo groupadd -g $SOCKET_GID docker
sudo usermod -a -G docker $USERNAME

# Execute whatever commands were passed in (if any). This allows us
# to set this script to ENTRYPOINT while still executing the default CMD.
set +e
exec "$@"

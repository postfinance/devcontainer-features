#!/usr/bin/env bash
set -e

# Create a matching docker group in the container and add the current user to it
SOCKET_GID=$(stat -c '%g' /var/run/docker.sock)
sudo groupadd -g $SOCKET_GID docker
sudo usermod -a -G docker $(whoami)

# Execute whatever commands were passed in (if any). This allows us
# to set this script to ENTRYPOINT while still executing the default CMD.
set +e
exec "$@"

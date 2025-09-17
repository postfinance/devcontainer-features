## Notes

### Cypress Binary Cache

By default, Cypress installs binaries somewhere in the user home. This directory is not a volume and therefore is cleared when the container is rebuild which means, Cypress needs to be reinstalled in each new container.

To prevent this, the Cypress cache folder can be set to a folder from the workspace which is a volume (basically the cloned repository) and with that, it does not need to be reinstalled each time the container is rebuilt.

To do this, add the following environment variable to your container:
```
"containerEnv": {
  "CYPRESS_CACHE_FOLDER": "${containerWorkspaceFolder}/.cypress_cache"
}
```

Also don't forget to add `.cypress_cache` to your `.gitignore` file.

### x11 Socket

The Cypress UI in the dev-container is displayed via x11. If the corresponding socket is not correcly forwarded into the container, this can lead to heavy performance loss due to higher CPU usage.

To correctly forward the x11 socket into the container, make sure to add the following variables to the `runArgs`:
```
"runArgs": [
  "-e", "DISPLAY=${localEnv:DISPLAY}",
  "-v", "/tmp/.X11-unix:/tmp/.X11-unix"
]
```

### System Compatibility

Debian, Ubuntu

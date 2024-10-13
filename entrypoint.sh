#!/bin/bash

# Start the SSH server if needed
#/usr/sbin/sshd -D &

# If additional arguments were passed to the container, execute them
if [ $# -gt 0 ]; then
    # Execute the provided command or script
    exec "$@"
else
    # Start an interactive shell if no arguments were provided
    exec /bin/bash
fi

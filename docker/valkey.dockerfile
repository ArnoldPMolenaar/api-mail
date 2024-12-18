# Use the official Valkey base image.
FROM valkey/valkey:8-alpine

# Copy script to execute linux commands for the container.
COPY docker/scripts/valkey.sh /usr/local/etc/valkey/valkey.sh

# Make the file executable.
RUN chmod +x /usr/local/etc/valkey/valkey.sh

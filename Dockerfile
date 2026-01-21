# syntax=docker/dockerfile:1.4

FROM ubuntu:latest

ARG VERSION=1.1.4
ARG RELEASES_URL="https://github.com/ErdemOzgen/blackdagger/releases"

EXPOSE 8080
EXPOSE 8090

# Avoid tzdata prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages
RUN apt-get update -y && \
    apt-get install -y wget curl bash sudo tzdata git && \
    echo "Etc/UTC" > /etc/timezone && \
    ln -fs /usr/share/zoneinfo/Etc/UTC /etc/localtime && \
    dpkg-reconfigure --frontend noninteractive tzdata

# Download, make executable, and run the blackdagger installer script
RUN curl -L https://raw.githubusercontent.com/ErdemOzgen/blackdagger/main/scripts/blackdagger-installer.sh -o blackdagger-installer.sh && \
    chmod +x blackdagger-installer.sh && \
    sudo bash blackdagger-installer.sh && \
    rm blackdagger-installer.sh

# Environment variables for the application
ENV BLACKDAGGER_HOST=0.0.0.0
ENV BLACKDAGGER_PORT=8080

# Copy the start script and set permissions
COPY ./startservices.sh /usr/local/bin/startservices.sh
RUN chmod +x /usr/local/bin/startservices.sh

CMD ["sh", "/usr/local/bin/startservices.sh"]

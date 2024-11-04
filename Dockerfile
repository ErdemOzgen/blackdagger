# syntax=docker/dockerfile:1.4
FROM --platform=$BUILDPLATFORM ubuntu:latest

ARG TARGETARCH
ARG VERSION=1.0.5
ARG RELEASES_URL="https://github.com/ErdemOzgen/blackdagger/releases"
ARG TARGET_FILE="blackdagger_${VERSION}_linux_${TARGETARCH}.tar.gz"

ARG USER="blackdagger"
ARG USER_UID=1000
ARG USER_GID=$USER_UID

EXPOSE 8080
EXPOSE 8090

# Avoid tzdata prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages and setup user before switching
RUN echo "Etc/UTC" > /etc/timezone && \
    ln -fs /usr/share/zoneinfo/Etc/UTC /etc/localtime && \
    apt-get update -y && \
    apt-get install -y sudo tzdata wget curl bash git && \
    dpkg-reconfigure --frontend noninteractive tzdata && \
    groupadd -g ${USER_GID} ${USER} && \
    useradd -m -d /home/${USER} -u ${USER_UID} -g ${USER_GID} -s /bin/bash ${USER} && \
    echo "${USER} ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/${USER} && \
    chmod 0440 /etc/sudoers.d/${USER}

# Switch back to root if needed for installing software or performing tasks that require root privileges
USER root

WORKDIR /home/${USER}

# Download and install the application
# RUN wget ${RELEASES_URL}/download/v${VERSION}/${TARGET_FILE} && \
#     tar -xf ${TARGET_FILE} && rm *.tar.gz && \
#     mv blackdagger /usr/local/bin/ && \
#     mkdir .blackdagger

RUN curl -L https://raw.githubusercontent.com/ErdemOzgen/blackdagger/main/scripts/blackdagger-installer.sh | sudo bash

ENV BLACKDAGGER_HOST=0.0.0.0
ENV BLACKDAGGER_PORT=8080
COPY ./startservices.sh /home/${USER}/startservices.sh
#CMD ["blackdagger", "server"]

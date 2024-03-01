# syntax=docker/dockerfile:1.4
FROM --platform=$BUILDPLATFORM ubuntu:22.04

ARG TARGETARCH
ARG VERSION=1.0.1
ARG RELEASES_URL="https://github.com/ErdemOzgen/blackdagger/releases"
ARG TARGET_FILE="dagu_${VERSION}_linux_${TARGETARCH}.tar.gz"

ARG USER="blackdagger"
ARG USER_UID=1000
ARG USER_GID=$USER_UID

EXPOSE 8080 8090

RUN <<EOF
    # User and permissions setup
    apt-get update
    apt-get install -y sudo tzdata wget
    groupadd -g ${USER_GID} ${USER} || true
    useradd -m ${USER} -u ${USER_UID} -g ${USER_GID} -s /bin/bash
    echo ${USER} ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/${USER}
    chmod 0440 /etc/sudoers.d/${USER}

    # Download and install gotty
    wget https://github.com/some-location/gotty/releases/download/v1.0.0/gotty_linux_amd64 -O /usr/local/bin/gotty
    chmod +x /usr/local/bin/gotty
EOF

USER blackdagger
WORKDIR /home/blackdagger
RUN <<EOF
    export TARGET_FILE="blackdagger_${VERSION}_Linux_${TARGETARCH}.tar.gz"
    wget ${RELEASES_URL}/download/v${VERSION}/${TARGET_FILE}
    tar -xf ${TARGET_FILE} && rm *.tar.gz
    sudo mv blackdagger /usr/local/bin/
    mkdir .blackdagger
EOF

ENV blackdagger_HOST=0.0.0.0
ENV blackdagger_PORT=8080

# Start blackdagger in the background and gotty in the foreground
CMD bash -c "blackdagger server &" && gotty -p 8090 -w --credential blackcart:blackcart /bin/bash

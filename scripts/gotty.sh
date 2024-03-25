#!/bin/sh

RELEASES_URL="https://github.com/yudai/gotty/releases"
TARGET_VERSION="v1.0.1"

echo "Preparing to download gotty $TARGET_VERSION..."

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    arm*) ARCH="arm" ;;
    aarch64) ARCH="arm64" ;;
    i386) ARCH="386" ;;
    i686) ARCH="386" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

TAR_FILE="gotty_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="$RELEASES_URL/download/$TARGET_VERSION/$TAR_FILE"
TMPDIR=$(mktemp -d)
DOWNLOAD_PATH="${TMPDIR}/${TAR_FILE}"

echo "Downloading $TAR_FILE from $DOWNLOAD_URL..."
curl -sfLo "$DOWNLOAD_PATH" "$DOWNLOAD_URL"

if [ ! -f "$DOWNLOAD_PATH" ]; then
    echo "Failed to download $TAR_FILE. Please check the URL and try again."
    exit 1
fi

echo "Extracting $TAR_FILE..."
tar -xzf "$DOWNLOAD_PATH" -C "$TMPDIR"

if [ ! -f "${TMPDIR}/gotty" ]; then
    echo "Failed to extract. The gotty binary is not found."
    exit 1
fi

# Check if /usr/bin/ exists, if not, create it
if [ ! -d "/usr/bin/" ]; then
    sudo mkdir -p /usr/bin/
fi

# Move gotty to /usr/bin to make it globally accessible
sudo mv "${TMPDIR}/gotty" /usr/bin/
echo "gotty has been downloaded, extracted, and moved to /usr/bin/ successfully."

# Cleanup
rm -rf "$TMPDIR"

# Uncomment to run gotty with passed arguments
# "./gotty" "$@"
#gotty -p 8090 -w --credential blackcart:blackcart /bin/bash
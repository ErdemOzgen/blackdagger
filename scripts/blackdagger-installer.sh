#!/bin/sh

# Ensure the script is running as root
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root. Please use sudo or log in as root."
  exit 1
fi

# Check Git installation
if ! command -v git &>/dev/null; then
    echo "Git is not installed. Please install Git to pull default yamls."
else
    echo "Git is already installed. It will pull default yamls."
fi

# Function to set OS and ARCH variables
set_os_arch() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64) ARCH="arm64" ;;
        armv7*|armhf|armv7l) ARCH="armv7" ;;
        armv6*) ARCH="armv6" ;;
        i386|i686) ARCH="386" ;;
        *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
    esac
}

set_os_arch

# Download and install gotty
GOTTY_VERSION="v1.5.0"
GOTTY_BASE_URL="https://github.com/sorenisanerd/gotty/releases/download/$GOTTY_VERSION"
GOTTY_TAR_FILE="gotty_${GOTTY_VERSION}_${OS}_${ARCH}.tar.gz"

TMPDIR=$(mktemp -d)

# gotty installation
echo "Downloading gotty from $GOTTY_BASE_URL/$GOTTY_TAR_FILE..."
curl -sfLo "$TMPDIR/$GOTTY_TAR_FILE" "$GOTTY_BASE_URL/$GOTTY_TAR_FILE"

if [ ! -f "$TMPDIR/$GOTTY_TAR_FILE" ]; then
    echo "Failed to download $GOTTY_TAR_FILE. Check the URL and retry."
    rm -rf "$TMPDIR"
    exit 1
fi

echo "Extracting gotty..."
tar -xzf "$TMPDIR/$GOTTY_TAR_FILE" -C "$TMPDIR"

if [ ! -f "$TMPDIR/gotty" ]; then
    echo "Extraction failed: gotty binary not found."
    rm -rf "$TMPDIR"
    exit 1
fi

mv "$TMPDIR/gotty" /usr/bin/
echo "gotty installed successfully."

# blackdagger installation
BLACKDAGGER_BASE_URL="https://github.com/erdemozgen/blackdagger/releases"

# Fetch latest blackdagger version
BLACKDAGGER_VERSION=$(curl -sfL -o /dev/null -w %{url_effective} "$BLACKDAGGER_BASE_URL/latest" | rev | cut -f1 -d'/' | rev)

if [ -z "$BLACKDAGGER_VERSION" ]; then
    echo "Unable to fetch latest blackdagger version."
    rm -rf "$TMPDIR"
    exit 1
fi

BLACKDAGGER_VERSION_NUMBER="${BLACKDAGGER_VERSION#v}"
BLACKDAGGER_TAR_FILE="blackdagger_${BLACKDAGGER_VERSION_NUMBER}_${OS}_${ARCH}.tar.gz"

# Download blackdagger
echo "Downloading blackdagger $BLACKDAGGER_VERSION from $BLACKDAGGER_BASE_URL/download/$BLACKDAGGER_VERSION/$BLACKDAGGER_TAR_FILE..."
curl -sfLo "$TMPDIR/$BLACKDAGGER_TAR_FILE" "$BLACKDAGGER_BASE_URL/download/$BLACKDAGGER_VERSION/$BLACKDAGGER_TAR_FILE"

if [ ! -f "$TMPDIR/$BLACKDAGGER_TAR_FILE" ]; then
    echo "Failed to download $BLACKDAGGER_TAR_FILE. Check the URL and retry."
    rm -rf "$TMPDIR"
    exit 1
fi

echo "Extracting blackdagger..."
tar -xf "$TMPDIR/$BLACKDAGGER_TAR_FILE" -C "$TMPDIR"

if [ ! -f "$TMPDIR/blackdagger" ]; then
    echo "Extraction failed: blackdagger binary not found."
    rm -rf "$TMPDIR"
    exit 1
fi

# Move blackdagger binary
rm -rf "/usr/bin/blackdagger"
mv "$TMPDIR/blackdagger" /usr/bin/
chmod +x /usr/bin/blackdagger

echo "blackdagger installed successfully."

# Cleanup
rm -rf "$TMPDIR"

echo "All installations complete!"

# Uncomment to run gotty or blackdagger
# gotty -p 8090 -w --credential blackcart:blackcart /bin/bash
# blackdagger "$@"
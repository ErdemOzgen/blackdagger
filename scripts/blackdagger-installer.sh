#!/bin/sh

# Check if the script is running as root
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root. Please use sudo or log in as root."
  exit 1
fi

if ! command -v git &>/dev/null; then
    echo "Git is not installed."
    echo "Please install git if you want to pull default yamls !!!"
else
    echo "Git is already installed. It will pull default yamls"
fi

RELEASES_URL="https://github.com/yudai/gotty/releases"
GOTTY_TARGET_VERSION="v1.0.1"

echo "Preparing to download gotty $GOTTY_TARGET_VERSION..."

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
DOWNLOAD_URL="$RELEASES_URL/download/$GOTTY_TARGET_VERSION/$TAR_FILE"
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
echo "GOTTY"
# Move gotty to /usr/bin to make it globally accessible
sudo mv "${TMPDIR}/gotty" /usr/bin/
echo "gotty has been downloaded, extracted, and moved to /usr/bin successfully."



# Uncomment to run gotty with passed arguments
# "./gotty" "$@"
#gotty -p 8090 -w --credential blackcart:blackcart /bin/bash

#!/bin/sh


RELEASES_URL="https://github.com/erdemozgen/blackdagger/releases"
FILE_BASENAME="blackdagger"


echo "Downloading the latest binary to the current directory..."

test -z "$VERSION" && VERSION="$(curl -sfL -o /dev/null -w %{url_effective} "$RELEASES_URL/latest" |
		rev |
		cut -f1 -d'/'|
		rev)"

test -z "$VERSION" && {
	echo "Unable to get blackdagger version." >&2
	exit 1
}

if [ "$( uname -m )" = "x86_64" ]
then
	ARCHITECTURE="amd64"
else
	ARCHITECTURE="$( uname -m )"
fi

test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"
export TAR_FILE="${TMPDIR}${FILE_BASENAME}_$(uname -s)_$ARCHITECTURE.tar.gz"

(
	cd "$TMPDIR"
	echo "Downloading blackdagger $VERSION..."
	curl -sfLo "$TAR_FILE" "$RELEASES_URL/download/$VERSION/${FILE_BASENAME}_${VERSION:1}_$(uname -s)_$ARCHITECTURE.tar.gz"
)

tar -xf "$TAR_FILE" -C "$TMPDIR"
cp "${TMPDIR}/blackdagger" ./


sudo mv "./blackdagger" /usr/bin/
echo "blackdagger has been downloaded, extracted, and moved to /usr/bin/ successfully."
# Cleanup
rm -rf "$TMPDIR"
"blackdagger" "$@"

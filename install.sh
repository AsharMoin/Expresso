#!/bin/bash

# Repository information
REPOOWNER="AsharMoin"
REPONAME="Expresso"
VERSION="v1.0.0"  # Hardcoded version for first release

# Determine OS type
KERNEL=$(uname -s 2>/dev/null || /usr/bin/uname -s)
case ${KERNEL} in
    "Linux"|"linux")
        KERNEL="linux"
        ;;
    "Darwin"|"darwin")
        KERNEL="darwin"
        ;;
    *)
        echo "OS '${KERNEL}' not supported" >&2
        exit 1
        ;;
esac

# Determine architecture
MACHINE=$(uname -m 2>/dev/null || /usr/bin/uname -m)
case ${MACHINE} in
    arm|armv7*)
        MACHINE="arm"
        ;;
    aarch64*|armv8*|arm64)
        MACHINE="arm64"
        ;;
    i[36]86)
        MACHINE="386"
        if [ "darwin" = "${KERNEL}" ]; then
            echo "Your architecture (${MACHINE}) is not supported on macOS" >&2
            exit 1
        fi
        ;;
    x86_64)
        MACHINE="amd64"
        ;;
    *)
        echo "Your architecture (${MACHINE}) is not currently supported" >&2
        exit 1
        ;;
esac

# Binary and installation directory settings
BINNAME="expresso"
BINDIR="${BINDIR:-/usr/local/bin}"
CONFIG_DIR="${HOME}/.config/expresso"

# Download URL for the release - make sure this matches your file naming convention
URL="https://github.com/$REPOOWNER/$REPONAME/releases/download/${VERSION}/expresso_${VERSION}_${KERNEL}_${MACHINE}.tar.gz"

echo "Installing Expresso version $VERSION..."
echo "Downloading from $URL"
echo

# Download the release archive
curl -q --fail --location --progress-bar --output "expresso_${KERNEL}_${MACHINE}.tar.gz" "$URL"

# Extract the archive
tar xzf "expresso_${KERNEL}_${MACHINE}.tar.gz"

# Make binary executable
chmod +x $BINNAME

# Create config directory if it doesn't exist
mkdir -p $CONFIG_DIR

# Install the binary
sudo mv $BINNAME $BINDIR/$BINNAME

# Clean up
rm "expresso_${KERNEL}_${MACHINE}.tar.gz"

echo
echo "Installation of Expresso version $VERSION complete!"
echo "Run 'expresso' to start using it."
echo "Note: You'll need to configure your OpenAI API key on first run."
#!/bin/bash

REPOOWNER="AsharMoin"
REPONAME="Expresso"
RELEASETAG=$(curl -s "https://api.github.com/repos/$REPOOWNER/$REPONAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

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

BINNAME="expresso"
TARBALL_BINNAME="Expresso"
BINDIR="${BINDIR:-/usr/local/bin}"
TARBALL="Expresso_${RELEASETAG}_${KERNEL}_${MACHINE}.tar.gz"
URL="https://github.com/$REPOOWNER/$REPONAME/releases/download/${RELEASETAG}/${TARBALL}"

echo "Installing ${BINNAME} version ${RELEASETAG}..."
echo "Downloading from $URL"
echo

curl --fail --location --progress-bar --output "${TARBALL}" "$URL"
tar xzf "${TARBALL}"
chmod +x "${TARBALL_BINNAME}"
sudo mv "${TARBALL_BINNAME}" "$BINDIR/$BINNAME"
rm "${TARBALL}"

echo
echo "Installation of expresso version ${RELEASETAG} complete!"
echo "Run 'expresso' to start using it."
echo "Note: You'll need to configure your OpenAI API key on first run."
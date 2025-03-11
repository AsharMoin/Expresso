#!/bin/bash

# Binary and configuration settings
BINNAME="${BINNAME:-expresso}"
BINDIR="${BINDIR:-/usr/local/bin}"
CONFIG_DIR="${HOME}/.config/expresso"

echo "Uninstalling Expresso..."
echo

# Remove binary
if [ -f "$BINDIR/$BINNAME" ]; then
    sudo rm "$BINDIR/$BINNAME"
    echo "Removed binary from $BINDIR/$BINNAME"
else
    echo "Binary not found at $BINDIR/$BINNAME"
fi

# Remove configuration directory
if [ -d "$CONFIG_DIR" ]; then
    rm -rf "$CONFIG_DIR"
    echo "Removed configuration directory at $CONFIG_DIR"
else
    echo "Configuration directory not found at $CONFIG_DIR"
fi

echo
echo "Uninstallation of Expresso complete!"
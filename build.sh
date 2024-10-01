#!/bin/bash

# Name of the binary
BINARY_NAME=vc-analyze

# Detect OS type
if [[ "$OSTYPE" == "linux-gnu"* || "$OSTYPE" == "darwin"* ]]; then
    # Linux or macOS
    INSTALL_DIR="/usr/local/bin"
    echo "Detected Unix-based OS (Linux/macOS)"
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    # Windows with Git Bash (msys)
    INSTALL_DIR="/c/Program Files/VC-Analyze"
    echo "Detected Windows with Git Bash"
else
    echo "Unsupported OS: $OSTYPE"
    exit 1
fi

echo "Building $BINARY_NAME..."
GO111MODULE=on go build -o $BINARY_NAME ./cmd/vc-analyze

# Create installation directory if it doesn't exist (on Windows)
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    mkdir -p "$INSTALL_DIR"
fi

echo "Installing $BINARY_NAME to $INSTALL_DIR..."
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo mv "$BINARY_NAME" "$INSTALL_DIR"
fi

echo "Installation completed successfully!"

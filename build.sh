#!/bin/bash

# Name of the binary
BINARY_NAME=vc-analyze
# Install path
INSTALL_DIR=/usr/local/bin

echo "Building $BINARY_NAME..."
GO111MODULE=on go build -o $BINARY_NAME ./cmd/vc-analyze

echo "Installing $BINARY_NAME to $INSTALL_DIR..."
sudo mv $BINARY_NAME $INSTALL_DIR

echo "Done!"

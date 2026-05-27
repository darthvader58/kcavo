#!/bin/bash

# kubectl-cost Installation Script
# This script helps install all dependencies and build the plugin

set -e  # Exit on error

echo "kcavo Installation Script"
echo "===================================="
echo ""

# Check Go installation
echo "Checking prerequisites..."
if ! command -v go &> /dev/null; then
    echo "Go is not installed!"
    echo "   Please install Go 1.21 or higher from https://go.dev/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo "Found Go: $GO_VERSION"
echo ""

# Check kubectl installation (optional but recommended)
if command -v kubectl &> /dev/null; then
    echo "Found kubectl: $(kubectl version --client --short 2>/dev/null || kubectl version --client)"
else
    echo "kubectl not found - you'll need it to use the plugin"
fi
echo ""

# Step 1: Clean any existing dependencies
echo "Cleaning existing dependencies..."
rm -f go.sum
rm -rf vendor/
echo "Cleaned"
echo ""

# Step 2: Download dependencies
echo "Downloading dependencies (this may take a minute)..."
go mod download
echo "Dependencies downloaded"
echo ""

# Step 3: Tidy up go.mod and generate go.sum
echo "Tidying up dependencies..."
go mod tidy
echo "Dependencies tidied"
echo ""

# Step 4: Verify dependencies
echo "Verifying dependencies..."
go mod verify
echo "Dependencies verified"
echo ""

# Step 5: Build the binary
echo "Building kubectl-cost..."
mkdir -p bin
go build -v -o bin/kubectl-cost .
echo "Build complete: bin/kubectl-cost"
echo ""

# Step 6: Install (optional)
read -p "Install to ~/.local/bin? (y/n) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    mkdir -p ~/.local/bin
    cp bin/kubectl-cost ~/.local/bin/
    chmod +x ~/.local/bin/kubectl-cost
    echo "Installed to ~/.local/bin/kubectl-cost"
    echo ""
    
    # Check if ~/.local/bin is in PATH
    if [[ ":$PATH:" == *":$HOME/.local/bin:"* ]]; then
        echo "~/.local/bin is in your PATH"
    else
        echo "Add ~/.local/bin to your PATH:"
        echo "   echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc"
        echo "   source ~/.bashrc"
    fi
fi
echo ""

# Step 7: Test the installation
echo "Testing installation..."
if [ -f bin/kubectl-cost ]; then
    ./bin/kubectl-cost --version
    echo "kubectl-cost is working!"
else
    echo "Build failed"
    exit 1
fi
echo ""

echo "Installation Complete!"
echo ""

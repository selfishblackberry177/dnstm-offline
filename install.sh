#!/bin/bash
set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="dnstm"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

error() {
    echo -e "${RED}Error: $1${NC}" >&2
    exit 1
}

success() {
    echo -e "${GREEN}$1${NC}"
}

warn() {
    echo -e "${YELLOW}$1${NC}"
}

usage() {
    echo "Usage: $0 [path-to-dnstm-binary]"
    echo ""
    echo "Offline installer for dnstm. No internet access required."
    echo ""
    echo "Arguments:"
    echo "  path-to-dnstm-binary   Path to the dnstm binary (default: ./dnstm)"
    echo ""
    echo "Example:"
    echo "  scp dnstm-linux-amd64 root@server:~/dnstm"
    echo "  ssh root@server 'bash install.sh ./dnstm'"
    exit 0
}

# Parse arguments
BINARY_PATH="${1:-./dnstm}"

if [ "$BINARY_PATH" = "--help" ] || [ "$BINARY_PATH" = "-h" ]; then
    usage
fi

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    error "Please run as root (sudo)"
fi

# Verify the binary exists
if [ ! -f "$BINARY_PATH" ]; then
    error "Binary not found at: $BINARY_PATH"
fi

# Install
echo "Installing dnstm..."
mkdir -p "$INSTALL_DIR"
cp "$BINARY_PATH" "${INSTALL_DIR}/${BINARY_NAME}"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
success "Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"
echo ""

# Run install
echo "Running dnstm install..."
"${INSTALL_DIR}/${BINARY_NAME}" install
echo ""

# Show SCP instructions for transport binaries
echo "==========================================="
echo "  Transport binaries"
echo "==========================================="
echo ""
echo "If any transport binaries are missing, copy them to ${INSTALL_DIR}/ on this server."
echo ""
echo "Example (from your local machine):"
echo "  scp dnstt-server slipstream-server ssserver microsocks sshtun-user root@server:${INSTALL_DIR}/"
echo ""
echo "Then re-run: ${BINARY_NAME} install --force"
echo ""

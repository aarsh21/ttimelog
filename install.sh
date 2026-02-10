#!/bin/sh
# ttimelog installer
# Usage: curl -fsSL https://raw.githubusercontent.com/aarsh21/ttimelog/main/install.sh | sh
#
# Environment variables:
#   TTIMELOG_INSTALL_DIR  - Override install directory (default: /usr/local/bin or ~/.local/bin)
#   TTIMELOG_VERSION      - Install a specific version (default: latest)

set -e

# --- Logging helpers -----------------------------------------------------------

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
RESET='\033[0m'

info()    { printf "${BLUE}[info]${RESET}    %s\n" "$*"; }
success() { printf "${GREEN}[ok]${RESET}      %s\n" "$*"; }
warn()    { printf "${YELLOW}[warn]${RESET}    %s\n" "$*"; }
error()   { printf "${RED}[error]${RESET}   %s\n" "$*" >&2; }
step()    { printf "\n${BOLD}==> %s${RESET}\n" "$*"; }

# --- Cleanup on failure --------------------------------------------------------

TMPDIR=""
cleanup() {
    if [ -n "$TMPDIR" ] && [ -d "$TMPDIR" ]; then
        info "Cleaning up temporary files: $TMPDIR"
        rm -rf "$TMPDIR"
    fi
}
trap cleanup EXIT

# --- Detect OS and architecture ------------------------------------------------

detect_platform() {
    step "Detecting platform"

    OS="$(uname -s)"
    ARCH="$(uname -m)"

    case "$OS" in
        Linux*)  OS="linux"  ;;
        Darwin*) OS="darwin" ;;
        *)
            error "Unsupported operating system: $OS"
            error "ttimelog supports Linux and macOS only."
            exit 1
            ;;
    esac
    info "Operating system: $OS"

    case "$ARCH" in
        x86_64|amd64)   ARCH="amd64"  ;;
        aarch64|arm64)   ARCH="arm64"  ;;
        *)
            error "Unsupported architecture: $ARCH"
            error "ttimelog supports amd64 and arm64 only."
            exit 1
            ;;
    esac
    info "Architecture: $ARCH"

    success "Platform: ${OS}/${ARCH}"
}

# --- Check for required tools --------------------------------------------------

check_dependencies() {
    step "Checking required tools"

    for cmd in curl tar; do
        if command -v "$cmd" >/dev/null 2>&1; then
            info "Found: $cmd ($(command -v "$cmd"))"
        else
            error "Required tool not found: $cmd"
            error "Please install '$cmd' and try again."
            exit 1
        fi
    done

    success "All required tools are available"
}

# --- Resolve version -----------------------------------------------------------

resolve_version() {
    step "Resolving version"

    if [ -n "${TTIMELOG_VERSION:-}" ]; then
        VERSION="$TTIMELOG_VERSION"
        info "Using user-specified version: $VERSION"
    else
        info "Fetching latest release from GitHub..."
        VERSION="$(curl -fsSL -H "Accept: application/vnd.github.v3+json" \
            https://api.github.com/repos/aarsh21/ttimelog/releases/latest 2>/dev/null \
            | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')" || true

        if [ -z "$VERSION" ]; then
            error "Could not determine the latest release version."
            error ""
            error "This can happen if:"
            error "  - No releases have been published yet"
            error "  - GitHub API rate limit was exceeded"
            error ""
            error "Alternatives:"
            error "  1. Set a specific version:  TTIMELOG_VERSION=v0.1.0 sh install.sh"
            error "  2. Install from source:     go install github.com/Rash419/ttimelog/cmd/ttimelog@latest"
            exit 1
        fi
        info "Latest release: $VERSION"
    fi

    # Ensure version starts with 'v'
    case "$VERSION" in
        v*) ;;
        *)  VERSION="v${VERSION}" ;;
    esac

    success "Will install ttimelog $VERSION"
}

# --- Determine install directory -----------------------------------------------

determine_install_dir() {
    step "Determining install location"

    if [ -n "${TTIMELOG_INSTALL_DIR:-}" ]; then
        INSTALL_DIR="$TTIMELOG_INSTALL_DIR"
        info "Using user-specified directory: $INSTALL_DIR"
    elif [ "$(id -u)" = "0" ]; then
        INSTALL_DIR="/usr/local/bin"
        info "Running as root -- installing to $INSTALL_DIR"
    elif [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
        info "Directory /usr/local/bin is writable -- installing there"
    else
        INSTALL_DIR="${HOME}/.local/bin"
        info "/usr/local/bin is not writable -- falling back to $INSTALL_DIR"
    fi

    # Create install dir if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        info "Creating directory: $INSTALL_DIR"
        mkdir -p "$INSTALL_DIR"
    fi

    success "Install directory: $INSTALL_DIR"

    # Warn if not in PATH
    case ":${PATH}:" in
        *":${INSTALL_DIR}:"*) ;;
        *)
            warn "$INSTALL_DIR is not in your PATH"
            warn "Add it by running:"
            warn "  export PATH=\"${INSTALL_DIR}:\$PATH\""
            warn "Or add that line to your shell profile (~/.bashrc, ~/.zshrc, etc.)"
            ;;
    esac
}

# --- Download and install ------------------------------------------------------

download_and_install() {
    step "Downloading ttimelog"

    ARCHIVE_NAME="ttimelog_${OS}_${ARCH}.tar.gz"
    # Strip the leading 'v' for the download URL path
    VERSION_NUM="${VERSION#v}"
    DOWNLOAD_URL="https://github.com/aarsh21/ttimelog/releases/download/${VERSION}/${ARCHIVE_NAME}"

    info "Archive: $ARCHIVE_NAME"
    info "URL: $DOWNLOAD_URL"

    TMPDIR="$(mktemp -d)"
    info "Downloading to temporary directory: $TMPDIR"

    HTTP_CODE=$(curl -fsSL -w "%{http_code}" -o "${TMPDIR}/${ARCHIVE_NAME}" "$DOWNLOAD_URL" 2>/dev/null) || true

    if [ ! -f "${TMPDIR}/${ARCHIVE_NAME}" ] || [ "$HTTP_CODE" = "404" ]; then
        error "Download failed (HTTP $HTTP_CODE)"
        error "URL: $DOWNLOAD_URL"
        error ""
        error "The release asset may not exist. Check:"
        error "  https://github.com/aarsh21/ttimelog/releases"
        exit 1
    fi

    success "Download complete ($(du -h "${TMPDIR}/${ARCHIVE_NAME}" | cut -f1 | xargs))"

    step "Extracting archive"
    info "Extracting ${ARCHIVE_NAME}..."
    tar -xzf "${TMPDIR}/${ARCHIVE_NAME}" -C "$TMPDIR"

    if [ ! -f "${TMPDIR}/ttimelog" ]; then
        error "Expected binary 'ttimelog' not found in archive."
        error "Archive contents:"
        ls -la "$TMPDIR" >&2
        exit 1
    fi
    success "Extracted ttimelog binary"

    step "Installing binary"
    info "Moving ttimelog to ${INSTALL_DIR}/ttimelog"
    mv "${TMPDIR}/ttimelog" "${INSTALL_DIR}/ttimelog"
    chmod +x "${INSTALL_DIR}/ttimelog"
    success "Installed to ${INSTALL_DIR}/ttimelog"
}

# --- Verify installation -------------------------------------------------------

verify_installation() {
    step "Verifying installation"

    if [ -x "${INSTALL_DIR}/ttimelog" ]; then
        info "Binary exists and is executable: ${INSTALL_DIR}/ttimelog"
        info "Size: $(du -h "${INSTALL_DIR}/ttimelog" | cut -f1 | xargs)"
    else
        error "Binary not found or not executable at ${INSTALL_DIR}/ttimelog"
        exit 1
    fi

    # Check if it's reachable from PATH
    if command -v ttimelog >/dev/null 2>&1; then
        WHICH_PATH="$(command -v ttimelog)"
        success "ttimelog is available in PATH: $WHICH_PATH"
    else
        warn "ttimelog is not yet in your PATH"
        warn "Run: export PATH=\"${INSTALL_DIR}:\$PATH\""
    fi
}

# --- Print summary -------------------------------------------------------------

print_summary() {
    printf "\n"
    printf "${GREEN}${BOLD}ttimelog ${VERSION} installed successfully!${RESET}\n"
    printf "\n"
    printf "  ${BOLD}Binary:${RESET}   ${INSTALL_DIR}/ttimelog\n"
    printf "  ${BOLD}Config:${RESET}   ~/.ttimelog/ttimelogrc\n"
    printf "  ${BOLD}Data:${RESET}     ~/.ttimelog/ttimelog.txt\n"
    printf "  ${BOLD}Logs:${RESET}     ~/.ttimelog/ttimelog.log\n"
    printf "\n"
    printf "  Run ${BOLD}ttimelog${RESET} to get started.\n"
    printf "\n"
}

# --- Main ----------------------------------------------------------------------

main() {
    printf "\n"
    printf "${BOLD}ttimelog installer${RESET}\n"
    printf "https://github.com/aarsh21/ttimelog\n"

    detect_platform
    check_dependencies
    resolve_version
    determine_install_dir
    download_and_install
    verify_installation
    print_summary
}

main

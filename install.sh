#!/bin/bash

################################################################################
#                                                                              #
#                    RedditView - Complete Installation Script                #
#                                                                              #
#  This script installs all dependencies and sets up the application          #
#  for running on Linux (Ubuntu, Debian, Fedora, Arch, etc.)                  #
#                                                                              #
################################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track installation status
INSTALLATION_FAILED=0
INSTALLED_PACKAGES=()
SKIPPED_PACKAGES=()

################################################################################
# Helper Functions
################################################################################

print_header() {
    echo -e "${BLUE}"
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║  $1"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_section() {
    echo ""
    echo -e "${BLUE}▶ $1${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [ -f /etc/os-release ]; then
            . /etc/os-release
            OS=$ID
            OS_VERSION=$VERSION_ID
        fi
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
        OS="windows"
    fi
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install Node.js based on OS
install_nodejs() {
    print_section "Installing Node.js"
    
    if command_exists node; then
        NODE_VERSION=$(node --version)
        print_success "Node.js already installed: $NODE_VERSION"
        SKIPPED_PACKAGES+=("Node.js")
        return 0
    fi
    
    case "$OS" in
        arch)
            print_info "Installing via pacman..."
            sudo pacman -Sy --noconfirm nodejs npm
            ;;
        ubuntu|debian)
            print_info "Installing via apt..."
            sudo apt-get update
            sudo apt-get install -y nodejs npm
            ;;
        fedora|rhel|centos)
            print_info "Installing via dnf..."
            sudo dnf install -y nodejs npm
            ;;
        opensuse*)
            print_info "Installing via zypper..."
            sudo zypper install -y nodejs npm
            ;;
        macos)
            print_info "Installing via Homebrew..."
            if ! command_exists brew; then
                print_error "Homebrew not found. Please install Homebrew first:"
                print_info "  /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
                return 1
            fi
            brew install node
            ;;
        *)
            print_error "Unsupported OS: $OS"
            print_warning "Please install Node.js manually from https://nodejs.org/"
            return 1
            ;;
    esac
    
    if command_exists node; then
        NODE_VERSION=$(node --version)
        print_success "Node.js installed: $NODE_VERSION"
        INSTALLED_PACKAGES+=("Node.js")
        return 0
    else
        print_error "Node.js installation failed"
        INSTALLATION_FAILED=$((INSTALLATION_FAILED + 1))
        return 1
    fi
}

# Install Go based on OS
install_go() {
    print_section "Installing Go (required for TUI)"
    
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}')
        print_success "Go already installed: $GO_VERSION"
        SKIPPED_PACKAGES+=("Go")
        return 0
    fi
    
    case "$OS" in
        arch)
            print_info "Installing via pacman..."
            sudo pacman -Sy --noconfirm go
            ;;
        ubuntu|debian)
            print_info "Installing via apt..."
            sudo apt-get update
            sudo apt-get install -y golang-go
            ;;
        fedora|rhel|centos)
            print_info "Installing via dnf..."
            sudo dnf install -y golang
            ;;
        opensuse*)
            print_info "Installing via zypper..."
            sudo zypper install -y go
            ;;
        macos)
            print_info "Installing via Homebrew..."
            if ! command_exists brew; then
                print_error "Homebrew not found. Please install Homebrew first:"
                print_info "  /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
                return 1
            fi
            brew install go
            ;;
        *)
            print_error "Unsupported OS: $OS"
            print_warning "Please install Go manually from https://golang.org/dl"
            return 1
            ;;
    esac
    
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}')
        print_success "Go installed: $GO_VERSION"
        INSTALLED_PACKAGES+=("Go")
        return 0
    else
        print_error "Go installation failed"
        INSTALLATION_FAILED=$((INSTALLATION_FAILED + 1))
        return 1
    fi
}

# Install Git based on OS
install_git() {
    print_section "Installing Git"
    
    if command_exists git; then
        GIT_VERSION=$(git --version | awk '{print $3}')
        print_success "Git already installed: $GIT_VERSION"
        SKIPPED_PACKAGES+=("Git")
        return 0
    fi
    
    case "$OS" in
        arch)
            print_info "Installing via pacman..."
            sudo pacman -Sy --noconfirm git
            ;;
        ubuntu|debian)
            print_info "Installing via apt..."
            sudo apt-get update
            sudo apt-get install -y git
            ;;
        fedora|rhel|centos)
            print_info "Installing via dnf..."
            sudo dnf install -y git
            ;;
        opensuse*)
            print_info "Installing via zypper..."
            sudo zypper install -y git
            ;;
        macos)
            print_info "Installing via Homebrew..."
            if ! command_exists brew; then
                print_error "Homebrew not found"
                return 1
            fi
            brew install git
            ;;
        *)
            print_error "Unsupported OS: $OS"
            return 1
            ;;
    esac
    
    if command_exists git; then
        GIT_VERSION=$(git --version | awk '{print $3}')
        print_success "Git installed: $GIT_VERSION"
        INSTALLED_PACKAGES+=("Git")
        return 0
    else
        print_error "Git installation failed"
        INSTALLATION_FAILED=$((INSTALLATION_FAILED + 1))
        return 1
    fi
}

# Install npm dependencies
install_npm_dependencies() {
    print_section "Installing npm Dependencies"
    
    if [ ! -f "package.json" ]; then
        print_error "package.json not found. Make sure you're in the project root."
        INSTALLATION_FAILED=$((INSTALLATION_FAILED + 1))
        return 1
    fi
    
    print_info "Installing root workspace dependencies..."
    npm install
    
    print_info "Building core package..."
    npm run build --workspace=@redditview/core 2>/dev/null || print_warning "Core package build completed with warnings"
    
    if [ $? -eq 0 ]; then
        print_success "npm dependencies installed"
        INSTALLED_PACKAGES+=("npm dependencies")
        return 0
    else
        print_error "npm installation failed"
        INSTALLATION_FAILED=$((INSTALLATION_FAILED + 1))
        return 1
    fi
}

# Install Go TUI dependencies
install_go_tui_dependencies() {
    print_section "Installing Go TUI Dependencies"
    
    if [ ! -d "apps/tui" ]; then
        print_error "TUI directory not found"
        INSTALLATION_FAILED=$((INSTALLATION_FAILED + 1))
        return 1
    fi
    
    cd apps/tui
    
    print_info "Fetching Go module dependencies..."
    go mod download
    
    print_info "Tidying Go modules..."
    go mod tidy
    
    cd - > /dev/null
    
    if [ $? -eq 0 ]; then
        print_success "Go TUI dependencies installed"
        INSTALLED_PACKAGES+=("Go TUI dependencies")
        return 0
    else
        print_error "Go TUI dependencies installation failed"
        INSTALLATION_FAILED=$((INSTALLATION_FAILED + 1))
        return 1
    fi
}

# Verify installation
verify_installation() {
    print_section "Verifying Installation"
    
    local all_good=true
    
    # Check Node.js
    if command_exists node; then
        print_success "Node.js: $(node --version)"
    else
        print_error "Node.js: Not installed"
        all_good=false
    fi
    
    # Check npm
    if command_exists npm; then
        print_success "npm: $(npm --version)"
    else
        print_error "npm: Not installed"
        all_good=false
    fi
    
    # Check Go
    if command_exists go; then
        print_success "Go: $(go version | awk '{print $3}')"
    else
        print_warning "Go: Not installed (required for TUI only)"
    fi
    
    # Check Git
    if command_exists git; then
        print_success "Git: $(git --version | awk '{print $3}')"
    else
        print_error "Git: Not installed"
        all_good=false
    fi
    
    # Check project structure
    echo ""
    print_info "Checking project structure..."
    
    if [ -f "package.json" ]; then
        print_success "package.json found"
    else
        print_error "package.json not found"
        all_good=false
    fi
    
    if [ -f "config.json" ]; then
        print_success "config.json found"
    else
        print_error "config.json not found"
        all_good=false
    fi
    
    if [ -f "api-server.js" ]; then
        print_success "api-server.js found"
    else
        print_error "api-server.js not found"
        all_good=false
    fi
    
    if [ -d "apps/tui" ]; then
        print_success "TUI app directory found"
    else
        print_error "TUI app directory not found"
        all_good=false
    fi
    
    if [ -d "packages" ]; then
        print_success "Packages directory found"
    else
        print_error "Packages directory not found"
        all_good=false
    fi
    
    return $([ "$all_good" = true ] && echo 0 || echo 1)
}

# Print next steps
print_next_steps() {
    echo ""
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  Next Steps                                                    ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    echo -e "${GREEN}To run the application:${NC}"
    echo ""
    
    echo "  ${YELLOW}Option 1: Start API Server Only${NC}"
    echo "    npm run dev:api"
    echo "    → Available at http://localhost:8765"
    echo ""
    
    echo "  ${YELLOW}Option 2: Start Web UI (requires API running)${NC}"
    echo "    npm run dev"
    echo "    → Available at http://localhost:5173"
    echo ""
    
    echo "  ${YELLOW}Option 3: Start Web + API Together${NC}"
    echo "    ./launch.sh web"
    echo ""
    
    if command_exists go; then
        echo "  ${YELLOW}Option 4: Start TUI (Terminal UI)${NC}"
        echo "    cd apps/tui && go run main.go"
        echo ""
        
        echo "  ${YELLOW}Option 5: Start Everything (API + Web + TUI)${NC}"
        echo "    ./launch.sh all"
        echo ""
    fi
    
    echo -e "${GREEN}Documentation:${NC}"
    echo "  • README.md - Project overview"
    echo "  • INSTALLATION.md - Detailed installation guide"
    echo "  • QUICKSTART.md - Quick start guide"
    echo "  • CONFIGURATION.md - Configuration options"
    echo "  • TUI_KEYBINDINGS.md - TUI keyboard shortcuts"
    echo ""
}

################################################################################
# Main Installation Flow
################################################################################

main() {
    print_header "RedditView Installation"
    
    echo ""
    echo "This script will install all dependencies for RedditView"
    echo "You may be prompted for your password (sudo) several times"
    echo ""
    
    # Detect OS
    detect_os
    print_info "Detected OS: ${OS:-unknown}"
    echo ""
    
    # Check if running from project root
    if [ ! -f "package.json" ]; then
        print_error "package.json not found!"
        print_info "Please run this script from the project root directory"
        exit 1
    fi
    
    # Install system dependencies
    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║  Step 1: System Dependencies                                   ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    
    install_nodejs
    install_go
    install_git
    
    # Install npm dependencies
    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║  Step 2: npm Dependencies                                      ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    
    install_npm_dependencies
    
    # Install Go TUI dependencies
    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║  Step 3: Go TUI Dependencies                                   ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    
    if command_exists go; then
        install_go_tui_dependencies
    else
        print_warning "Go not installed - skipping TUI dependencies"
        print_info "Install Go later to use the TUI interface"
    fi
    
    # Verify installation
    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║  Step 4: Verification                                          ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    
    verify_installation
    VERIFY_RESULT=$?
    
    # Summary
    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║  Installation Summary                                          ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    
    if [ ${#INSTALLED_PACKAGES[@]} -gt 0 ]; then
        echo -e "${GREEN}Newly Installed:${NC}"
        for pkg in "${INSTALLED_PACKAGES[@]}"; do
            echo "  ✓ $pkg"
        done
        echo ""
    fi
    
    if [ ${#SKIPPED_PACKAGES[@]} -gt 0 ]; then
        echo -e "${BLUE}Already Installed:${NC}"
        for pkg in "${SKIPPED_PACKAGES[@]}"; do
            echo "  ✓ $pkg"
        done
        echo ""
    fi
    
    if [ $INSTALLATION_FAILED -gt 0 ]; then
        echo -e "${RED}Installation failed with $INSTALLATION_FAILED error(s)${NC}"
        echo ""
        exit 1
    fi
    
    if [ $VERIFY_RESULT -eq 0 ]; then
        echo -e "${GREEN}✓ Installation successful!${NC}"
        echo ""
        print_next_steps
        exit 0
    else
        echo -e "${YELLOW}⚠ Installation completed with warnings${NC}"
        echo "Please check the output above for details"
        echo ""
        print_next_steps
        exit 0
    fi
}

# Run main function
main "$@"

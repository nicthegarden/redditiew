#!/bin/bash

##############################################################################
# RedditView Systemd Setup Script
# 
# This script helps users install and manage RedditView as a systemd service.
# It supports multiple deployment modes:
# - API Server Only (web-only)
# - API Server + TUI (with tmux)
# - Full Installation with systemd service management
#
# Usage: ./setup.sh [OPTIONS]
#        ./setup.sh --help
##############################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Default values
INSTALL_MODE="both"          # both, api-only, web-only
SERVICE_MANAGER="systemd"    # systemd, manual
INSTALL_PATH=""
USERNAME=$(whoami)
ENABLE_SERVICE=false
START_SERVICE=false
HELP=false
VERBOSE=false

##############################################################################
# Helper Functions
##############################################################################

print_header() {
    echo -e "${CYAN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║         RedditView Systemd Service Setup${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════╝${NC}"
}

print_section() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}▶ $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

log_verbose() {
    if [ "$VERBOSE" = true ]; then
        echo -e "${CYAN}[DEBUG] $1${NC}"
    fi
}

show_help() {
    cat << EOF
${CYAN}RedditView Systemd Service Setup${NC}

${BLUE}USAGE:${NC}
    ./setup.sh [OPTIONS]

${BLUE}OPTIONS:${NC}
    -m, --mode MODE              Installation mode: both, api-only, web-only
                                 Default: both (API + TUI with tmux)
    
    -p, --path PATH              Installation path (where redditiew is installed)
                                 Default: current directory
    
    -u, --user USERNAME          System user to run service as
                                 Default: current user
    
    -e, --enable                 Enable service to start on boot
    
    -s, --start                  Start service immediately after installation
    
    -v, --verbose                Enable verbose output
    
    -h, --help                   Show this help message

${BLUE}MODES:${NC}
    both                         Install both API server and TUI (requires tmux)
    api-only                     Install only API server
    web-only                     Install only web interface

${BLUE}EXAMPLES:${NC}
    # Interactive setup (guided)
    ./setup.sh

    # Setup with custom path and enable on boot
    ./setup.sh --path /opt/redditiew --enable

    # Setup API only, start immediately, enable on boot
    ./setup.sh --mode api-only --start --enable

    # Setup with verbose output
    ./setup.sh --verbose

${BLUE}REQUIREMENTS:${NC}
    - systemd (Linux with systemd support)
    - tmux (for API + TUI mode)
    - sudo/root access (to install systemd services)
    - Node.js and Go (for running the application)

${BLUE}FILES CREATED:${NC}
    ~/.config/systemd/user/redditview-api.service
    ~/.config/systemd/user/redditview-tui.service  (if using both/tui mode)
    ~/.config/systemd/user/redditview-web.service  (if using web-only mode)

${BLUE}COMMANDS AFTER INSTALLATION:${NC}
    # Check service status
    systemctl --user status redditview-api

    # Start/stop services
    systemctl --user start redditview-api
    systemctl --user stop redditview-api

    # View logs
    journalctl --user -u redditview-api -f

    # Access TUI session (if installed)
    tmux attach-session -t redditview

EOF
}

##############################################################################
# Dependency Checks
##############################################################################

check_dependencies() {
    print_section "Checking Dependencies"
    
    local missing=()
    
    # Check systemd
    if ! command -v systemctl &> /dev/null; then
        missing+=("systemd")
    else
        print_success "systemd is installed"
    fi
    
    # Check tmux (only needed for TUI mode)
    if [ "$INSTALL_MODE" = "both" ] || [ "$INSTALL_MODE" = "tui" ]; then
        if ! command -v tmux &> /dev/null; then
            missing+=("tmux")
        else
            print_success "tmux is installed"
        fi
    fi
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        missing+=("Node.js")
    else
        NODE_VERSION=$(node --version)
        print_success "Node.js is installed ($NODE_VERSION)"
    fi
    
    # Check Go (only if TUI binary doesn't exist)
    if [ "$INSTALL_MODE" = "both" ]; then
        if [ ! -f "$INSTALL_PATH/apps/tui/redditview" ]; then
            if ! command -v go &> /dev/null; then
                missing+=("Go")
            else
                GO_VERSION=$(go version | awk '{print $3}')
                print_success "Go is installed ($GO_VERSION)"
            fi
        else
            print_success "TUI binary found (Go not required)"
        fi
    fi
    
    # Check git
    if ! command -v git &> /dev/null; then
        missing+=("git")
    else
        print_success "git is installed"
    fi
    
    if [ ${#missing[@]} -gt 0 ]; then
        print_error "Missing dependencies: ${missing[*]}"
        echo ""
        print_info "Install missing dependencies:"
        echo "  Ubuntu/Debian: sudo apt-get install ${missing[@]}"
        echo "  Fedora/RHEL:   sudo dnf install ${missing[@]}"
        echo "  Arch:          sudo pacman -S ${missing[@]}"
        return 1
    fi
    
    return 0
}

##############################################################################
# Configuration Functions
##############################################################################

configure_install_path() {
    print_section "Configuration: Installation Path"
    
    if [ -z "$INSTALL_PATH" ]; then
        INSTALL_PATH=$(pwd)
        print_info "Using current directory: $INSTALL_PATH"
    fi
    
    if [ ! -f "$INSTALL_PATH/config.json" ]; then
        print_warning "config.json not found in $INSTALL_PATH"
        echo "Is this the correct RedditView installation directory? (y/n)"
        read -r response
        if [ "$response" != "y" ]; then
            print_error "Please specify correct path using -p option"
            return 1
        fi
    else
        print_success "Found config.json"
    fi
    
    log_verbose "Install path: $INSTALL_PATH"
    return 0
}

configure_user() {
    print_section "Configuration: System User"
    
    if ! id "$USERNAME" &>/dev/null; then
        print_error "User '$USERNAME' does not exist"
        echo "Create this user? (y/n)"
        read -r response
        if [ "$response" = "y" ]; then
            if sudo useradd -m -s /bin/bash "$USERNAME"; then
                print_success "Created user '$USERNAME'"
            else
                print_error "Failed to create user"
                return 1
            fi
        else
            return 1
        fi
    else
        print_success "User '$USERNAME' exists"
    fi
    
    # Verify ownership
    if [ ! -O "$INSTALL_PATH" ]; then
        print_warning "Installation path is not owned by '$USERNAME'"
        echo "Change ownership? (y/n)"
        read -r response
        if [ "$response" = "y" ]; then
            if sudo chown -R "$USERNAME:$USERNAME" "$INSTALL_PATH"; then
                print_success "Changed ownership to '$USERNAME'"
            else
                print_error "Failed to change ownership"
                return 1
            fi
        fi
    fi
    
    log_verbose "Service user: $USERNAME"
    return 0
}

configure_mode() {
    print_section "Configuration: Installation Mode"
    
    if [ -z "$INSTALL_MODE" ]; then
        echo "Select installation mode:"
        echo "  1) both       - API Server + TUI with tmux (recommended)"
        echo "  2) api-only   - API Server only"
        echo "  3) web-only   - Web interface only"
        read -p "Enter choice (1-3) [1]: " choice
        choice=${choice:-1}
        
        case $choice in
            1) INSTALL_MODE="both" ;;
            2) INSTALL_MODE="api-only" ;;
            3) INSTALL_MODE="web-only" ;;
            *) print_error "Invalid choice"; return 1 ;;
        esac
    fi
    
    case $INSTALL_MODE in
        both)
            print_success "Mode: API Server + TUI (with tmux)"
            ;;
        api-only)
            print_success "Mode: API Server only"
            ;;
        web-only)
            print_success "Mode: Web interface only"
            ;;
        *)
            print_error "Invalid mode: $INSTALL_MODE"
            return 1
            ;;
    esac
    
    return 0
}

##############################################################################
# Service Installation Functions
##############################################################################

locate_binaries() {
    print_section "Locating Binaries"
    
    # Find Node.js
    NODE_BIN=$(which node)
    print_success "Node.js: $NODE_BIN"
    
    # Find tmux (if needed)
    if [ "$INSTALL_MODE" = "both" ]; then
        TMUX_BIN=$(which tmux)
        if [ -z "$TMUX_BIN" ]; then
            print_error "tmux not found but required for TUI mode"
            return 1
        fi
        print_success "tmux: $TMUX_BIN"
    fi
    
    # Locate TUI binary
    if [ "$INSTALL_MODE" = "both" ]; then
        if [ -f "$INSTALL_PATH/apps/tui/redditview" ]; then
            TUI_BIN="$INSTALL_PATH/apps/tui/redditview"
            print_success "TUI binary: $TUI_BIN"
        else
            print_warning "TUI binary not found at $INSTALL_PATH/apps/tui/redditview"
            echo "Build it now? (y/n)"
            read -r response
            if [ "$response" = "y" ]; then
                build_tui_binary || return 1
            else
                print_error "TUI binary required for 'both' mode"
                return 1
            fi
        fi
    fi
    
    return 0
}

build_tui_binary() {
    print_section "Building TUI Binary"
    
    if [ ! -f "$INSTALL_PATH/apps/tui/main.go" ]; then
        print_error "main.go not found"
        return 1
    fi
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        return 1
    fi
    
    print_info "Building TUI application..."
    cd "$INSTALL_PATH/apps/tui"
    
    if go build -o redditview .; then
        print_success "TUI binary built successfully"
        TUI_BIN="$INSTALL_PATH/apps/tui/redditview"
        return 0
    else
        print_error "Failed to build TUI binary"
        return 1
    fi
}

create_systemd_service() {
    local service_name=$1
    local template_file=$2
    local output_file="$HOME/.config/systemd/user/$service_name"
    
    # Create directory if it doesn't exist
    mkdir -p "$HOME/.config/systemd/user"
    
    # Copy and substitute variables
    log_verbose "Creating $service_name from $template_file"
    
    # Determine which binaries to use
    local node_bin=$(which node)
    local tmux_bin=$(which tmux)
    local tui_bin="$INSTALL_PATH/apps/tui/redditview"
    
    # Create the service file with variable substitution
    sed -e "s|__USERNAME__|$USERNAME|g" \
        -e "s|__INSTALL_PATH__|$INSTALL_PATH|g" \
        -e "s|__NODE_BIN__|$node_bin|g" \
        -e "s|__TMUX_BIN__|$tmux_bin|g" \
        -e "s|__TUI_BIN__|$tui_bin|g" \
        "$template_file" > "$output_file"
    
    chmod 644 "$output_file"
    print_success "Created: $output_file"
    
    return 0
}

install_services() {
    print_section "Installing Systemd Services"
    
    # Ensure template directory exists
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    TEMPLATE_DIR="$SCRIPT_DIR/systemd-templates"
    
    if [ ! -d "$TEMPLATE_DIR" ]; then
        print_error "Template directory not found: $TEMPLATE_DIR"
        return 1
    fi
    
    case $INSTALL_MODE in
        both)
            # Install API service
            create_systemd_service "redditview-api.service" \
                "$TEMPLATE_DIR/redditview-api.service" || return 1
            
            # Install TUI service
            create_systemd_service "redditview-tui.service" \
                "$TEMPLATE_DIR/redditview-tui.service" || return 1
            ;;
        api-only)
            # Install only API service
            create_systemd_service "redditview-api.service" \
                "$TEMPLATE_DIR/redditview-api.service" || return 1
            ;;
        web-only)
            # Install only web service
            create_systemd_service "redditview-web.service" \
                "$TEMPLATE_DIR/redditview-web.service" || return 1
            ;;
    esac
    
    # Reload systemd daemon
    print_info "Reloading systemd daemon..."
    systemctl --user daemon-reload
    print_success "Systemd daemon reloaded"
    
    return 0
}

##############################################################################
# Service Management
##############################################################################

enable_services() {
    print_section "Enabling Services"
    
    case $INSTALL_MODE in
        both)
            print_info "Enabling redditview-api..."
            systemctl --user enable redditview-api.service
            print_success "Enabled redditview-api"
            
            print_info "Enabling redditview-tui..."
            systemctl --user enable redditview-tui.service
            print_success "Enabled redditview-tui"
            
            # Enable lingering for user services (if needed for boot)
            if systemctl --user is-enabled --quiet redditview-api; then
                echo "Enable services to start at boot? (y/n)"
                read -r response
                if [ "$response" = "y" ]; then
                    print_info "Enabling user session to run at boot..."
                    sudo loginctl enable-linger "$USERNAME"
                    print_success "User session will start at boot"
                fi
            fi
            ;;
        api-only)
            print_info "Enabling redditview-api..."
            systemctl --user enable redditview-api.service
            print_success "Enabled redditview-api"
            ;;
        web-only)
            print_info "Enabling redditview-web..."
            systemctl --user enable redditview-web.service
            print_success "Enabled redditview-web"
            ;;
    esac
    
    return 0
}

start_services() {
    print_section "Starting Services"
    
    sleep 1  # Brief pause to ensure daemon is ready
    
    case $INSTALL_MODE in
        both)
            print_info "Starting redditview-api..."
            systemctl --user start redditview-api.service
            print_success "Started redditview-api"
            
            sleep 2  # Wait for API to start
            
            print_info "Starting redditview-tui..."
            systemctl --user start redditview-tui.service
            print_success "Started redditview-tui"
            ;;
        api-only)
            print_info "Starting redditview-api..."
            systemctl --user start redditview-api.service
            print_success "Started redditview-api"
            ;;
        web-only)
            print_info "Starting redditview-web..."
            systemctl --user start redditview-web.service
            print_success "Started redditview-web"
            ;;
    esac
    
    sleep 1
    
    return 0
}

##############################################################################
# Status and Verification
##############################################################################

show_status() {
    print_section "Service Status"
    
    case $INSTALL_MODE in
        both)
            echo ""
            print_info "API Service Status:"
            systemctl --user status redditview-api.service --no-pager || true
            echo ""
            print_info "TUI Service Status:"
            systemctl --user status redditview-tui.service --no-pager || true
            ;;
        api-only)
            echo ""
            print_info "API Service Status:"
            systemctl --user status redditview-api.service --no-pager || true
            ;;
        web-only)
            echo ""
            print_info "Web Service Status:"
            systemctl --user status redditview-web.service --no-pager || true
            ;;
    esac
    
    return 0
}

show_next_steps() {
    print_section "Next Steps"
    
    case $INSTALL_MODE in
        both)
            cat << EOF

${GREEN}✓ RedditView services installed successfully!${NC}

${BLUE}To access your services:${NC}

  ${CYAN}API Server Logs:${NC}
    journalctl --user -u redditview-api -f

  ${CYAN}TUI Logs:${NC}
    journalctl --user -u redditview-tui -f

  ${CYAN}Connect to TUI session:${NC}
    tmux attach-session -t redditview

  ${CYAN}API Server Status:${NC}
    systemctl --user status redditview-api

  ${CYAN}TUI Service Status:${NC}
    systemctl --user status redditview-tui

${BLUE}Common Commands:${NC}

  ${CYAN}Stop services:${NC}
    systemctl --user stop redditview-api redditview-tui

  ${CYAN}Restart services:${NC}
    systemctl --user restart redditview-api redditview-tui

  ${CYAN}Start services:${NC}
    systemctl --user start redditview-api redditview-tui

  ${CYAN}Disable from boot:${NC}
    systemctl --user disable redditview-api redditview-tui

  ${CYAN}View service files:${NC}
    cat ~/.config/systemd/user/redditview-api.service
    cat ~/.config/systemd/user/redditview-tui.service

${BLUE}Web Access:${NC}
  Open http://localhost:3000 in your browser

${BLUE}Documentation:${NC}
  See SYSTEMD_SETUP.md for detailed information

EOF
            ;;
        api-only)
            cat << EOF

${GREEN}✓ RedditView API service installed successfully!${NC}

${BLUE}To manage your service:${NC}

  ${CYAN}View logs:${NC}
    journalctl --user -u redditview-api -f

  ${CYAN}Service status:${NC}
    systemctl --user status redditview-api

  ${CYAN}Stop service:${NC}
    systemctl --user stop redditview-api

  ${CYAN}Restart service:${NC}
    systemctl --user restart redditview-api

${BLUE}Web Access:${NC}
  Open http://localhost:3000 in your browser

${BLUE}Documentation:${NC}
  See SYSTEMD_SETUP.md for detailed information

EOF
            ;;
        web-only)
            cat << EOF

${GREEN}✓ RedditView Web service installed successfully!${NC}

${BLUE}To manage your service:${NC}

  ${CYAN}View logs:${NC}
    journalctl --user -u redditview-web -f

  ${CYAN}Service status:${NC}
    systemctl --user status redditview-web

  ${CYAN}Stop service:${NC}
    systemctl --user stop redditview-web

  ${CYAN}Restart service:${NC}
    systemctl --user restart redditview-web

${BLUE}Web Access:${NC}
  Open http://localhost:3000 in your browser

${BLUE}Documentation:${NC}
  See SYSTEMD_SETUP.md for detailed information

EOF
            ;;
    esac
}

##############################################################################
# Main Functions
##############################################################################

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -m|--mode)
                INSTALL_MODE="$2"
                shift 2
                ;;
            -p|--path)
                INSTALL_PATH="$2"
                shift 2
                ;;
            -u|--user)
                USERNAME="$2"
                shift 2
                ;;
            -e|--enable)
                ENABLE_SERVICE=true
                shift
                ;;
            -s|--start)
                START_SERVICE=true
                shift
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -h|--help)
                HELP=true
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                return 1
                ;;
        esac
    done
    return 0
}

main() {
    print_header
    
    # Parse arguments
    if ! parse_arguments "$@"; then
        echo ""
        show_help
        return 1
    fi
    
    # Show help if requested
    if [ "$HELP" = true ]; then
        show_help
        return 0
    fi
    
    # Run setup steps
    if ! check_dependencies; then
        return 1
    fi
    
    if ! configure_install_path; then
        return 1
    fi
    
    if ! configure_user; then
        return 1
    fi
    
    if ! configure_mode; then
        return 1
    fi
    
    if ! locate_binaries; then
        return 1
    fi
    
    if ! install_services; then
        return 1
    fi
    
    # Optional: enable and start services
    if [ "$ENABLE_SERVICE" = true ]; then
        if ! enable_services; then
            return 1
        fi
    else
        echo ""
        echo "Enable services to start on boot? (y/n)"
        read -r response
        if [ "$response" = "y" ]; then
            if ! enable_services; then
                return 1
            fi
        fi
    fi
    
    if [ "$START_SERVICE" = true ]; then
        if ! start_services; then
            return 1
        fi
    else
        echo ""
        echo "Start services now? (y/n)"
        read -r response
        if [ "$response" = "y" ]; then
            if ! start_services; then
                return 1
            fi
        fi
    fi
    
    # Show status and next steps
    show_status
    show_next_steps
    
    print_section "Setup Complete!"
    print_success "RedditView systemd setup finished successfully"
    
    return 0
}

# Run main function with all arguments
main "$@"
exit $?

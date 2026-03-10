#!/bin/bash

################################################################################
# RedditView SystemD Service Setup Script
#
# This script installs RedditView as systemd services (user or system-level).
# It supports:
# - User-level services: ~/.config/systemd/user/ (no sudo required)
# - System-level services: /etc/systemd/system/ (requires sudo)
#
# Service modes:
# - both: API Server + TUI (with tmux)
# - api-only: API Server only
# - web-only: Web interface only
#
# Usage:
#   ./setup-systemctl.sh [OPTIONS]
#   sudo ./setup-systemctl.sh --scope system [OPTIONS]
#
# Examples:
#   # User-level installation (recommended)
#   ./setup-systemctl.sh --scope user --mode both --enable --start
#
#   # System-level installation
#   sudo ./setup-systemctl.sh --scope system --mode both --user redditview --enable
#
################################################################################

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Default values
SCOPE="user"                           # user or system
MODE="both"                            # both, api-only, web-only
INSTALL_PATH="$(pwd)"
SERVICE_USER="$(whoami)"
ENABLE_SERVICES=false
START_SERVICES=false
VERBOSE=false
HELP=false

################################################################################
# Helper Functions
################################################################################

print_header() {
    echo -e "${CYAN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║    RedditView SystemD Service Setup${NC}"
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
${CYAN}RedditView SystemD Service Setup${NC}

${BLUE}USAGE:${NC}
    ./setup-systemctl.sh [OPTIONS]
    sudo ./setup-systemctl.sh --scope system [OPTIONS]

${BLUE}OPTIONS:${NC}
    -s, --scope SCOPE              Installation scope: user, system
                                   Default: user (no sudo needed)
    
    -m, --mode MODE                Service mode: both, api-only, web-only
                                   Default: both (API + TUI)
    
    -p, --path PATH                Installation path (redditiew directory)
                                   Default: current directory
    
    -u, --user USERNAME            System user to run services
                                   Default: current user
                                   (requires sudo for system scope)
    
    -e, --enable                   Enable services at boot
                                   Default: false
    
    --start                        Start services immediately after setup
                                   Default: false
    
    -v, --verbose                  Verbose output for debugging
                                   Default: false
    
    -h, --help                     Show this help message

${BLUE}EXAMPLES:${NC}
    # User-level installation (recommended)
    ./setup-systemctl.sh --scope user --mode both --enable --start

    # System-level with automatic start
    sudo ./setup-systemctl.sh --scope system --mode both --enable --start

    # API-only service
    ./setup-systemctl.sh --scope user --mode api-only --enable

    # Custom installation path
    ./setup-systemctl.sh --scope user --path /opt/redditiew --enable --start

${BLUE}SCOPES:${NC}
    ${CYAN}user${NC}
      • Location: \$HOME/.config/systemd/user/
      • No root required
      • Services run as current user
      • Best for: development, desktops, personal machines
    
    ${CYAN}system${NC}
      • Location: /etc/systemd/system/
      • Requires sudo
      • Services run as specified user
      • Best for: servers, production deployments

${BLUE}MODES:${NC}
    ${CYAN}both${NC} (default)
      • Starts: API Server + TUI (in tmux)
      • Components: redditview-api, redditview-tui
      • Best for: full-featured systems with display
    
    ${CYAN}api-only${NC}
      • Starts: API Server only
      • Components: redditview-api
      • Best for: headless servers, web-only access
    
    ${CYAN}web-only${NC}
      • Starts: Web UI only (Vite dev server)
      • Components: redditview-web
      • Best for: web-only deployments

${BLUE}NOTES:${NC}
    • API runs on port 8765
    • Web UI runs on port 5174
    • TUI requires: tmux, terminal
    • All services are configured to auto-restart on failure
    • Logs available via: journalctl --user -u SERVICE_NAME

${BLUE}QUICK START:${NC}
    1. Run: ./setup-systemctl.sh --enable --start
    2. Check: systemctl --user status redditview-api
    3. Access Web UI: http://localhost:5174
    4. Attach to TUI: tmux attach-session -t redditview

EOF
}

################################################################################
# Parse Arguments
################################################################################

parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -s|--scope)
                SCOPE="$2"
                shift 2
                ;;
            -m|--mode)
                MODE="$2"
                shift 2
                ;;
            -p|--path)
                INSTALL_PATH="$2"
                shift 2
                ;;
            -u|--user)
                SERVICE_USER="$2"
                shift 2
                ;;
            -e|--enable)
                ENABLE_SERVICES=true
                shift
                ;;
            --start)
                START_SERVICES=true
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
                echo "Use -h or --help for usage information"
                exit 1
                ;;
        esac
    done
}

################################################################################
# Validation Functions
################################################################################

validate_inputs() {
    print_section "Validating Configuration"
    
    # Check scope
    if [[ "$SCOPE" != "user" && "$SCOPE" != "system" ]]; then
        print_error "Invalid scope: $SCOPE (must be 'user' or 'system')"
        exit 1
    fi
    print_success "Scope: $SCOPE"
    
    # Check mode
    if [[ "$MODE" != "both" && "$MODE" != "api-only" && "$MODE" != "web-only" ]]; then
        print_error "Invalid mode: $MODE (must be 'both', 'api-only', or 'web-only')"
        exit 1
    fi
    print_success "Mode: $MODE"
    
    # Check if install path exists
    if [ ! -d "$INSTALL_PATH" ]; then
        print_error "Installation path does not exist: $INSTALL_PATH"
        exit 1
    fi
    print_success "Installation path: $INSTALL_PATH"
    
    # Check if redditiew files exist
    if [ ! -f "$INSTALL_PATH/api-server.js" ]; then
        print_error "api-server.js not found in $INSTALL_PATH"
        exit 1
    fi
    print_success "api-server.js found"
    
    if [ ! -f "$INSTALL_PATH/apps/tui/redditview" ]; then
        print_warning "TUI binary not found (will be required for 'both' mode)"
    else
        print_success "TUI binary found"
    fi
    
    # Validate service user
    if [ "$SCOPE" = "system" ] && [ "$(id -u)" -ne 0 ]; then
        print_error "System scope requires root privileges. Use: sudo $0 --scope system ..."
        exit 1
    fi
    
    if [ "$SCOPE" = "system" ]; then
        if ! id "$SERVICE_USER" &>/dev/null; then
            print_error "User '$SERVICE_USER' does not exist"
            exit 1
        fi
        print_success "Service user: $SERVICE_USER"
    else
        print_success "Service user: $SERVICE_USER"
    fi
    
    # Check for required binaries
    if ! command -v node &> /dev/null; then
        print_error "node.js not found in PATH"
        exit 1
    fi
    print_success "node.js found: $(command -v node)"
    
    if ! command -v tmux &> /dev/null; then
        if [ "$MODE" = "both" ]; then
            print_error "tmux not found but required for 'both' mode"
            exit 1
        fi
        print_warning "tmux not found (not needed for this mode)"
    else
        print_success "tmux found: $(command -v tmux)"
    fi
    
    # Check systemd
    if ! command -v systemctl &> /dev/null; then
        print_error "systemctl not found (systemd required)"
        exit 1
    fi
    print_success "systemd found"
}

################################################################################
# Service File Functions
################################################################################

create_service_files() {
    print_section "Creating Service Files"
    
    local SERVICE_DIR
    local SERVICE_CMD
    
    if [ "$SCOPE" = "user" ]; then
        SERVICE_DIR="$HOME/.config/systemd/user"
        SERVICE_CMD="systemctl --user"
        mkdir -p "$SERVICE_DIR"
    else
        SERVICE_DIR="/etc/systemd/system"
        SERVICE_CMD="systemctl"
    fi
    
    log_verbose "Service directory: $SERVICE_DIR"
    
    # Create API service
    if [[ "$MODE" == "both" || "$MODE" == "api-only" ]]; then
        create_api_service "$SERVICE_DIR"
        print_success "Created redditview-api.service"
    fi
    
    # Create TUI service
    if [ "$MODE" = "both" ]; then
        create_tui_service "$SERVICE_DIR"
        print_success "Created redditview-tui.service"
    fi
    
    # Create Web service
    if [ "$MODE" = "web-only" ]; then
        create_web_service "$SERVICE_DIR"
        print_success "Created redditview-web.service"
    fi
    
    # Reload systemd daemon
    if [ "$SCOPE" = "user" ]; then
        systemctl --user daemon-reload
    else
        systemctl daemon-reload
    fi
    print_success "Systemd daemon reloaded"
}

create_api_service() {
    local SERVICE_DIR=$1
    local NODE_BIN=$(command -v node)
    
    cat > "$SERVICE_DIR/redditview-api.service" << EOF
[Unit]
Description=RedditView API Server
Documentation=https://github.com/nicthegarden/redditiew
After=network.target

[Service]
Type=simple
User=$SERVICE_USER
WorkingDirectory=$INSTALL_PATH
ExecStart=$NODE_BIN api-server.js

# Environment variables
Environment="NODE_ENV=production"
Environment="PORT=8765"

# Auto-restart on failure
Restart=on-failure
RestartSec=10
StartLimitInterval=60
StartLimitBurst=3

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=redditview-api

# Process management
KillMode=process
KillSignal=SIGTERM
TimeoutStopSec=10

[Install]
WantedBy=default.target
EOF
}

create_tui_service() {
    local SERVICE_DIR=$1
    local TMUX_BIN=$(command -v tmux)
    local TUI_BIN="$INSTALL_PATH/apps/tui/redditview"
    
    cat > "$SERVICE_DIR/redditview-tui.service" << EOF
[Unit]
Description=RedditView Terminal User Interface (TUI)
Documentation=https://github.com/nicthegarden/redditiew
After=network.target redditview-api.service
Requires=redditview-api.service

[Service]
Type=simple
User=$SERVICE_USER
WorkingDirectory=$INSTALL_PATH

# Phase 1: Create tmux session (idempotent - won't fail if exists)
ExecStartPre=/usr/bin/bash -c '$TMUX_BIN new-session -d -s redditview -x 120 -y 40 || true'

# Phase 2: Wait for session to exist and be ready (up to 3 seconds)
ExecStartPre=/usr/bin/bash -c 'for i in {1..30}; do $TMUX_BIN has-session -t redditview 2>/dev/null && break; sleep 0.1; done'

# Phase 3: Send command to session to start TUI
ExecStartPre=$TMUX_BIN send-keys -t redditview:0 'cd $INSTALL_PATH && $TUI_BIN' Enter

# Phase 4: Keep service running (monitor session existence)
ExecStart=/usr/bin/bash -c 'while $TMUX_BIN has-session -t redditview 2>/dev/null; do sleep 10; done'

# Auto-restart on failure
Restart=on-failure
RestartSec=10
StartLimitInterval=60
StartLimitBurst=3

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=redditview-tui

# Process management
KillMode=mixed
KillSignal=SIGTERM
TimeoutStopSec=10

# Clean up tmux session on stop
ExecStopPost=/usr/bin/bash -c '$TMUX_BIN kill-session -t redditview 2>/dev/null || true'

[Install]
WantedBy=default.target
EOF
}

create_web_service() {
    local SERVICE_DIR=$1
    local NODE_BIN=$(command -v node)
    
    cat > "$SERVICE_DIR/redditview-web.service" << EOF
[Unit]
Description=RedditView Web Interface
Documentation=https://github.com/nicthegarden/redditiew
After=network.target

[Service]
Type=simple
User=$SERVICE_USER
WorkingDirectory=$INSTALL_PATH

# Run Vite development server on port 5174
ExecStart=$NODE_BIN ./node_modules/vite/bin/vite.js --host 0.0.0.0 --port 5174

# Environment variables
Environment="NODE_ENV=production"
Environment="VITE_API_BASE_URL=http://localhost:8765/api"

# Auto-restart on failure
Restart=on-failure
RestartSec=10
StartLimitInterval=60
StartLimitBurst=3

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=redditview-web

# Process management
KillMode=process
KillSignal=SIGTERM
TimeoutStopSec=10

[Install]
WantedBy=default.target
EOF
}

################################################################################
# Service Management Functions
################################################################################

enable_services() {
    if [ "$ENABLE_SERVICES" != true ]; then
        return
    fi
    
    print_section "Enabling Services at Boot"
    
    local CMD
    if [ "$SCOPE" = "user" ]; then
        CMD="systemctl --user"
    else
        CMD="systemctl"
    fi
    
    if [[ "$MODE" == "both" || "$MODE" == "api-only" ]]; then
        $CMD enable redditview-api.service
        print_success "Enabled redditview-api"
    fi
    
    if [ "$MODE" = "both" ]; then
        $CMD enable redditview-tui.service
        print_success "Enabled redditview-tui"
    fi
    
    if [ "$MODE" = "web-only" ]; then
        $CMD enable redditview-web.service
        print_success "Enabled redditview-web"
    fi
}

start_services() {
    if [ "$START_SERVICES" != true ]; then
        return
    fi
    
    print_section "Starting Services"
    
    local CMD
    if [ "$SCOPE" = "user" ]; then
        CMD="systemctl --user"
    else
        CMD="systemctl"
    fi
    
    if [[ "$MODE" == "both" || "$MODE" == "api-only" ]]; then
        $CMD start redditview-api.service
        print_success "Started redditview-api"
        sleep 2
    fi
    
    if [ "$MODE" = "both" ]; then
        $CMD start redditview-tui.service
        print_success "Started redditview-tui"
    fi
    
    if [ "$MODE" = "web-only" ]; then
        $CMD start redditview-web.service
        print_success "Started redditview-web"
    fi
}

verify_installation() {
    print_section "Verifying Installation"
    
    local CMD
    if [ "$SCOPE" = "user" ]; then
        CMD="systemctl --user"
    else
        CMD="systemctl"
    fi
    
    # List services
    echo -e "\n${BLUE}Installed Services:${NC}"
    if [[ "$MODE" == "both" || "$MODE" == "api-only" ]]; then
        $CMD list-unit-files | grep "redditview-api" || true
    fi
    if [ "$MODE" = "both" ]; then
        $CMD list-unit-files | grep "redditview-tui" || true
    fi
    if [ "$MODE" = "web-only" ]; then
        $CMD list-unit-files | grep "redditview-web" || true
    fi
    
    # Check status if started
    if [ "$START_SERVICES" = true ]; then
        echo -e "\n${BLUE}Service Status:${NC}"
        if [[ "$MODE" == "both" || "$MODE" == "api-only" ]]; then
            $CMD status redditview-api.service || true
        fi
        if [ "$MODE" = "both" ]; then
            $CMD status redditview-tui.service || true
        fi
        if [ "$MODE" = "web-only" ]; then
            $CMD status redditview-web.service || true
        fi
    fi
}

print_next_steps() {
    print_section "Next Steps"
    
    local CMD
    if [ "$SCOPE" = "user" ]; then
        CMD="systemctl --user"
    else
        CMD="sudo systemctl"
    fi
    
    echo -e "\n${CYAN}📋 Service Management:${NC}"
    echo "  View status:    $CMD status redditview-api"
    echo "  Start service:  $CMD start redditview-api"
    echo "  Stop service:   $CMD stop redditview-api"
    echo "  Restart:        $CMD restart redditview-api"
    echo "  View logs:      journalctl -u redditview-api -f"
    
    if [ "$MODE" = "both" ]; then
        echo -e "\n${CYAN}🖥️  Accessing TUI:${NC}"
        echo "  Attach tmux:    tmux attach-session -t redditview"
        echo "  Detach (Ctrl+B D)"
        echo "  Kill session:   tmux kill-session -t redditview"
    fi
    
    echo -e "\n${CYAN}🌐 Web Interface:${NC}"
    echo "  Access:         http://localhost:5174"
    echo "  External:       http://<your-ip>:5174"
    
    echo -e "\n${CYAN}🔌 API Server:${NC}"
    echo "  Access:         http://localhost:8765"
    echo "  Health check:   curl http://localhost:8765/health"
    
    echo -e "\n${CYAN}📚 Documentation:${NC}"
    echo "  Full guide:     cat SYSTEMD_SETUP.md"
    echo "  View service:   cat ~/.config/systemd/user/redditview-api.service"
}

################################################################################
# Main Execution
################################################################################

main() {
    parse_arguments "$@"
    
    if [ "$HELP" = true ]; then
        show_help
        exit 0
    fi
    
    print_header
    validate_inputs
    create_service_files
    enable_services
    start_services
    verify_installation
    print_next_steps
    
    echo ""
    print_success "Setup complete!"
}

# Run main function
main "$@"

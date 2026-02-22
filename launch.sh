#!/bin/bash

# RedditView Launcher Script
# Starts API server, Web app, and optional TUI

set -e

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}╭─────────────────────────────────────────────────╮${NC}"
echo -e "${BLUE}│  RedditView - Multi-Platform Launcher           │${NC}"
echo -e "${BLUE}╰─────────────────────────────────────────────────╯${NC}"
echo ""

# Check what to run
if [ "$1" == "all" ]; then
    LAUNCH_WEB=true
    LAUNCH_TUI=true
    LAUNCH_API=true
elif [ "$1" == "tui" ]; then
    LAUNCH_TUI=true
    LAUNCH_API=true
elif [ "$1" == "web" ]; then
    LAUNCH_WEB=true
    LAUNCH_API=true
elif [ "$1" == "api" ]; then
    LAUNCH_API=true
else
    # Default: web app only
    LAUNCH_WEB=true
    LAUNCH_API=true
fi

# Kill all background processes on exit
cleanup() {
    echo ""
    echo -e "${YELLOW}Shutting down...${NC}"
    pkill -P $$ || true
    wait
}
trap cleanup EXIT

# Start API Server
if [ "$LAUNCH_API" = true ]; then
    echo -e "${GREEN}▶ Starting API Server (port 3002)${NC}"
    node api-server.js &
    sleep 2
    
    # Verify API is running
    if curl -s http://localhost:3002/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ API Server running${NC}"
    else
        echo -e "${YELLOW}⚠ API Server may not be responding yet${NC}"
    fi
    echo ""
fi

# Start Web App
if [ "$LAUNCH_WEB" = true ]; then
    echo -e "${GREEN}▶ Starting Web App (port 5173)${NC}"
    echo -e "${YELLOW}  → Open http://localhost:5173${NC}"
    echo ""
    npm run dev &
    sleep 3
fi

# Start TUI
if [ "$LAUNCH_TUI" = true ]; then
    echo -e "${GREEN}▶ Starting TUI (Terminal UI)${NC}"
    echo ""
    cd apps/tui
    go run main.go &
    cd - > /dev/null
fi

# Wait for all background processes
wait

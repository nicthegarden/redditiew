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
    pkill -P $$ 2>/dev/null || true
    # Also kill any processes using our ports
    lsof -ti:3002 2>/dev/null | xargs kill -9 2>/dev/null || true
    lsof -ti:5173 2>/dev/null | xargs kill -9 2>/dev/null || true
    wait 2>/dev/null || true
}
trap cleanup EXIT

# Function to kill process on port
kill_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        echo -e "${YELLOW}  Cleaning up port $port...${NC}"
        lsof -ti :$port | xargs kill -9 2>/dev/null || true
        sleep 1
    fi
}

# Start API Server
if [ "$LAUNCH_API" = true ]; then
    echo -e "${GREEN}▶ Starting API Server (port 3002)${NC}"
    kill_port 3002
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
    kill_port 5173
    npm run dev &
    sleep 3
fi

# Start TUI (in foreground - no & since it needs interactive TTY)
if [ "$LAUNCH_TUI" = true ]; then
    echo -e "${GREEN}▶ Starting TUI (Terminal UI)${NC}"
    echo -e "${YELLOW}Controls: ↑↓/jk=navigate, q=quit${NC}"
    echo ""
    cd apps/tui
    go run main.go
    cd - > /dev/null
fi

# Wait for remaining background processes (API and Web if running)
wait 2>/dev/null || true

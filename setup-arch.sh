#!/bin/bash
# RedditView - Arch Linux Setup Script
# Run this script to install all dependencies

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║                                                                ║"
echo "║      RedditView - Arch Linux Setup Script                     ║"
echo "║                                                                ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# Check if running on Arch Linux
if ! command -v pacman &> /dev/null; then
    echo "❌ Error: pacman not found. This script requires Arch Linux or Arch-based distro."
    exit 1
fi

echo "📦 Step 1: Updating system packages..."
echo "   Command: sudo pacman -Syu"
echo ""
read -p "Press Enter to continue (you may be asked for password)..."
sudo pacman -Syu

echo ""
echo "📦 Step 2: Installing Node.js and npm..."
echo "   Command: sudo pacman -S nodejs npm"
echo ""
if ! command -v node &> /dev/null; then
    sudo pacman -S nodejs npm
    echo "✅ Node.js and npm installed"
else
    echo "✅ Node.js and npm already installed"
    node --version
    npm --version
fi

echo ""
echo "📦 Step 3: Installing Go..."
echo "   Command: sudo pacman -S go"
echo ""
if ! command -v go &> /dev/null; then
    sudo pacman -S go
    echo "✅ Go installed"
else
    echo "✅ Go already installed"
    go version
fi

echo ""
echo "📦 Step 4: Installing build tools (optional)..."
echo "   Command: sudo pacman -S base-devel"
echo ""
read -p "Install base-devel? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    sudo pacman -S base-devel
    echo "✅ Build tools installed"
else
    echo "⏭️  Skipped build tools"
fi

echo ""
echo "✅ All system dependencies installed!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📚 Next steps:"
echo ""
echo "  1. Navigate to project:"
echo "     cd /path/to/redditiew-local"
echo ""
echo "  2. Install npm dependencies:"
echo "     npm install"
echo ""
echo "  3. Build core package:"
echo "     npm run build"
echo ""
echo "  4. Run the app:"
echo ""
echo "     Option A - Web app:"
echo "       npm run dev"
echo "       → Open http://localhost:5173"
echo ""
echo "     Option B - TUI app:"
echo "       Terminal 1: npm run dev:api"
echo "       Terminal 2: cd apps/tui && go run main.go"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "✨ Setup complete!"

#!/bin/bash
# RedditView Web Server Startup Script
set -e

cd /home/edve/2/redditiew

echo "Starting RedditView Web Server..."
echo "Web UI will be available at: http://localhost:5174"
echo "Press Ctrl+C to stop"
echo ""

# Run the server in the foreground
node web-server.js

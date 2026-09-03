#!/bin/sh
# Start the demo server in the foreground (keeps the terminal window open).
# Usage: run-server.sh [port]
DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$DIR/bin/server"

echo "Starting demo server on port ${1:-8080}..."
./server --port "${1:-8080}"

# If we get here the server stopped (error or Ctrl+C) — show it and
# keep the window open so the message is not lost
echo
echo "Server stopped. If it failed, the error is in the lines above."
read -r -p "Press Enter to close this window..." _

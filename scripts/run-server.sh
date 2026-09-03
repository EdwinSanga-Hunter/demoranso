#!/bin/sh
# Start the demo server in the foreground (keeps the terminal window open).
# Usage: run-server.sh [port]
set -e
DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$DIR/bin/server"
echo "Starting demo server on port ${1:-8080}..."
exec ./server --port "${1:-8080}"

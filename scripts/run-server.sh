#!/bin/sh
# Start the demo server in the foreground (keeps the terminal window open).
# Usage: run-server.sh [port]
set -e
DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$DIR/bin/server"
exec ./server --port "${1:-8080}"

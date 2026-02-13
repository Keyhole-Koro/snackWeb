#!/usr/bin/env bash
# ============================================================
#  run_snack_web.sh — Launch SnackWeb (Backend + Frontend)
# ============================================================
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo "🍿 Starting SnackWeb..."
echo ""

# Proactive Cleanup
echo "▸ Cleaning up any existing processes on port 13579..."
fuser -k 13579/tcp 2>/dev/null || true
sleep 1

# ── Backend ──
# ── Backend ──
echo "▸ Starting Backend (Go) on :13579 ..."
cd "$ROOT_DIR/snackWeb/backend"
echo "  Building Go server..."
go build -o server cmd/server/main.go
./server &
BACKEND_PID=$!
echo "  PID: $BACKEND_PID"

echo ""
echo "═══════════════════════════════════════════"
echo "  🍿 SnackWeb is running!"
echo ""
echo "  Backend  : http://localhost:13579"
echo "  API Docs : http://localhost:13579/docs"
echo "═══════════════════════════════════════════"
echo ""
echo "Press Ctrl+C to stop all services."

# Trap Ctrl+C to kill both processes and their children
cleanup() {
  echo ""
  echo "🛑 Shutting down..."
  # Kill children of the stored PIDs (e.g. uvicorn worker, vite process)
  pkill -P $BACKEND_PID 2>/dev/null || true
  # Kill the stored PIDs themselves
  kill $BACKEND_PID 2>/dev/null || true
  wait $BACKEND_PID 2>/dev/null || true
  echo "Done."
}

trap cleanup SIGINT SIGTERM EXIT

wait

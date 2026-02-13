#!/usr/bin/env bash
# ============================================================
#  run_snack_web.sh — Launch SnackWeb (Backend + Frontend)
# ============================================================
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo "🍿 Starting SnackWeb..."
echo ""

# ── Backend ──
echo "▸ Starting Backend (FastAPI) on :8000 ..."
cd "$ROOT_DIR"
uvicorn snackWeb.backend.main:app --host 0.0.0.0 --port 8000 --reload &
BACKEND_PID=$!
echo "  PID: $BACKEND_PID"

# ── Frontend ──
echo "▸ Starting Frontend (Vite) on :5173 ..."
cd "$SCRIPT_DIR/frontend"
npm run dev -- --host &
FRONTEND_PID=$!
echo "  PID: $FRONTEND_PID"

echo ""
echo "═══════════════════════════════════════════"
echo "  🍿 SnackWeb is running!"
echo ""
echo "  Frontend : http://localhost:5173"
echo "  Backend  : http://localhost:8000"
echo "  API Docs : http://localhost:8000/docs"
echo "═══════════════════════════════════════════"
echo ""
echo "Press Ctrl+C to stop all services."

# Trap Ctrl+C to kill both processes
cleanup() {
  echo ""
  echo "🛑 Shutting down..."
  kill $BACKEND_PID 2>/dev/null
  kill $FRONTEND_PID 2>/dev/null
  wait
  echo "Done."
}

trap cleanup SIGINT SIGTERM

wait

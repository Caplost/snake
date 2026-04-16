#!/bin/bash
# E2E Screenshot Script for Snake Game
# Records terminal session to prove game runs correctly
# Usage: ./e2e-screenshot.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SCREENSHOTS_DIR="$PROJECT_DIR/screenshots"

# Create screenshots directory if it doesn't exist
mkdir -p "$SCREENSHOTS_DIR"

# Generate timestamp
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
OUTPUT_FILE="$SCREENSHOTS_DIR/${TIMESTAMP}-snake-gameplay.cast"

echo "Recording snake game session to $OUTPUT_FILE"
echo "Press Ctrl+C to stop early, or wait 10 seconds for auto-stop"

# Check if asciinema is available (preferred)
if command -v asciinema &> /dev/null; then
    asciinema rec "$OUTPUT_FILE" --overwrite --timeout 10 &
    PID=$!
    sleep 10
    kill $PID 2>/dev/null || true
    echo "Recording saved to $OUTPUT_FILE"
else
    # Fallback: use script command
    if command -v script &> /dev/null; then
        script -q -f "$OUTPUT_FILE" &
        PID=$!
        sleep 10
        kill $PID 2>/dev/null || true
        echo "Recording saved to $OUTPUT_FILE"
    else
        echo "Neither asciinema nor script command available"
        echo "Please install asciinema: brew install asciinema"
        exit 1
    fi
fi

echo "Done! Recording: $OUTPUT_FILE"

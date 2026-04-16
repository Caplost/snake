#!/bin/bash
# E2E Screenshot Script for Snake Game
# Records terminal session to prove game runs correctly
# Usage: ./e2e-screenshot.sh
#
# Compatible with both Linux (script -f) and macOS (script -F)

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
echo "Recording for 3 seconds..."

# Determine platform and appropriate script flags
if [[ "$(uname)" == "Darwin" ]]; then
    # macOS: use -F for flush (different semantics than Linux -f)
    SCRIPT_FLAGS="-F"
else
    # Linux: use -f for force flush
    SCRIPT_FLAGS="-f"
fi

# Check if asciinema is available (preferred)
if command -v asciinema &> /dev/null; then
    asciinema rec "$OUTPUT_FILE" --overwrite --timeout 3 &
    PID=$!
    sleep 3
    kill $PID 2>/dev/null || true
    echo "Recording saved to $OUTPUT_FILE"
else
    # Fallback: use script command with platform-specific flags
    if command -v script &> /dev/null; then
        script $SCRIPT_FLAGS -q "$OUTPUT_FILE" &
        PID=$!
        sleep 3
        kill $PID 2>/dev/null || true
        echo "Recording saved to $OUTPUT_FILE"
    else
        echo "Neither asciinema nor script command available"
        echo "Please install asciinema: brew install asciinema"
        exit 1
    fi
fi

echo "Done! Recording: $OUTPUT_FILE"

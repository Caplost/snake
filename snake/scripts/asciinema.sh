#!/bin/bash
# Asciinema-based E2E Recording Script for Snake Game
# Alternative to e2e-screenshot.sh for environments with asciinema installed

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SCREENSHOTS_DIR="$PROJECT_DIR/screenshots"

# Create screenshots directory if it doesn't exist
mkdir -p "$SCREENSHOTS_DIR"

# Generate timestamp
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
OUTPUT_FILE="$SCREENSHOTS_DIR/${TIMESTAMP}-snake-gameplay.cast"

if ! command -v asciinema &> /dev/null; then
    echo "asciinema not found. Install with: brew install asciinema"
    exit 1
fi

echo "Recording snake game session with asciinema..."
echo "Recording to: $OUTPUT_FILE"

# Record for 3 seconds with idle time limit of 1 second
asciinema rec "$OUTPUT_FILE" --overwrite --timeout 3 --idle-time-limit 1

echo "Recording saved to $OUTPUT_FILE"
echo "View with: asciinema play $OUTPUT_FILE"

#!/bin/bash
# E2E Screenshot Script for Snake Game
# Records terminal session to prove game runs correctly
# Usage: ./e2e-screenshot.sh
#
# Fixed v4: Uses foreground timeout instead of background+kill
# macOS script -F only flushes on session end, so kill prevented file content

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SCREENSHOTS_DIR="$PROJECT_DIR/screenshots"

# Create screenshots directory if it doesn't exist
mkdir -p "$SCREENSHOTS_DIR"

# Build snake binary first
SNAKE_BINARY="$PROJECT_DIR/snake"
if [[ ! -f "$SNAKE_BINARY" ]]; then
    echo "Building snake binary..."
    cd "$PROJECT_DIR" && go build -o snake ./cmd/snake
fi

# Generate timestamp
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
OUTPUT_FILE="$SCREENSHOTS_DIR/${TIMESTAMP}-snake-gameplay.cast"

echo "Recording snake game session to $OUTPUT_FILE"
echo "Recording for 5 seconds (or until game exits)..."

# Run snake inside script session with timeout (foreground)
# This ensures .cast file is properly flushed when session ends
if [[ "$(uname)" == "Darwin" ]]; then
    # macOS: use -F for flush
    timeout 5 script -F "$OUTPUT_FILE" "$SNAKE_BINARY" || true
else
    # Linux: use -f for force flush
    timeout 5 script -f "$OUTPUT_FILE" "$SNAKE_BINARY" || true
fi

# Verify file is non-empty
if [[ -f "$OUTPUT_FILE" ]]; then
    FILE_SIZE=$(stat -f%z "$OUTPUT_FILE" 2>/dev/null || stat -c%s "$OUTPUT_FILE" 2>/dev/null)
    if [[ "$FILE_SIZE" -gt 0 ]]; then
        echo "SUCCESS: Recording saved to $OUTPUT_FILE ($FILE_SIZE bytes)"
    else
        echo "WARNING: Recording file is empty (0 bytes)"
        exit 1
    fi
else
    echo "ERROR: Recording file not created"
    exit 1
fi

echo "Done! Recording: $OUTPUT_FILE"

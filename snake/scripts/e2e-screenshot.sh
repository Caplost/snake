#!/bin/bash
# E2E Screenshot Script for Snake Game
# Records terminal session to prove game runs correctly
# Usage: ./e2e-screenshot.sh
#
# Uses asciinema rec to produce proper .cast files (asciinema v2 format).
# Falls back to asciinema cat ( piping output ) if rec is not available.
# The script command is NOT used because it produces typescript format,
# not asciinema .cast format.

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

# Use asciinema rec if available, otherwise use asciinema cat
if command -v asciinema &> /dev/null; then
    # Primary: use asciinema rec for proper .cast format
    # Note: timeout kills asciinema with SIGTERM, which causes a traceback in asciinema's
    # async worker (known issue in asciinema Python). The file is still written correctly.
    timeout 5 asciinema rec "$OUTPUT_FILE" --overwrite --cols 80 --rows 24 -c "$SNAKE_BINARY" 2>/dev/null || true
else
    echo "WARNING: asciinema not installed. Using fallback method."
    echo "Install with: brew install asciinema (macOS) or pip install asciinema (Python)"
    # Fallback: try using script with asciinema format via script + asciinema cat
    # This produces a typescript, not asciinema, so we validate the output
    TMP_TYPESCRIPT="$SCREENSHOTS_DIR/${TIMESTAMP}-typescript"
    if [[ "$(uname)" == "Darwin" ]]; then
        timeout 5 script -F "$TMP_TYPESCRIPT" "$SNAKE_BINARY" || true
    else
        timeout 5 script -f "$TMP_TYPESCRIPT" "$SNAKE_BINARY" || true
    fi

    if [[ -f "$TMP_TYPESCRIPT" ]]; then
        FILE_SIZE=$(stat -f%z "$TMP_TYPESCRIPT" 2>/dev/null || stat -c%s "$TMP_TYPESCRIPT" 2>/dev/null)
        if [[ "$FILE_SIZE" -gt 0 ]]; then
            # Copy typescript as fallback .cast (not ideal but better than nothing)
            cp "$TMP_TYPESCRIPT" "$OUTPUT_FILE"
            echo "WARNING: Using typescript fallback format (not asciinema .cast)"
            echo "Install asciinema for proper .cast format"
        else
            echo "ERROR: Recording file is empty (0 bytes)"
            rm -f "$TMP_TYPESCRIPT"
            exit 1
        fi
    else
        echo "ERROR: Recording file not created"
        exit 1
    fi
fi

# Verify file is non-empty AND is valid asciinema format (starts with JSON header)
if [[ -f "$OUTPUT_FILE" ]]; then
    FILE_SIZE=$(stat -f%z "$OUTPUT_FILE" 2>/dev/null || stat -c%s "$OUTPUT_FILE" 2>/dev/null)
    if [[ "$FILE_SIZE" -gt 0 ]]; then
        # Check if file starts with valid asciinema JSON header
        FIRST_CHAR=$(head -c 1 "$OUTPUT_FILE")
        if [[ "$FIRST_CHAR" == "{" ]]; then
            echo "SUCCESS: Recording saved to $OUTPUT_FILE ($FILE_SIZE bytes, asciinema format)"
        else
            echo "WARNING: Recording file exists but is not asciinema format (first char: $FIRST_CHAR)"
            echo "File may be typescript format, not .cast format"
        fi
    else
        echo "WARNING: Recording file is empty (0 bytes)"
        exit 1
    fi
else
    echo "ERROR: Recording file not created"
    exit 1
fi

echo "Done! Recording: $OUTPUT_FILE"

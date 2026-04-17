#!/bin/bash

# Configuration
SEARCH_DIR="internal/resources"
ERROR_FOUND=0

echo "Checking for @ID in all Swagger-documented endpoints..."

# Find all Go files that contain a @Router tag
FILES=$(grep -rl "@Router" "$SEARCH_DIR")

for FILE in $FILES; do
    # Count occurrences of @Router and @ID in the file
    ROUTER_COUNT=$(grep -c "@Router" "$FILE")
    ID_COUNT=$(grep -c "@ID" "$FILE")

    if [ "$ROUTER_COUNT" -ne "$ID_COUNT" ]; then
        echo "❌ ERROR: $FILE has $ROUTER_COUNT @Router(s) but $ID_COUNT @ID(s)"
        ERROR_FOUND=1
    fi
done

if [ "$ERROR_FOUND" -eq 1 ]; then
    echo "Summary: Missing @ID tags found. Please add @ID to all handlers."
    exit 1
fi

echo "✅ All endpoints have matching @ID tags."
exit 0

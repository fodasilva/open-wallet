#!/bin/bash

# gen-repos.sh - Scanner script to find all .go files with @gen_repo and call gen-repo.sh

echo "Scanning project for @gen_repo tag..."

# Search for all .go files containing the @gen_repo tag, excluding scripts directory
FILES=$(grep -rla "@gen_repo" . --include="*.go" | grep -v "scripts/")

if [ -z "$FILES" ]; then
    echo "No files found with @gen_repo tag."
    exit 0
fi

for file in $FILES; do
    # Call the singular master script for each file
    bash scripts/repository/gen-repo.sh "$file"
done

echo "All repositories updated!"

#!/bin/bash

# gen-select.sh - Specific script for Select operation

METHOD_NAME=$1
STRUCT_NAME=$2
TABLE_NAME=$3
IN_COLS=$4
IN_VALS=$5
COLUMNS=$6
SCAN_FIELDS=$7 # This will now contain newlines from the master script

REPO_NAME=$8
RETURN_TYPE=$9
PAYLOAD=${10}

TEMPLATE="templates/repository/select.txt"

if [ ! -f "$TEMPLATE" ]; then
    echo "Template $TEMPLATE not found"
    exit 1
fi

# Prepare temp files
SCAN_FILE="/tmp/gen_repo_scan_${STRUCT_NAME}_$METHOD_NAME.txt"
# Write SCAN_FIELDS to the temp file
echo "$SCAN_FIELDS" > "$SCAN_FILE"

# Replace placeholders in template
# We replace most things with sed, then insert the scan fields file
cat "$TEMPLATE" | \
    sed "s/{{StructName}}/$STRUCT_NAME/g" | \
    sed "s/{{TableName}}/$TABLE_NAME/g" | \
    sed "s/{{Columns}}/$COLUMNS/g" | \
    sed "s/{{RepoName}}/${REPO_NAME}/g" | \
    sed "s/{{ReturnType}}/${RETURN_TYPE}/g" | \
    sed -e "/{{ScanFields}}/{r $SCAN_FILE" -e "d;}"

rm -f "$SCAN_FILE"

#!/bin/bash

# gen-count.sh - Specific script for Count operation

METHOD_NAME=$1
STRUCT_NAME=$2
TABLE_NAME=$3
# IN_COLS, IN_VALS, COLUMNS, SCAN_FIELDS are not used for Count as it's hardcoded to COUNT(*)
# but we receive them anyway from the master script

REPO_NAME=$8
RETURN_TYPE=$9

TEMPLATE="templates/repository/count.txt"

if [ ! -f "$TEMPLATE" ]; then
    echo "Template $TEMPLATE not found"
    exit 1
fi

# Replace placeholders in template
cat "$TEMPLATE" | \
    sed "s/{{RepoName}}/${REPO_NAME}/g" | \
    sed "s/{{TableName}}/$TABLE_NAME/g" | \
    sed "s/{{ReturnType}}/${RETURN_TYPE}/g"

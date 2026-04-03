#!/bin/bash

# gen-delete.sh - Specific script for Delete operation

METHOD_NAME=$1
STRUCT_NAME=$2
TABLE_NAME=$3

REPO_NAME=$8

TEMPLATE="templates/repository/delete.txt"

if [ ! -f "$TEMPLATE" ]; then
    echo "Template $TEMPLATE not found"
    exit 1
fi

# Replace placeholders in template
cat "$TEMPLATE" | \
    sed "s/{{MethodName}}/$METHOD_NAME/g" | \
    sed "s/{{RepoName}}/${REPO_NAME}/g" | \
    sed "s/{{TableName}}/$TABLE_NAME/g"

#!/bin/bash

# gen-update.sh - Specific script for Update operation

METHOD_NAME=$1
STRUCT_NAME=$2
TABLE_NAME=$3
# IN_COLS, IN_VALS, OUT_COLS, OUT_SCAN are not used directly here
# REPO_NAME, RETURN_TYPE, PAYLOAD, RAW_OUT_COLS are arguments 8-11
REPO_NAME=$8
RETURN_TYPE=$9
PAYLOAD=${10}
IN_UPDATE=${12}

TEMPLATE="templates/repository/update.txt"

if [ ! -f "$TEMPLATE" ]; then
    echo "Template $TEMPLATE not found"
    exit 1
fi

# Prepare temp files for the block of Set calls
UPDATE_FILE="/tmp/gen_repo_update_${STRUCT_NAME}_$METHOD_NAME.txt"
echo "$IN_UPDATE" > "$UPDATE_FILE"

# Replace placeholders in template
cat "$TEMPLATE" | \
    sed "s/{{MethodName}}/$METHOD_NAME/g" | \
    sed "s/{{RepoName}}/${REPO_NAME}/g" | \
    sed "s/{{TableName}}/$TABLE_NAME/g" | \
    sed "s/{{PayloadType}}/$PAYLOAD/g" | \
    sed -e "/{{UpdateSets}}/{r $UPDATE_FILE" -e "d;}"

rm -f "$UPDATE_FILE"

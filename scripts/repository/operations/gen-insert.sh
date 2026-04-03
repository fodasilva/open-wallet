#!/bin/bash

# gen-insert.sh - Specific script for Insert operation

METHOD_NAME=$1
ENTITY=$2
TABLE_NAME=$3
IN_COLS_VALS=$4

REPO_NAME=$8
PAYLOAD=${10}

TEMPLATE="templates/repository/insert.txt"

if [ ! -f "$TEMPLATE" ]; then
    echo "Template $TEMPLATE not found"
    exit 1
fi

IN_COLS_VALS_FILE="/tmp/gen_repo_incolvals_${ENTITY}_$METHOD_NAME.txt"
echo "$IN_COLS_VALS" > "$IN_COLS_VALS_FILE"

cat "$TEMPLATE" | \
    sed "s/{{MethodName}}/$METHOD_NAME/g" | \
    sed "s/{{RepoName}}/${REPO_NAME}/g" | \
    sed "s/{{TableName}}/$TABLE_NAME/g" | \
    sed "s/{{PayloadType}}/${PAYLOAD}/g" | \
    sed -e "/{{InColsVals}}/{r $IN_COLS_VALS_FILE" -e "d;}"

rm -f "$IN_COLS_VALS_FILE"

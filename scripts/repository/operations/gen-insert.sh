#!/bin/bash

# gen-insert.sh - Specific script for Insert operation

METHOD_NAME=$1
ENTITY=$2
TABLE_NAME=$3
IN_COLS_VALS=$4
OUT_SCAN=$7

REPO_NAME=$8
RETURN_TYPE=$9
PAYLOAD=${10}
RAW_OUT_COLS=${11}

TEMPLATE="templates/repository/insert.txt"

if [ ! -f "$TEMPLATE" ]; then
    echo "Template $TEMPLATE not found"
    exit 1
fi

IN_COLS_VALS_FILE="/tmp/gen_repo_incolvals_${ENTITY}_$METHOD_NAME.txt"
echo "$IN_COLS_VALS" > "$IN_COLS_VALS_FILE"

OUT_SCAN_FILE="/tmp/gen_repo_outscan_${ENTITY}_$METHOD_NAME.txt"
echo "$OUT_SCAN" > "$OUT_SCAN_FILE"

# Extract the base return type without slice or pointer to instantiate a zero value correctly
ZERO_VAL_TYPE=$(echo "$RETURN_TYPE" | sed -E 's/\[\]|\*//g')

# Special case for outscan: replacing "item." with "result." because insert usually has a singular result
sed -i 's/&item./\&result./g' "$OUT_SCAN_FILE"

cat "$TEMPLATE" | \
    sed "s/{{Entity}}/$ENTITY/g" | \
    sed "s/{{TableName}}/$TABLE_NAME/g" | \
    sed "s/{{OutCols}}/$OUT_COLS/g" | \
    sed "s/{{RepoName}}/${REPO_NAME}/g" | \
    sed "s/{{ReturnType}}/${RETURN_TYPE}/g" | \
    sed "s/{{PayloadType}}/${PAYLOAD}/g" | \
    sed "s/{{ZeroValType}}/${ZERO_VAL_TYPE}/g" | \
    sed "s/{{RawOutCols}}/${RAW_OUT_COLS}/g" | \
    sed -e "/{{InColsVals}}/{r $IN_COLS_VALS_FILE" -e "d;}" | \
    sed -e "/{{OutScan}}/{r $OUT_SCAN_FILE" -e "d;}"

rm -f "$IN_COLS_VALS_FILE" "$OUT_SCAN_FILE"

rm -f "$IN_VALS_FILE" "$OUT_SCAN_FILE"

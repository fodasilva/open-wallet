#!/bin/bash

# gen-repo.sh - Master script to generate repositories for a SINGLE file

MODEL_FILE=$1
if [ -z "$MODEL_FILE" ]; then
    echo "Usage: $0 <models.go>"
    exit 1
fi

if [ ! -f "$MODEL_FILE" ]; then
    echo "File $MODEL_FILE not found"
    exit 1
fi

DIR=$(dirname "$MODEL_FILE")
PACKAGE=$(grep "package " "$MODEL_FILE" | awk '{print $2}')

echo "--- Generating repository for $MODEL_FILE ---"

# Remove all previously generated files to avoid appending to old runs
rm -f "$DIR"/zz_generated_*.go

# Find all @table: tags and process them
TABLE_LINES=($(grep -n "@table:" "$MODEL_FILE" | cut -d: -f1))
NUM_TABLES=${#TABLE_LINES[@]}

for (( i=0; i<$NUM_TABLES; i++ )); do
    LINE_NUM=${TABLE_LINES[$i]}
    NEXT_LINE_NUM=""
    if [ $((i + 1)) -lt $NUM_TABLES ]; then
        NEXT_LINE_NUM=${TABLE_LINES[$((i + 1))]}
    fi

    # Define range for current block
    if [ -n "$NEXT_LINE_NUM" ]; then
        BLOCK_LENGTH=$((NEXT_LINE_NUM - LINE_NUM))
    else
        BLOCK_LENGTH=999999 # Until end of file
    fi

    TABLE=$(sed -n "${LINE_NUM}p" "$MODEL_FILE" | sed -E 's/.*@table: ([^ ]+).*/\1/')
    
    # Find block-specific metadata tags within the block scope
    ENTITY=$(sed -n "${LINE_NUM},+$((BLOCK_LENGTH - 1))p" "$MODEL_FILE" | grep "@entity:" | head -n 1 | sed -E 's/.*@entity: ([^ ]+).*/\1/')

    # If entity is empty, search 5 lines BEFORE too just in case it was placed before the tag
    if [ -z "$ENTITY" ]; then
        START_LOOK=$((LINE_NUM - 5))
        if [ $START_LOOK -lt 1 ]; then START_LOOK=1; fi
        ENTITY=$(sed -n "${START_LOOK},${LINE_NUM}p" "$MODEL_FILE" | grep "@entity:" | head -n 1 | sed -E 's/.*@entity: ([^ ]+).*/\1/')
    fi

    REPO_NAME=$(sed -n "${LINE_NUM},+$((BLOCK_LENGTH - 1))p" "$MODEL_FILE" | grep "@name:" | head -n 1 | sed -E 's/.*@name: ([^ ]+).*/\1/')

    # If repo name is empty, search 5 lines BEFORE too
    if [ -z "$REPO_NAME" ]; then
        START_LOOK=$((LINE_NUM - 5))
        if [ $START_LOOK -lt 1 ]; then START_LOOK=1; fi
        REPO_NAME=$(sed -n "${START_LOOK},${LINE_NUM}p" "$MODEL_FILE" | grep "@name:" | head -n 1 | sed -E 's/.*@name: ([^ ]+).*/\1/')
    fi

    # Find the first struct AFTER this table tag within the block
    STRUCT_INFO=$(sed -n "${LINE_NUM},+$((BLOCK_LENGTH - 1))p" "$MODEL_FILE" | grep -m 1 "type .* struct")
    STRUCT_NAME=$(echo "$STRUCT_INFO" | awk '{print $2}')

    # If entity is empty, use struct name as default
    if [ -z "$ENTITY" ]; then
        ENTITY="$STRUCT_NAME"
    fi

    # Fallback to {StructName}Repo for the interface name if not provided.
    if [ -z "$REPO_NAME" ]; then
        REPO_NAME="${STRUCT_NAME}Repo"
    fi

    echo "  > Processing entity $ENTITY (Table: $TABLE, Repo: $REPO_NAME)"

    # Find methods for THIS struct (only those within the block range)
    sed -n "${LINE_NUM},+$((BLOCK_LENGTH - 1))p" "$MODEL_FILE" | grep -n "@method:" | while read -r method_line; do
        METHOD_RAW=$(echo "$method_line" | cut -d: -f2-)
        
        # Parse Piped structure: @method: Insert | in: name:Name | out: id:ID | payload: CreateUserDTO | return: User
        METHOD_NAME=$(echo "$METHOD_RAW" | cut -d'|' -f1 | sed -E 's/.*@method: ([^ ]+).*/\1/')
        
        FIELDS_IN=$(echo "$METHOD_RAW" | grep -E -o "fields:[ ]*[^|]+" | sed -E 's/fields:[ ]*//' | xargs)
        
        FIELDS_OUT=""
        if echo "$METHOD_RAW" | grep -q "return:[ ]*"; then
            FIELDS_OUT=$(echo "$METHOD_RAW" | grep -E -o "return:[ ]*[^|]+" | sed -E 's/return:[ ]*//' | xargs)
        elif [ "$METHOD_NAME" == "Select" ]; then
            FIELDS_OUT=$FIELDS_IN
        fi
        
        # Resolve operation from method name
        OPERATION=$METHOD_NAME
        if [[ "$METHOD_NAME" == Select* ]]; then OPERATION="Select"; fi
        if [[ "$METHOD_NAME" == Insert* ]]; then OPERATION="Insert"; fi
        if [[ "$METHOD_NAME" == Update* ]]; then OPERATION="Update"; fi
        if [[ "$METHOD_NAME" == Delete* ]]; then OPERATION="Delete"; fi
        if [[ "$METHOD_NAME" == Count* ]]; then OPERATION="Count"; fi
        if [[ "$METHOD_NAME" == GetByID* ]]; then OPERATION="GetByID"; fi

        PAYLOAD=$(echo "$METHOD_RAW" | grep -E -o "payload:[ ]*[^|]+" | sed -E 's/payload:[ ]*//' | xargs)
        
        # Hardcode return type inference since 'return:' now contains fields mapping
        RETURN_TYPE=""
        if [ "$OPERATION" == "Select" ]; then
            RETURN_TYPE="[]$ENTITY"
        elif [ "$OPERATION" == "Insert" ] || [ "$OPERATION" == "GetByID" ]; then
            RETURN_TYPE="$ENTITY"
        elif [ "$OPERATION" == "Count" ]; then
            RETURN_TYPE="int"
        elif [ "$OPERATION" == "Update" ] || [ "$OPERATION" == "Delete" ]; then
            RETURN_TYPE="error"
        fi

        # Parse IN Fields
        rm -f "/tmp/${ENTITY}_in_cols_vals.txt" "/tmp/${ENTITY}_in_update.txt"
        touch "/tmp/${ENTITY}_in_cols_vals.txt" "/tmp/${ENTITY}_in_update.txt"
        # Initialize Insert slices
        echo "var columns []string" >> "/tmp/${ENTITY}_in_cols_vals.txt"
        echo "var values []interface{}" >> "/tmp/${ENTITY}_in_cols_vals.txt"

        echo "$FIELDS_IN" | tr ',' '\n' | while read -r field_pair; do
            COL=$(echo "$field_pair" | cut -d: -f1 | xargs)
            GO_FIELD_RAW=$(echo "$field_pair" | cut -d: -f2 | xargs)
            
            # Check for ? suffix meaning optional
            GO_FIELD=$(echo "$GO_FIELD_RAW" | sed 's/?//')

            if [ -n "$COL" ] && [ -n "$GO_FIELD" ]; then
                if [[ "$GO_FIELD_RAW" == *"?"* ]]; then
                    # Support for OptionalNullable: if it has '?', only set if Set is true
                    echo "if data.$GO_FIELD.Set { columns = append(columns, \"$COL\"); values = append(values, data.$GO_FIELD.Value) }" >> "/tmp/${ENTITY}_in_cols_vals.txt"
                    echo "if data.$GO_FIELD.Set { query = query.Set(\"$COL\", data.$GO_FIELD.Value) }" >> "/tmp/${ENTITY}_in_update.txt"
                else
                    echo "columns = append(columns, \"$COL\"); values = append(values, data.$GO_FIELD)" >> "/tmp/${ENTITY}_in_cols_vals.txt"
                    echo "query = query.Set(\"$COL\", data.$GO_FIELD)" >> "/tmp/${ENTITY}_in_update.txt"
                fi
            fi
        done
        # Apply Insert columns and values
        echo "query = query.Columns(columns...).Values(values...)" >> "/tmp/${ENTITY}_in_cols_vals.txt"
        IN_COLS_VALS=""
        IN_UPDATE=""
        if [ -s "/tmp/${ENTITY}_in_cols_vals.txt" ]; then
            IN_COLS_VALS=$(cat "/tmp/${ENTITY}_in_cols_vals.txt")
            IN_UPDATE=$(cat "/tmp/${ENTITY}_in_update.txt")
        fi
        rm -f "/tmp/${ENTITY}_in_cols_vals.txt" "/tmp/${ENTITY}_in_update.txt"

        # Parse OUT Fields
        rm -f "/tmp/${ENTITY}_out_cols.txt" "/tmp/${ENTITY}_out_scan.txt"
        echo "$FIELDS_OUT" | tr ',' '\n' | while read -r field_pair; do
            COL=$(echo "$field_pair" | cut -d: -f1 | xargs)
            GO_FIELD=$(echo "$field_pair" | cut -d: -f2 | xargs)
            if [ -n "$COL" ] && [ -n "$GO_FIELD" ]; then
                echo "$COL" >> "/tmp/${ENTITY}_out_cols.txt"
                echo "&item.$GO_FIELD," >> "/tmp/${ENTITY}_out_scan.txt"
            fi
        done
        OUT_COLS=""
        OUT_SCAN=""
        RAW_OUT_COLS=""
        if [ -s "/tmp/${ENTITY}_out_cols.txt" ]; then
            RAW_OUT_COLS=$(paste -sd "," "/tmp/${ENTITY}_out_cols.txt" | sed 's/,/, /g')
            OUT_COLS=$(paste -sd "," "/tmp/${ENTITY}_out_cols.txt" | sed 's/,/", "/g' | sed 's/^/"/' | sed 's/$/"/')
            OUT_SCAN=$(cat "/tmp/${ENTITY}_out_scan.txt")
        fi
        rm -f "/tmp/${ENTITY}_out_cols.txt" "/tmp/${ENTITY}_out_scan.txt"

        METHOD_LOWER=$(echo "$METHOD_NAME" | tr '[:upper:]' '[:lower:]')
        OUT_FILE="$DIR/zz_generated_${METHOD_LOWER}.go"
        
        # Create the specific file with package and imports (if it doesn't exist yet)
        if [ ! -f "$OUT_FILE" ]; then
            cat <<EOF > "$OUT_FILE"
// Code generated. DO NOT EDIT.

package $PACKAGE

import (
	"github.com/Masterminds/squirrel"
	"github.com/felipe1496/open-wallet/internal/utils"
)

EOF
        fi

        # Call the dynamically resolved sub-script and append to file
        OPERATION_LOWER=$(echo "$OPERATION" | tr '[:upper:]' '[:lower:]')
        SCRIPT_PATH="scripts/repository/operations/gen-${OPERATION_LOWER}.sh"
        if [ -f "$SCRIPT_PATH" ]; then
            bash "$SCRIPT_PATH" \
                "$METHOD_NAME" "$ENTITY" "$TABLE" "$IN_COLS_VALS" "" "$OUT_COLS" "$OUT_SCAN" "$REPO_NAME" "$RETURN_TYPE" "$PAYLOAD" "$RAW_OUT_COLS" "$IN_UPDATE" >> "$OUT_FILE"
            
            # Format the newly created file
            if command -v gofmt &> /dev/null; then
                gofmt -w "$OUT_FILE"
            fi
            echo "    - Generated $OUT_FILE"
        else
            echo "    ! Warning: Operation script not found: $SCRIPT_PATH"
        fi
    done
done

echo "Done generating for $MODEL_FILE"

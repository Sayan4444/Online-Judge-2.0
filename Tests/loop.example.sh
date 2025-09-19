#!/bin/bash

PROBLEM_ID=""
LANGUAGE=""
SOURCE_CODE_FILE="./code.txt"
TOKEN=""
URL=""

args=()
for i in {1..10}; do
    COMMAND="bash -c './runner.sh \"$PROBLEM_ID\" \"$LANGUAGE\" \"$SOURCE_CODE_FILE\" \"$TOKEN\" \"$URL\"; exec bash'"
    args+=(--tab --command="$COMMAND")
done

gnome-terminal "${args[@]}" 2>/dev/null
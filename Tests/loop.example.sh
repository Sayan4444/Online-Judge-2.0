#!/bin/bash

args=()
for i in {1..1}; do
    args+=(--tab --command="bash -c './runner.sh; exec bash'")
done

gnome-terminal "${args[@]}" 2>/dev/null
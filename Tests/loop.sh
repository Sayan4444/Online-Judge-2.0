#!/bin/bash

args=()
for i in {1..120}; do
    args+=(--tab --command="bash -c './test.sh; exec bash'")
done

gnome-terminal "${args[@]}" 2>/dev/null
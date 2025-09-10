#!/bin/bash

args=()
for i in {1..15}; do
    args+=(--tab --command="bash -c './temp.sh; exec bash'")
done

gnome-terminal "${args[@]}" 2>/dev/null
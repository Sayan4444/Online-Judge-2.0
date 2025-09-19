#!/bin/bash

# Record start time
start_time=$(date +%s.%N)

# --- Configuration from arguments ---
if [ "$#" -lt 3 ]; then
    echo "Usage: $0 <PROBLEM_ID> <LANGUAGE> <SOURCE_CODE_FILE> [TOKEN] [URL]"
    echo "Example: $0 750e8400... C++ ../Test_Case_Generator/two_sum.cpp eyJhbGci..."
    exit 1
fi

PROBLEM_ID="$1"
LANGUAGE="$2"
SOURCE_CODE_FILE="$3"
TOKEN="${4:-eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzIwNWI3MTgtZTFkZC00ZDFiLTk1YmEtNmU2NGFiNDNmZGRkIiwidXNlcm5hbWUiOiJ0ZXN0X3VzZXIiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3NTg0NDg5NzYsImlhdCI6MTc1ODE4OTc3Nn0.O5oF0N0cgvJyngRpvXQZBd14yP_SuRDdymztfoveca8}"
URL="${5:-http://localhost:8080}"

# Check if the source code file exists
if [ ! -f "$SOURCE_CODE_FILE" ]; then
    echo "Error: Source code file not found at '$SOURCE_CODE_FILE'"
    exit 1
fi

# The source code to be submitted.
SOURCE_CODE=$(cat "$SOURCE_CODE_FILE")

# --- Step 1: Submit the code via POST request ---
echo "Submitting code to $URL/api/submit/$PROBLEM_ID..."

# Construct the JSON payload using the variables.
# Using jq is a robust way to handle JSON creation.
JSON_PAYLOAD=$(jq -n \
                  --arg sc "$SOURCE_CODE" \
                  --arg lang "$LANGUAGE" \
                  '{source_code: $sc, language: $lang}')

# Perform the POST request and store the server's response.
# The -s flag silences the progress meter.
RESPONSE=$(curl -s -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$JSON_PAYLOAD" \
    "$URL/api/submit/$PROBLEM_ID")

# Check if the curl command was successful.
if [ $? -ne 0 ]; then
    echo "Error: The POST request failed. Please check the URL and network connection."
    exit 1
fi

# --- Step 2: Extract the Submission ID ---
# Use jq to parse the JSON response and extract the 'submission_id'.
# The -r flag outputs the raw string without quotes.
SUBMISSION_ID=$(echo "$RESPONSE" | jq -r '.submission_id')

# Validate the submission ID.
if [ -z "$SUBMISSION_ID" ] || [ "$SUBMISSION_ID" == "null" ]; then
    echo "Error: Could not extract submission_id from the response."
    echo "Server Response: $RESPONSE"
    exit 1
fi

echo "Submission successful. Submission ID: $SUBMISSION_ID"
echo "---------------------------------------------------"

# --- Step 3: Connect to the Server-Sent Events (SSE) stream ---
echo "Connecting to event stream..."

# Use curl to make a GET request to the events endpoint.
# -N disables buffering, which is crucial for streaming responses.
curl -N -X GET \
    -H "Authorization: Bearer $TOKEN" \
    -H "Accept: text/event-stream" \
    "$URL/api/submission/events/$SUBMISSION_ID"

echo -e "\nStream finished."

# Record end time
end_time=$(date +%s.%N)

# Calculate and display the total time
total_time=$(awk "BEGIN {print $end_time - $start_time}")
echo "Total time taken: $total_time seconds"
#!/bin/bash

# --- Configuration ---
# The base URL of your API endpoint.
URL="http://localhost:8080"

# Your JWT authorization token.
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzIwNWI3MTgtZTFkZC00ZDFiLTk1YmEtNmU2NGFiNDNmZGRkIiwidXNlcm5hbWUiOiJ0ZXN0X3VzZXIiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3NTc3ODMxMjN9.JJKnU3mNMZq5meR8g_7JO3jgmiSXqEUbRDR0bWXc3Cw"

# The unique identifier for the problem.
PROBLEM_ID="750e8400-e29b-41d4-a716-446655440001"

# The programming language of the source code.
LANGUAGE="C++"

# The source code to be submitted.
# Using a variable makes the curl command cleaner.
SOURCE_CODE=$(cat <<'EOF'
#include <iostream>
#include <vector>
#include <unordered_map>
using namespace std;

class Solution
{
public:
    vector<int> twoSum(const vector<int> &nums, int target)
    {
        unordered_map<int, int> num_map;
        for (int i = 0; i < nums.size(); ++i)
        {
            int complement = target - nums[i];
            if (num_map.find(complement) != num_map.end())
            {
                return {num_map[complement], i};
            }
            num_map[nums[i]] = i;
        }
        return {};
    }
};

int main()
{
    Solution sol;
    int n;
    vector<int> nums;
    int target;
    cin >> n;
    for (int i = 0; i < n; i++)
    {
        int num;
        cin >> num;
        nums.push_back(num);
    }
    cin >> target;
    vector<int> result = sol.twoSum(nums, target);
    if (!result.empty())
    {
        cout << result[0] << " " << result[1];
    }

    return 0;
}
EOF
)

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

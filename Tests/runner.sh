#!/bin/bash

start_time=$(date +%s)

URL="http://localhost:1323"

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzIwNWI3MTgtZTFkZC00ZDFiLTk1YmEtNmU2NGFiNDNmZGRkIiwidXNlcm5hbWUiOiJ0ZXN0X3VzZXIiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3NTg1Njk5OTcsImlhdCI6MTc1ODMxMDc5N30.Tg3U0PO3nqFfAaTpEp34pPKVcmXGbMIsZeM0P2GQAdc"

PROBLEM_ID="750e8400-e29b-41d4-a716-446655440001"

LANGUAGE="C++"
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

echo "Submitting code to $URL/api/submit/$PROBLEM_ID..."

JSON_PAYLOAD=$(jq -n \
                  --arg sc "$SOURCE_CODE" \
                  --arg lang "$LANGUAGE" \
                  '{source_code: $sc, language: $lang}')

RESPONSE=$(curl -s -X POST \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$JSON_PAYLOAD" \
    "$URL/api/submit/$PROBLEM_ID")

if [ $? -ne 0 ]; then
    echo "Error: The POST request failed. Please check the URL and network connection."
    exit 1
fi

SUBMISSION_ID=$(echo "$RESPONSE" | jq -r '.submission_id')

if [ -z "$SUBMISSION_ID" ] || [ "$SUBMISSION_ID" == "null" ]; then
    echo "Error: Could not extract submission_id from the response."
    echo "Server Response: $RESPONSE"
    exit 1
fi

echo "Submission successful. Submission ID: $SUBMISSION_ID"
echo "---------------------------------------------------"

echo "Connecting to event stream..."

curl -N -X GET \
    -H "Authorization: Bearer $TOKEN" \
    -H "Accept: text/event-stream" \
    "$URL/api/submission/events/$SUBMISSION_ID"

echo -e "\nStream finished."

# Record end time
end_time=$(date +%s)

# Calculate and display the total time
total_time=$((end_time - start_time))
echo "Total time taken: $total_time seconds"
import http from 'k6/http';
import { check, sleep } from 'k6';

// --- Read Target VUs from command line or use a default value ---
const targetVUs = __ENV.TARGET_VUS || 5;

// --- Test Configuration ---
export const options = {
    stages: [
        // Ramp-up: Gradually increase virtual users to the target value over 1 minute.
        { duration: '1m', target: targetVUs },

        // Hold Load: Maintain the target number of virtual users for 1 minute.
        { duration: '2m', target: targetVUs },

        // Ramp-down: Gradually decrease users back to 0 over 1 minute.
        { duration: '1m', target: 0 },
    ],
    thresholds: {
        'http_req_failed': ['rate<0.01'], // less than 1% of requests should fail
    },
};

// --- Test Data ---
// const BASE_URL = 'http://localhost:8080';
const BASE_URL = 'http://64.225.84.213:80';
const TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzIwNWI3MTgtZTFkZC00ZDFiLTk1YmEtNmU2NGFiNDNmZGRkIiwidXNlcm5hbWUiOiJ0ZXN0X3VzZXIiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE3NTg1Njk5OTcsImlhdCI6MTc1ODMxMDc5N30.Tg3U0PO3nqFfAaTpEp34pPKVcmXGbMIsZeM0P2GQAdc';
const PROBLEM_ID = '750e8400-e29b-41d4-a716-446655440001';
const LANGUAGE = 'C++';
const SOURCE_CODE = `
#include <iostream>
#include <vector>
#include <unordered_map>
using namespace std;

class Solution {
public:
    vector<int> twoSum(const vector<int> &nums, int target) {
        unordered_map<int, int> num_map;
        for (int i = 0; i < nums.size(); ++i) {
            int complement = target - nums[i];
            if (num_map.find(complement) != num_map.end()) {
                return {num_map[complement], i};
            }
            num_map[nums[i]] = i;
        }
        return {};
    }
};

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);

    Solution sol;
    int t;
    cin >> t;
    while (t--) {
        int n;
        cin >> n;
        vector<int> nums(n);
        int target;
        for (int i = 0; i < n; i++) {
            cin >> nums[i];
        }
        cin >> target;
        vector<int> result = sol.twoSum(nums, target);
        if (!result.empty()) {
            cout << result[0] << " " << result[1] << "\\n";
        }
    }
    return 0;
}
`;


// --- Main Test Logic ---
export default function () {
    const startTime = Date.now();

    const submitUrl = `${BASE_URL}/api/submit/${PROBLEM_ID}`;

    const payload = JSON.stringify({
        source_code: SOURCE_CODE,
        language: LANGUAGE,
    });

    const params = {
        headers: {
            'Authorization': `Bearer ${TOKEN}`,
            'Content-Type': 'application/json',
        },
    };

    // 1. Submit the code via a POST request
    const submitRes = http.post(submitUrl, payload, params);

    const submissionSuccess = check(submitRes, {
        'submission POST status is 200 or 201': (r) => r.status === 200 || r.status === 201,
        'submission response has submission_id': (r) => r.json('submission_id') !== undefined,
    });

    if (!submissionSuccess) {
        console.error(`Submission failed. VU: ${__VU}, Iter: ${__ITER}. Response: ${submitRes.body}`);
        return;
    }

    const submissionId = submitRes.json('submission_id');

    // 2. Connect to the event stream to get results
    const eventsUrl = `${BASE_URL}/api/submission/events/${submissionId}`;
    const eventParams = {
        headers: {
            'Authorization': `Bearer ${TOKEN}`,
            'Accept': 'text-event-stream',
        },
    };

    const eventRes = http.get(eventsUrl, eventParams);

    check(eventRes, {
        'event stream GET status is 200': (r) => r.status === 200,
    });

    // --- MODIFIED: Calculate duration and log protocol ---
    const endTime = Date.now();
    const durationInSeconds = (endTime - startTime) / 1000;
    
    // Get protocol from each response
    const submitProto = submitRes.proto;
    const eventProto = eventRes.proto;

    console.log(`✅ Iteration complete. VU: ${__VU} | Sub ID: ${submissionId} | Time: ${durationInSeconds.toFixed(3)}s | Protocols (Submit/Events): ${submitProto}/${eventProto}`);

    // 3. Pause for a short time to simulate a real user's think time
    sleep(1);
}

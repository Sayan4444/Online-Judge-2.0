package isolatejob_test

import (
	"context"
	"fmt"
	"testing"

	uuid "github.com/google/uuid"

	"OJ-Worker/config"
	isolateJob "OJ-Worker/isolateJob"
	"OJ-Worker/schema"
)


func TestProcessSubmission(t *testing.T) {
	code := "#include <bits/stdc++.h>\nusing namespace std;\n\nclass Solution {\n    vector<int> val, wt;\n    int n, W;\n    vector<vector<int>> dp;\n\n    // Recursive function with memoization\n    int f(int ind, int weight) {\n        if (ind == n) return 0;  // Base case: no items left\n\n        if (dp[ind][weight] != -1) return dp[ind][weight]; // Already computed\n\n        // Option 1: Do not take this item\n        int nontake = f(ind + 1, weight);\n\n        // Option 2: Take this item (only if it fits)\n        int take = 0;\n        if (weight + wt[ind] <= W)\n            take = val[ind] + f(ind + 1, weight + wt[ind]);\n\n        // Store and return the maximum\n        return dp[ind][weight] = max(take, nontake);\n    }\n\npublic:\n    int knapsack(int W, vector<int> &val, vector<int> &wt) {\n        this->W = W;\n        this->val = val;\n        this->wt = wt;\n        n = val.size();\n\n        // Initialize DP with -1\n        dp.assign(n + 1, vector<int>(W + 1, -1));\n\n        return f(0, 0);\n    }\n};\n\n// ---------------- DRIVER CODE ----------------\nint main() {\n    // Fast I/O\n    ios_base::sync_with_stdio(false);\n    cin.tie(NULL);\n\n    int t;\n    cin >> t; // Read the number of test cases\n\n    while (t--) { // Loop for each test case\n        int n, W;\n        cin >> n;\n        vector<int> val(n), wt(n);\n\n        for (int i = 0; i < n; i++) cin >> val[i];\n\n        for (int i = 0; i < n; i++) cin >> wt[i];\n\n        cin >> W;\n\n        Solution obj;\n        int maxProfit = obj.knapsack(W, val, wt);\n\n        cout << maxProfit << \"\\n\"; // Print the result for the current test case\n    }\n    return 0;\n}"

    payload := schema.RabbitMQPayload{
        SubmissionID:   uuid.MustParse("f84b84f4-7053-4ece-a484-d96c4e5fdccd"),
        ProblemID:     	uuid.MustParse("56d4fc14-a544-4910-916d-b3d7968ffb2e"),	
        UserID:         uuid.MustParse("f84b84f4-7053-4ece-a484-d96c4e5fdccd"),
		Language:       "cpp",
		SourceCode:     code,
	}
	_, err := config.ConnectDB()
	if err != nil {
		t.Fatal("Expected no error but got ", err)
	}
	response := schema.JudgeResponse{}
	ctx := context.Background()
	err = isolateJob.ProcessSubmission(&payload, &response, ctx)
	fmt.Println()
	fmt.Println("Response:-", response)
	fmt.Println()
	if err != nil {
		t.Fatal("Expected no error but got ", err)
	}
}


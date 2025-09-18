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
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);

    Solution sol;
    // int t;
    // cin >> t;
    // while (t--)
    {
        int n;
        cin >> n;
        vector<int> nums(n);
        int target;
        for (int i = 0; i < n; i++)
        {
            cin >> nums[i];
        }
        cin >> target;
        vector<int> result = sol.twoSum(nums, target);
        if (!result.empty())
        {
            cout << result[0] << " " << result[1] << "\n"; 
        }
    }
    return 0;
}

-- Sample data for Online Judge Platform

-- Create database schema first
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    oauth_id VARCHAR UNIQUE,
    provider VARCHAR,
    image VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create contests table
CREATE TABLE IF NOT EXISTS contests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create problems table
CREATE TABLE IF NOT EXISTS problems (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contest_id UUID NOT NULL,
    title VARCHAR NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE
);

-- Create submissions table
CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    problem_id UUID NOT NULL,
    user_id UUID NOT NULL,
    contest_id UUID NOT NULL,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    result VARCHAR NOT NULL,
    language VARCHAR NOT NULL,
    source_code TEXT NOT NULL,
    score INTEGER DEFAULT 0,
    std_output TEXT,
    std_error TEXT,
    wrong_test_case UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    compile_output TEXT,
    exit_code INTEGER,
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE
);

-- Create test_cases table
CREATE TABLE IF NOT EXISTS test_cases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    problem_id UUID NOT NULL,
    input TEXT NOT NULL,
    output TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);

-- Create languages table
CREATE TABLE IF NOT EXISTS languages (
    id INTEGER PRIMARY KEY,
    name VARCHAR NOT NULL,
    compile_command VARCHAR NOT NULL,
    run_command VARCHAR NOT NULL,
    time_limit INTEGER NOT NULL,
    memory_limit INTEGER,
    wall_limit INTEGER,
    stack_limit INTEGER,
    output_limit INTEGER,
    src_file VARCHAR NOT NULL
);

-- Create contest_users junction table for many-to-many relationship
CREATE TABLE IF NOT EXISTS contest_users (
    contest_id UUID,
    user_id UUID,
    PRIMARY KEY (contest_id, user_id),
    FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Insert sample languages
INSERT INTO languages (id, name, compile_command, run_command, time_limit, memory_limit, wall_limit, stack_limit, output_limit, src_file) VALUES
(1, 'C++', 'g++ -o solution solution.cpp -std=c++17', './solution', 2000, 256, 5, 8, 64, 'solution.cpp'),
(2, 'Java', 'javac Solution.java', 'java Solution', 3000, 512, 10, 8, 64, 'Solution.java'),
(3, 'Python', '', 'python3 solution.py', 5000, 256, 10, 8, 64, 'solution.py'),
(4, 'C', 'gcc -o solution solution.c', './solution', 2000, 256, 5, 8, 64, 'solution.c'),
(5, 'JavaScript', '', 'node solution.js', 3000, 256, 8, 8, 64, 'solution.js');

-- Insert sample users
INSERT INTO users (id, username, email, oauth_id, provider, image, created_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'alice_coder', 'alice@example.com', 'github_123456', 'github', 'https://avatar.githubusercontent.com/alice', '2024-01-15 10:30:00'),
('550e8400-e29b-41d4-a716-446655440002', 'bob_solver', 'bob@example.com', 'google_789012', 'google', 'https://lh3.googleusercontent.com/bob', '2024-01-16 14:20:00'),
('550e8400-e29b-41d4-a716-446655440003', 'charlie_dev', 'charlie@example.com', 'github_345678', 'github', 'https://avatar.githubusercontent.com/charlie', '2024-01-17 09:15:00'),
('550e8400-e29b-41d4-a716-446655440004', 'diana_prog', 'diana@example.com', 'google_901234', 'google', 'https://lh3.googleusercontent.com/diana', '2024-01-18 16:45:00'),
('550e8400-e29b-41d4-a716-446655440005', 'eve_hacker', 'eve@example.com', 'github_567890', 'github', 'https://avatar.githubusercontent.com/eve', '2024-01-19 11:30:00');

-- Insert sample contests
INSERT INTO contests (id, name, description, start_time, end_time, created_at) VALUES
('650e8400-e29b-41d4-a716-446655440001', 'Weekly Programming Contest #1', 'A beginner-friendly contest with basic algorithmic problems', '2024-03-01 10:00:00', '2024-03-01 13:00:00', '2024-02-25 15:30:00'),
('650e8400-e29b-41d4-a716-446655440002', 'Advanced Algorithms Challenge', 'Advanced contest featuring dynamic programming and graph algorithms', '2024-03-08 14:00:00', '2024-03-08 18:00:00', '2024-03-01 12:00:00'),
('650e8400-e29b-41d4-a716-446655440003', 'Data Structures Showdown', 'Contest focused on implementation of various data structures', '2024-03-15 09:00:00', '2024-03-15 12:00:00', '2024-03-08 10:15:00'),
('650e8400-e29b-41d4-a716-446655440004', 'Math and Logic Contest', 'Mathematical problem-solving contest', '2024-03-22 11:00:00', '2024-03-22 15:00:00', '2024-03-15 13:45:00'),
('650e8400-e29b-41d4-a716-446655440005', 'Speed Coding Challenge', 'Fast-paced contest with simple but tricky problems', '2024-03-29 16:00:00', '2024-03-29 18:00:00', '2024-03-22 14:20:00');

-- Insert sample problems
INSERT INTO problems (id, contest_id, title, description, created_at) VALUES
('750e8400-e29b-41d4-a716-446655440001', '650e8400-e29b-41d4-a716-446655440001', 'Two Sum', 'Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target. You may assume that each input would have exactly one solution, and you may not use the same element twice.', '2024-02-25 16:00:00'),
('750e8400-e29b-41d4-a716-446655440002', '650e8400-e29b-41d4-a716-446655440001', 'Palindrome Check', 'Given a string, determine if it is a palindrome, considering only alphanumeric characters and ignoring cases.', '2024-02-25 16:15:00'),
('750e8400-e29b-41d4-a716-446655440003', '650e8400-e29b-41d4-a716-446655440002', 'Longest Common Subsequence', 'Given two strings text1 and text2, return the length of their longest common subsequence.', '2024-03-01 12:30:00'),
('750e8400-e29b-41d4-a716-446655440004', '650e8400-e29b-41d4-a716-446655440002', 'Dijkstra Shortest Path', 'Implement Dijkstra algorithm to find the shortest path between two nodes in a weighted graph.', '2024-03-01 12:45:00'),
('750e8400-e29b-41d4-a716-446655440005', '650e8400-e29b-41d4-a716-446655440003', 'Binary Search Tree Implementation', 'Implement a binary search tree with insert, delete, and search operations.', '2024-03-08 10:30:00'),
('750e8400-e29b-41d4-a716-446655440006', '650e8400-e29b-41d4-a716-446655440004', 'Prime Factorization', 'Given a positive integer n, return all the prime factors of n.', '2024-03-15 14:00:00'),
('750e8400-e29b-41d4-a716-446655440007', '650e8400-e29b-41d4-a716-446655440005', 'Array Rotation', 'Given an array, rotate the array to the right by k steps, where k is non-negative.', '2024-03-22 14:35:00');

-- Insert sample test cases
INSERT INTO test_cases (id, problem_id, input, output, created_at) VALUES
-- Test cases for Two Sum
('850e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', '4\n2 7 11 15\n9', '0 1', '2024-02-25 16:05:00'),
('850e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440001', '3\n3 2 4\n6', '1 2', '2024-02-25 16:05:00'),
('850e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440001', '2\n3 3\n6', '0 1', '2024-02-25 16:05:00'),

-- Test cases for Palindrome Check
('850e8400-e29b-41d4-a716-446655440004', '750e8400-e29b-41d4-a716-446655440002', 'A man a plan a canal Panama', 'true', '2024-02-25 16:20:00'),
('850e8400-e29b-41d4-a716-446655440005', '750e8400-e29b-41d4-a716-446655440002', 'race a car', 'false', '2024-02-25 16:20:00'),
('850e8400-e29b-41d4-a716-446655440006', '750e8400-e29b-41d4-a716-446655440002', 'Was it a car or a cat I saw', 'true', '2024-02-25 16:20:00'),

-- Test cases for Longest Common Subsequence
('850e8400-e29b-41d4-a716-446655440007', '750e8400-e29b-41d4-a716-446655440003', 'abcde\nace', '3', '2024-03-01 12:35:00'),
('850e8400-e29b-41d4-a716-446655440008', '750e8400-e29b-41d4-a716-446655440003', 'abc\nabc', '3', '2024-03-01 12:35:00'),
('850e8400-e29b-41d4-a716-446655440009', '750e8400-e29b-41d4-a716-446655440003', 'abc\ndef', '0', '2024-03-01 12:35:00'),

-- Test cases for Prime Factorization
('850e8400-e29b-41d4-a716-446655440010', '750e8400-e29b-41d4-a716-446655440006', '12', '2 2 3', '2024-03-15 14:05:00'),
('850e8400-e29b-41d4-a716-446655440011', '750e8400-e29b-41d4-a716-446655440006', '15', '3 5', '2024-03-15 14:05:00'),
('850e8400-e29b-41d4-a716-446655440012', '750e8400-e29b-41d4-a716-446655440006', '17', '17', '2024-03-15 14:05:00'),

-- Test cases for Array Rotation
('850e8400-e29b-41d4-a716-446655440013', '750e8400-e29b-41d4-a716-446655440007', '7\n1 2 3 4 5 6 7\n3', '5 6 7 1 2 3 4', '2024-03-22 14:40:00'),
('850e8400-e29b-41d4-a716-446655440014', '750e8400-e29b-41d4-a716-446655440007', '3\n-1 -100 3\n2', '3 -1 -100', '2024-03-22 14:40:00'),
('850e8400-e29b-41d4-a716-446655440015', '750e8400-e29b-41d4-a716-446655440007', '2\n1 2\n1', '2 1', '2024-03-22 14:40:00');

-- Insert contest_users relationships (many-to-many)
INSERT INTO contest_users (contest_id, user_id) VALUES
('650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001'),
('650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440002'),
('650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440003'),
('650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001'),
('650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440004'),
('650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440005'),
('650e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440002'),
('650e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440003'),
('650e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440005'),
('650e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440001'),
('650e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440003'),
('650e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440004'),
('650e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440002'),
('650e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440004'),
('650e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440005');

-- Insert sample submissions using dollar-quoted strings to avoid escaping issues
INSERT INTO submissions (id, problem_id, user_id, contest_id, submitted_at, result, language, source_code, score, std_output, std_error, wrong_test_case, created_at, compile_output, exit_code) VALUES
('950e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', '650e8400-e29b-41d4-a716-446655440001', '2024-03-01 10:15:00', 'AC', 'C++', $$#include<iostream>
#include<vector>
#include<unordered_map>
using namespace std;

int main(){
    int n, target;
    cin >> n;
    vector<int> nums(n);
    for(int i = 0; i < n; i++) cin >> nums[i];
    cin >> target;
    
    unordered_map<int, int> mp;
    for(int i = 0; i < n; i++){
        if(mp.count(target - nums[i])){
            cout << mp[target - nums[i]] << " " << i << endl;
            return 0;
        }
        mp[nums[i]] = i;
    }
    return 0;
}$$, 100, '0 1', '', '00000000-0000-0000-0000-000000000000', '2024-03-01 10:15:00', 'Compilation successful', 0),

('950e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', '650e8400-e29b-41d4-a716-446655440001', '2024-03-01 10:30:00', 'AC', 'Python', $$def is_palindrome(s):
    cleaned = "".join(char.lower() for char in s if char.isalnum())
    return cleaned == cleaned[::-1]

s = input().strip()
print("true" if is_palindrome(s) else "false")$$, 100, 'true', '', '00000000-0000-0000-0000-000000000000', '2024-03-01 10:30:00', '', 0),

('950e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440003', '650e8400-e29b-41d4-a716-446655440001', '2024-03-01 10:45:00', 'WA', 'Java', $$import java.util.*;
public class Solution {
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        int n = sc.nextInt();
        int[] nums = new int[n];
        for(int i = 0; i < n; i++) nums[i] = sc.nextInt();
        int target = sc.nextInt();
        
        for(int i = 0; i < n; i++){
            for(int j = i+1; j < n; j++){
                if(nums[i] + nums[j] == target){
                    System.out.println(j + " " + i); // Wrong order!
                    return;
                }
            }
        }
    }
}$$, 0, '1 0', '', '850e8400-e29b-41d4-a716-446655440001', '2024-03-01 10:45:00', 'Compilation successful', 0),

('950e8400-e29b-41d4-a716-446655440004', '750e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440001', '650e8400-e29b-41d4-a716-446655440002', '2024-03-08 14:20:00', 'AC', 'C++', $$#include<iostream>
#include<vector>
#include<string>
using namespace std;

int lcs(string text1, string text2) {
    int m = text1.length(), n = text2.length();
    vector<vector<int>> dp(m + 1, vector<int>(n + 1, 0));
    
    for (int i = 1; i <= m; i++) {
        for (int j = 1; j <= n; j++) {
            if (text1[i-1] == text2[j-1]) {
                dp[i][j] = dp[i-1][j-1] + 1;
            } else {
                dp[i][j] = max(dp[i-1][j], dp[i][j-1]);
            }
        }
    }
    return dp[m][n];
}

int main() {
    string text1, text2;
    cin >> text1 >> text2;
    cout << lcs(text1, text2) << endl;
    return 0;
}$$, 100, '3', '', '00000000-0000-0000-0000-000000000000', '2024-03-08 14:20:00', 'Compilation successful', 0),

('950e8400-e29b-41d4-a716-446655440005', '750e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440004', '650e8400-e29b-41d4-a716-446655440004', '2024-03-22 11:25:00', 'AC', 'Python', $$def prime_factors(n):
    factors = []
    d = 2
    while d * d <= n:
        while n % d == 0:
            factors.append(d)
            n //= d
        d += 1
    if n > 1:
        factors.append(n)
    return factors

n = int(input())
result = prime_factors(n)
print(" ".join(map(str, result)))$$, 100, '2 2 3', '', '00000000-0000-0000-0000-000000000000', '2024-03-22 11:25:00', '', 0),

('950e8400-e29b-41d4-a716-446655440006', '750e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440005', '650e8400-e29b-41d4-a716-446655440005', '2024-03-29 16:15:00', 'AC', 'JavaScript', $$const readline = require("readline");
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

let lines = [];
rl.on("line", (line) => {
    lines.push(line);
    if (lines.length === 3) {
        const n = parseInt(lines[0]);
        const nums = lines[1].split(" ").map(Number);
        const k = parseInt(lines[2]);
        
        function rotate(nums, k) {
            k = k % nums.length;
            return nums.slice(-k).concat(nums.slice(0, -k));
        }
        
        const result = rotate(nums, k);
        console.log(result.join(" "));
        rl.close();
    }
});$$, 100, '5 6 7 1 2 3 4', '', '00000000-0000-0000-0000-000000000000', '2024-03-29 16:15:00', '', 0),

('950e8400-e29b-41d4-a716-446655440007', '750e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440003', '650e8400-e29b-41d4-a716-446655440001', '2024-03-01 11:00:00', 'CE', 'C++', $$#include<iostream>
#include<string>
using namespace std;

int main(){
    string s;
    getline(cin, s);
    
    string cleaned = "";
    for(char c : s){
        if(isalnum(c)){
            cleaned += tolower(c);
        }
    }
    
    string reversed = cleaned;
    reverse(reversed.begin(), reversed.end()); // Missing #include<algorithm>
    
    if(cleaned == reversed){
        cout << "true" << endl;
    } else {
        cout << "false" << endl;
    }
    
    return 0;
}$$, 0, '', $$solution.cpp:15:5: error: 'reverse' was not declared in this scope$$, '00000000-0000-0000-0000-000000000000', '2024-03-01 11:00:00', $$Compilation failed: solution.cpp:15:5: error: 'reverse' was not declared in this scope
15 |     reverse(reversed.begin(), reversed.end());
   |     ^~~~~~~$$, 1),

('950e8400-e29b-41d4-a716-446655440008', '750e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440002', '650e8400-e29b-41d4-a716-446655440005', '2024-03-29 16:30:00', 'TLE', 'Python', $$n = int(input())
nums = list(map(int, input().split()))
k = int(input())

# Inefficient approach - rotating one by one
for _ in range(k):
    last = nums.pop()
    nums.insert(0, last)

print(" ".join(map(str, nums)))$$, 0, '', 'Time Limit Exceeded', '850e8400-e29b-41d4-a716-446655440013', '2024-03-29 16:30:00', '', 124);
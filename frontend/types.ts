type ContestType = {
  id: string;
  name: string;
  description: string;
  start_time: Date;
  end_time: Date;
  created_at: Date;
};

type ProblemType = {
  id: string;
  title: string;
  description: string;
  contest_id: string;
  created_at: Date;
  tests?: TestcaseType[];
};

type TestcaseType = {
  id: string;
  input: string;
  output: string;
  problem_id: string;
  created_at: Date;
};

type LeaderboardEntryType = {
  user_id: string;
  username: string;
  total_score: number;
};

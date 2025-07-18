import { API_URL } from "@/lib/apiEndpoints";
import axios from "axios";

export const fetchLeaderboard = async (
  contestId: string,
  token: string
): Promise<LeaderboardEntryType[]> => {
  if (!contestId || !token) return [];
  try {
    const response = await axios.get(`${API_URL}/leaderboard/${contestId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    const leaderboardData = response.data;
    return leaderboardData ?? [];
  } catch (error) {
    console.error("Fetch leaderboard error:", error);
    return [];
  }
};

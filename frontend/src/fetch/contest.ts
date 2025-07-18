import { ADMIN_URL, API_URL, BASE_URL } from "@/lib/apiEndpoints";
import axios from "axios";

interface UpdateContestType {
  name?: string;
  description?: string;
  start_time?: Date;
  end_time?: Date;
}

export const updateContest = async (
  token: string,
  payload: UpdateContestType,
  contestId: string
): Promise<UpdateContestType | null> => {
  if (!token) {
    console.error("User token is not available");
    return null;
  }
  try {
    const response = await axios.put(
      `${ADMIN_URL}/contest/${contestId}`,
      payload,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    return response.data;
  } catch (error) {
    console.error("Update contest error:", error);
    return null;
  }
};

export const deleteContest = async (
  token: string,
  contestId: string
): Promise<boolean> => {
  if (!token) {
    console.error("User token is not available");
    return false;
  }
  try {
    await axios.delete(`${ADMIN_URL}/contest/${contestId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return true;
  } catch (error) {
    console.error("Delete contest error:", error);
    return false;
  }
};

export const fetchContests = async (): Promise<ContestType[]> => {
  try {
    const response = await axios.get(`${BASE_URL}/contests`);
    const contests = response.data;
    return contests ?? [];
  } catch (error) {
    console.error("Fetch contests error:", error);
    return [];
  }
};

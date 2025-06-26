import { ADMIN_URL } from "@/lib/apiEndpoints";
import axios from "axios";

export const fetchContests = async (token: string): Promise<ContestType[]> => {
  if (!token) {
    console.error("User token is not available");
    return [];
  }
  try {
    const response = await axios.get(`${ADMIN_URL}/contests`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    const contests = response.data;
    return contests ?? [];
  } catch (error) {
    console.error("Fetch contests error:", error);
    return [];
  }
};

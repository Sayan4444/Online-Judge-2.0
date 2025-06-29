import { API_URL } from "@/lib/apiEndpoints";
import axios from "axios";

export const submitCode = async (
  problemId: string,
  userId: string,
  token: string,
  code: string,
  language: string
) => {
  try {
    const payload = {
      source_code: code,
      result: "pending",
      language: language,
      score: 0,
    };
    const data = await axios.post(
      `${API_URL}/submit/${userId}/${problemId}`,
      payload,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    console.log("Code submitted successfully:", data.data);
  } catch (error) {
    console.error("Error submitting code:", error);
    throw new Error("Failed to submit code");
  }
};

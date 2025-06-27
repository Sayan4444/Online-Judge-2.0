import { ADMIN_URL } from "@/lib/apiEndpoints";
import axios from "axios";

interface CreateProblemType {
  title: string;
  description: string;
}

export const fetchProblemsByContestId = async (
  id: string,
  token: string
): Promise<ProblemType[]> => {
  const contestId = id;
  try {
    const data = await axios.get(`${ADMIN_URL}/problems/${contestId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    const problems = data.data;
    return problems ?? [];
  } catch (error) {
    console.error("Fetch problems error:", error);
    return [];
  }
};

export const CreateProblemByContestId = async (
  id: string,
  payload: CreateProblemType,
  token: string
) => {
  const contestId = id;
  try {
    const data = await axios.post(
      `${ADMIN_URL}/create-problem/${contestId}`,
      payload,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    return data.data;
  } catch (error) {
    console.error("Create problem error:", error);
    return null;
  }
};

export const updateProblemByProblemId = async (
  problemId: string,
  payload: CreateProblemType,
  token: string
): Promise<CreateProblemType | null> => {
  try {
    const data = await axios.put(`${ADMIN_URL}/problem/${problemId}`, payload, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return data.data;
  } catch (error) {
    console.error("Update problem error:", error);
    return null;
  }
};

export const deleteProblemByProblemId = async (
  problemId: string,
  token: string
): Promise<boolean> => {
  try {
    await axios.delete(`${ADMIN_URL}/problem/${problemId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return true;
  } catch (error) {
    console.error("Delete problem error:", error);
    return false;
  }
};

import { ADMIN_URL } from "@/lib/apiEndpoints";
import axios from "axios";

export const fetchTestCasesByProblemId = async (
  problemId: string,
  token: string
): Promise<TestcaseType[]> => {
  if (!token) {
    console.error("User token is not available");
    return [];
  }
  try {
    const data = await axios.get(`${ADMIN_URL}/testcases/${problemId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    const testCases = data.data as TestcaseType[];
    return testCases ?? [];
  } catch (error) {
    console.error("Fetch test cases error:", error);
    return [];
  }
};

export const createTestCaseByProblemID = async (
  problemId: string,
  payload: { input: string; output: string },
  token: string
) => {
  if (!token) {
    console.error("User token is not available");
    return null;
  }
  try {
    const response = await axios.post(
      `${ADMIN_URL}/create-testcase/${problemId}`,
      payload,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    const data = await response.data;
    return data;
  } catch (error) {
    console.error("Create test case error:", error);
    return null;
  }
};

export const updateTestCaseByTestCaseID = async (
  testCaseId: string,
  payload: { input: string; output: string },
  token: string
): Promise<TestcaseType | null> => {
  if (!token) {
    console.error("User token is not available");
    return null;
  }
  try {
    const response = await axios.put(
      `${ADMIN_URL}/testcase/${testCaseId}`,
      payload,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    return response.data;
  } catch (error) {
    console.error("Update test case error:", error);
    return null;
  }
};

export const deleteTestCaseByTestCaseID = async (
  testCaseId: string,
  token: string
): Promise<boolean> => {
  if (!token) {
    console.error("User token is not available");
    return false;
  }
  try {
    await axios.delete(`${ADMIN_URL}/testcase/${testCaseId}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    return true;
  } catch (error) {
    console.error("Delete test case error:", error);
    return false;
  }
};

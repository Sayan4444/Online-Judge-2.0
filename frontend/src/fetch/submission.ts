import { API_URL, BASE_URL } from "@/lib/apiEndpoints";
import axios from "axios";

export interface SubmissionResponse {
  id: string;
  problem_id: string;
  user_id: string;
  contest_id: string;
  submitted_at: string;
  result: string;
  language: string;
  source_code: string;
  score: number;
  callback_url: string;
}

export interface SubmissionUpdate {
  submission_id: string;
  result: string;
  score: number;
  std_output: string;
  std_error: string;
  compile_output: string;
  exit_signal: number;
  exit_code: number;
  time: string;
  memory: string;
  message: string;
  status: string;
}

export const submitCode = async (
  problemId: string,
  userId: string,
  token: string,
  code: string,
  language: string
): Promise<SubmissionResponse> => {
  try {
    const payload = {
      source_code: code,
      result: "pending",
      language: language,
      score: 0,
    };
    const response = await axios.post(
      `${API_URL}/submit/${userId}/${problemId}`,
      payload,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    console.log("Code submitted successfully:", response.data);
    return response.data;
  } catch (error) {
    console.error("Error submitting code:", error);
    throw new Error("Failed to submit code");
  }
};

export const subscribeToSubmissionUpdates = (
  userId: string,
  submissionId: string,
  token: string,
  onUpdate: (update: SubmissionUpdate) => void,
  onError?: (error: Error) => void,
  onComplete?: () => void
): EventSource => {
  const eventSource = new EventSource(
    `${BASE_URL}/submission/${userId}/${submissionId}/events`,
    {
      withCredentials: false,
    }
  );

  eventSource.onmessage = (event) => {
    try {
      const update: SubmissionUpdate = JSON.parse(event.data);
      console.log("Received SSE update:", update);

      if (update.status === "connected") {
        console.log("Connected to submission updates");
        return;
      }

      onUpdate(update);

      // Close connection after receiving the final result
      if (update.status === "completed") {
        eventSource.close();
        onComplete?.();
      }
    } catch (error) {
      console.error("Error parsing SSE data:", error);
      onError?.(new Error("Failed to parse submission update"));
    }
  };

  eventSource.onerror = (event) => {
    console.error("SSE connection error:", event);
    onError?.(new Error("Connection to submission updates failed"));
    eventSource.close();
  };

  // Auto-close connection after 5 minutes to prevent hanging connections
  setTimeout(() => {
    if (eventSource.readyState !== EventSource.CLOSED) {
      console.log("Auto-closing SSE connection after timeout");
      eventSource.close();
      onComplete?.();
    }
  }, 5 * 60 * 1000);

  return eventSource;
};

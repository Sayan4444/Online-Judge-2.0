"use client";
import React, { useState, useRef } from "react";
import AceEditor from "react-ace";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "../ui/button";

import "ace-builds/src-noconflict/mode-java";
import "ace-builds/src-noconflict/mode-c_cpp";
import "ace-builds/src-noconflict/mode-python";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/ext-language_tools";
import { CustomUser } from "@/app/api/auth/[...nextauth]/options";
import {
  submitCode,
  subscribeToSubmissionUpdates,
  SubmissionUpdate,
} from "@/fetch/submission";
import { toast } from "sonner";

interface SubmissionResult {
  result: string;
  score: number;
  message: string;
  std_output?: string;
  std_error?: string;
  compile_output?: string;
  time?: string;
  memory?: string;
}

const CodeEditor = ({
  user,
  problem,
}: {
  user: CustomUser;
  problem: ProblemType;
}) => {
  const [language, setLanguage] = useState("c_cpp");
  const [code, setCode] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submissionResult, setSubmissionResult] =
    useState<SubmissionResult | null>(null);
  const eventSourceRef = useRef<EventSource | null>(null);

  const handleCodeChange = (newValue: string) => {
    setCode(newValue);
  };

  const handleSubmit = async () => {
    if (!user || !problem) {
      console.error("User or problem data is missing");
      return;
    }

    if (code === "") {
      toast.info("No code to submit");
      return;
    }

    setIsSubmitting(true);
    setSubmissionResult(null);

    try {
      // Submit the code and get submission details
      const submission = await submitCode(
        problem.id,
        user.id!,
        user.token!,
        code,
        language
      );

      console.log("Submission created:", submission);
      toast.success("Code submitted successfully! Waiting for results...");

      // Subscribe to real-time updates
      eventSourceRef.current = subscribeToSubmissionUpdates(
        user.id!,
        submission.id,
        user.token!,
        (update: SubmissionUpdate) => {
          console.log("Received submission update:", update);

          // Update the submission result state
          setSubmissionResult({
            result: update.result,
            score: update.score,
            message: update.message,
            std_output: update.std_output,
            std_error: update.std_error,
            compile_output: update.compile_output,
            time: update.time,
            memory: update.memory,
          });

          // Show appropriate toast based on result
          switch (update.result) {
            case "AC":
              toast.success(`Accepted! Score: ${update.score}/100`);
              break;
            case "WA":
              toast.error("Wrong Answer");
              break;
            case "TLE":
              toast.error("Time Limit Exceeded");
              break;
            case "CE":
              toast.error("Compilation Error");
              break;
            case "RE":
              toast.error("Runtime Error");
              break;
            default:
              toast.info(`Result: ${update.result}`);
          }
        },
        (error: Error) => {
          console.error("SSE Error:", error);
          toast.error("Connection error while waiting for results");
          setIsSubmitting(false);
        },
        () => {
          // On completion
          console.log("SSE connection completed");
          setIsSubmitting(false);
        }
      );
    } catch (error) {
      console.error("Error submitting code:", error);
      toast.error("Failed to submit code. Please try again.");
      setIsSubmitting(false);
    }
  };

  // Clean up SSE connection on component unmount
  React.useEffect(() => {
    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }
    };
  }, []);

  const getResultColor = (result: string) => {
    switch (result) {
      case "AC":
        return "text-green-600";
      case "WA":
        return "text-red-600";
      case "TLE":
        return "text-yellow-600";
      case "CE":
        return "text-orange-600";
      case "RE":
        return "text-red-600";
      default:
        return "text-gray-600";
    }
  };

  return (
    <div className="w-full space-y-4">
      <div className="flex justify-between items-center">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline">
              {language === "c_cpp"
                ? "C++"
                : language === "java"
                ? "Java"
                : "Python"}
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem onClick={() => setLanguage("c_cpp")}>
              C++
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => setLanguage("java")}>
              Java
            </DropdownMenuItem>
            <DropdownMenuItem onClick={() => setLanguage("python")}>
              Python
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <Button
          onClick={handleSubmit}
          disabled={isSubmitting || code === ""}
          className="bg-blue-600 hover:bg-blue-700"
        >
          {isSubmitting ? "Submitting..." : "Submit Code"}
        </Button>
      </div>

      <AceEditor
        mode={language}
        theme="monokai"
        name="editor"
        fontSize={16}
        lineHeight={19}
        showPrintMargin={true}
        showGutter={true}
        highlightActiveLine={true}
        value={code}
        onChange={handleCodeChange}
        setOptions={{
          enableBasicAutocompletion: true,
          enableLiveAutocompletion: true,
          enableSnippets: true,
          showLineNumbers: true,
          tabSize: 2,
        }}
      />

      {/* Submission Result Display */}
      {submissionResult && (
        <div className="mt-4 p-4 border rounded-lg bg-gray-50">
          <div className="flex items-center gap-2 mb-2">
            <h3 className="font-semibold">Submission Result:</h3>
            <span
              className={`font-bold ${getResultColor(submissionResult.result)}`}
            >
              {submissionResult.result}
            </span>
            <span className="text-gray-600">
              Score: {submissionResult.score}/100
            </span>
          </div>

          {submissionResult.message && (
            <p className="text-sm text-gray-700 mb-2">
              {submissionResult.message}
            </p>
          )}

          {submissionResult.time && submissionResult.memory && (
            <div className="text-sm text-gray-600 mb-2">
              Time: {submissionResult.time}ms | Memory:{" "}
              {submissionResult.memory}KB
            </div>
          )}

          {submissionResult.compile_output && (
            <div className="mt-2">
              <h4 className="font-medium text-red-600">Compilation Output:</h4>
              <pre className="text-sm bg-gray-100 p-2 rounded overflow-x-auto">
                {submissionResult.compile_output}
              </pre>
            </div>
          )}

          {submissionResult.std_output && (
            <div className="mt-2">
              <h4 className="font-medium text-green-600">Output:</h4>
              <pre className="text-sm bg-gray-100 p-2 rounded overflow-x-auto">
                {submissionResult.std_output}
              </pre>
            </div>
          )}

          {submissionResult.std_error && (
            <div className="mt-2">
              <h4 className="font-medium text-red-600">Error:</h4>
              <pre className="text-sm bg-gray-100 p-2 rounded overflow-x-auto">
                {submissionResult.std_error}
              </pre>
            </div>
          )}
        </div>
      )}

      {isSubmitting && (
        <div className="mt-4 p-4 border rounded-lg bg-blue-50">
          <div className="flex items-center gap-2">
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
            <span className="text-blue-700">Processing your submission...</span>
          </div>
        </div>
      )}
    </div>
  );
};

export default CodeEditor;

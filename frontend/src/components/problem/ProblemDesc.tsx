import React from "react";

const ProblemDesc = ({ problem }: { problem: ProblemType }) => {
  const test = problem.tests ? problem.tests[0] : null;
  return (
    <div className="w-full h-screen p-6 border-2 rounded-lg shadow-md">
      <h1 className="text-2xl font-bold mb-4">{problem.title}</h1>
      <p className="mb-4">{problem.description}</p>
      <h2 className="text-xl font-semibold mb-2">Example</h2>
      <pre className="p-4 rounded">
        {test?.created_at && (
          <>
            <strong>Input:</strong> {test.input}
            <br />
          </>
        )}
        {test?.created_at && (
          <>
            <strong>Output:</strong> {test.output}
          </>
        )}
      </pre>
    </div>
  );
};

export default ProblemDesc;

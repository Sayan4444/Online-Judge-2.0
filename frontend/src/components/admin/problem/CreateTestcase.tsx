"use client";
import { Input } from "@/components/ui/input";
import { Label } from "@radix-ui/react-label";
import React, { FormEvent } from "react";
import { Button } from "@/components/ui/button";
import { createTestCaseByProblemID } from "@/fetch/testcase";

const CreateTestcase = ({
  problemID,
  token,
}: {
  problemID: string;
  token: string;
}) => {
  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const input = formData.get("input") as string;
    const output = formData.get("output") as string;
    const payload = {
      input: input,
      output: output,
    };

    try {
      const data = await createTestCaseByProblemID(problemID, payload, token);

      // Handle success (e.g., redirect or show a success message)
      console.log("Testcase created successfully", data);
    } catch (error) {
      console.error("Error creating problem:", error);
    }
  };
  return (
    <div className="max-w-2xl mx-auto p-6 bg-white shadow-md rounded-lg">
      <h1 className="text-2xl font-bold mb-4">Create New Testcase</h1>
      <form className="space-y-4" onSubmit={handleSubmit}>
        <div>
          <Label
            htmlFor="input"
            className="block text-sm font-medium text-gray-700"
          >
            Input
          </Label>
          <Input id="input" name="input" required />
        </div>
        <div>
          <Label
            htmlFor="output"
            className="block text-sm font-medium text-gray-700"
          >
            Output
          </Label>
          <Input id="output" name="output" required></Input>
        </div>
        <Button type="submit">Create Testcase</Button>
      </form>
    </div>
  );
};

export default CreateTestcase;

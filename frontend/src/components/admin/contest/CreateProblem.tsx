"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { CreateProblemByContestId } from "@/fetch/problem";
import React, { FormEvent } from "react";

const CreateProblem = ({ id, token }: { id: string; token: string }) => {
  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const title = formData.get("title") as string;
    const description = formData.get("description") as string;
    const payload = {
      title,
      description,
    };

    try {
      const data = await CreateProblemByContestId(id, payload, token);

      // Handle success (e.g., redirect or show a success message)
      console.log("Problem created successfully", data);
    } catch (error) {
      console.error("Error creating problem:", error);
    }
  };
  return (
    <div className="max-w-2xl mx-auto p-6 bg-white shadow-md rounded-lg">
      <h1 className="text-2xl font-bold mb-4">Create New Problem</h1>
      <form className="space-y-4" onSubmit={handleSubmit}>
        <div>
          <Label
            htmlFor="title"
            className="block text-sm font-medium text-gray-700"
          >
            Problem Title
          </Label>
          <Input type="text" id="title" name="title" required />
        </div>
        <div>
          <label
            htmlFor="description"
            className="block text-sm font-medium text-gray-700"
          >
            Description
          </label>
          <textarea
            id="description"
            name="description"
            rows={4}
            required
          ></textarea>
        </div>

        <Button type="submit">Create Problem</Button>
      </form>
    </div>
  );
};

export default CreateProblem;

"use client";
import React, { FormEvent, useState } from "react";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Pencil } from "lucide-react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { updateTestCaseByTestCaseID } from "@/fetch/testcase";

const EditTestcase = ({
  testcaseID,
  token,
}: {
  testcaseID: string;
  token: string;
}) => {
  const [open, setOpen] = useState(false);
  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const input = formData.get("input") as string;
    const output = formData.get("output") as string;
    const payload = {
      input,
      output,
    };

    try {
      const data = await updateTestCaseByTestCaseID(testcaseID, payload, token);

      console.log("Testcase updated successfully", data);
    } catch (error) {
      console.error("Error updating Testcase:", error);
    }
  };
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Pencil className="h-4 w-4" />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Are you absolutely sure?</DialogTitle>
          <DialogDescription>
            This action cannot be undone. This will permanently delete your
            account and remove your data from our servers.
          </DialogDescription>
        </DialogHeader>
        <form className="space-y-4" onSubmit={handleSubmit}>
          <div>
            <Label htmlFor="input">Input</Label>
            <Input type="text" id="input" name="input" />
          </div>
          <div>
            <Label htmlFor="output">Output</Label>
            <Input id="output" name="output"></Input>
          </div>

          <Button type="submit">Update Testcase</Button>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default EditTestcase;

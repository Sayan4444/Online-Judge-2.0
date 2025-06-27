import React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import EditTestcase from "./EditTestcase";
import DeleteTestcase from "./DeleteTestcase";
import { Label } from "@/components/ui/label";

const TestcaseCard = ({
  testcase,
  token,
}: {
  testcase: TestcaseType;
  token: string;
}) => {
  return (
    <Card className="w-full h-auto">
      <CardHeader>
        <Label>Input</Label>
        <CardTitle>{testcase.input}</CardTitle>
        <Label>Output</Label>
        <CardTitle>{testcase.output}</CardTitle>
      </CardHeader>
      <CardContent>
        <p>Created at: {new Date(testcase.created_at).toLocaleString()}</p>
      </CardContent>
      <div className="flex gap-2 justify-end mr-2">
        <EditTestcase testcaseID={testcase.id} token={token} />
        <DeleteTestcase testcaseID={testcase.id} token={token} />
      </div>
    </Card>
  );
};

export default TestcaseCard;

import React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import EditProblem from "./EditProblem";
import DeleteProblem from "./DeleteProblem";

const ProblemCard = ({
  problem,
  token,
}: {
  problem: ProblemType;
  token: string;
}) => {
  return (
    <Card className="w-full max-w-md h-auto cursor-pointer">
      <CardHeader>
        <CardTitle>{problem.title}</CardTitle>
        <CardDescription>{problem.description}</CardDescription>
        <div className="flex gap-2 justify-end mr-2">
          <EditProblem problemID={problem.id} token={token} />
          <DeleteProblem problemID={problem.id} token={token} />
        </div>
      </CardHeader>
      <CardContent>
        <p>Created at: {new Date(problem.created_at).toLocaleString()}</p>
      </CardContent>
      <CardFooter>
        <Button variant="outline">
          <Link href={`problem/${problem.id}`}>View Problem</Link>
        </Button>
      </CardFooter>
    </Card>
  );
};

export default ProblemCard;

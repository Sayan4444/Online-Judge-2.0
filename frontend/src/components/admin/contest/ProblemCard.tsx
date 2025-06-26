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

const ProblemCard = ({ problem }: { problem: ProblemType }) => {
  return (
    <Card className="w-full max-w-md h-auto cursor-pointer">
      <CardHeader>
        <CardTitle>{problem.title}</CardTitle>
        <CardDescription>{problem.description}</CardDescription>
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

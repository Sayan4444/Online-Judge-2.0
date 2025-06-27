"use client";
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
import EditContest from "./EditContest";
import DeleteContest from "./DeleteContest";

const ContestCard = ({
  contest,
  token,
}: {
  contest: ContestType;
  token: string;
}) => {
  return (
    <Card className="w-full max-w-md h-auto cursor-pointer">
      <CardHeader>
        <CardTitle>{contest.name}</CardTitle>
        <CardDescription>{contest.description}</CardDescription>
      </CardHeader>
      <div className="flex gap-2 justify-end mr-2">
        <EditContest contestID={contest.id} token={token} />
        <DeleteContest contestID={contest.id} token={token} />
      </div>
      <CardContent>
        <p>Start time: {new Date(contest.start_time).toUTCString()}</p>
        <p>End time: {new Date(contest.end_time).toUTCString()}</p>
      </CardContent>
      <CardFooter>
        <Button variant="outline">
          <Link href={`contest/${contest.id}`}>View Contest</Link>
        </Button>
      </CardFooter>
    </Card>
  );
};

export default ContestCard;

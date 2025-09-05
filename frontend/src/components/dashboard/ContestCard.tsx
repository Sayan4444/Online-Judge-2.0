import React from "react";
import Link from "next/link";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const ContestCard = ({ contest }: { contest: ContestType }) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{contest.name}</CardTitle>
        <CardDescription>{contest.description}</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-sm mb-4">
          Start Date:{" "}
          {new Date(contest.start_time).toLocaleDateString("en-US", {
            year: "numeric",
            month: "2-digit",
            day: "2-digit",
          })}
        </p>
      </CardContent>
      <CardFooter>
        <Link href={`/contest/${contest.id}`} className=" hover:underline">
          View Contest
        </Link>
      </CardFooter>
    </Card>
  );
};

export default ContestCard;

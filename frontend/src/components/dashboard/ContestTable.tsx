import React from "react";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import Link from "next/link";

const ContestTable = ({ contests }: { contests: Array<ContestType> }) => {
  return (
    <Table>
      <TableCaption>Current contests</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[100px]">Contest Name</TableHead>
          <TableHead>Description</TableHead>
          <TableHead>Start date</TableHead>
          <TableHead>Link</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {contests.map((contest) => (
          <TableRow key={contest.id}>
            <TableCell className="font-medium">{contest.name}</TableCell>
            <TableCell>{contest.description}</TableCell>
            <TableCell>
              {new Date(contest.start_time).toLocaleDateString("en-US", {
                year: "numeric",
                month: "2-digit",
                day: "2-digit",
              })}
            </TableCell>
            <TableCell>
              <Link href={`/contest/${contest.id}`}>View Contest</Link>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};

export default ContestTable;

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

const ProblemTable = ({ problems }: { problems: Array<ProblemType> }) => {
  return (
    <Table>
      <TableCaption>Current problems</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[100px]">Problem Title</TableHead>
          <TableHead>Link to problem</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {problems.map((problem) => (
          <TableRow key={problem.id}>
            <TableCell className="font-medium">{problem.title}</TableCell>
            <TableCell>
              <Link href={`/contest/problem/${problem.id}`}>View Problem</Link>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};

export default ProblemTable;

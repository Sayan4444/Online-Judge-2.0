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
import { RotateCcw } from "lucide-react";

const ProblemTable = ({ problems }: { problems: Array<ProblemType> }) => {
  if (!problems || problems.length === 0) {
    return (
      <div className="rounded-lg border p-6 shadow">
        <div className="flex items-center gap-2 mb-2">
          <RotateCcw className="h-5 w-5 text-blue-500" />
          <h2 className="text-lg font-bold">Problem Statements</h2>
        </div>
        <p className="text-sm text-gray-500">No problems available...</p>
      </div>
    );
  }

  return (
    <div className="rounded-lg border w-full p-6 shadow">
      {/* Header with icon */}
      <div className="flex items-center gap-2 mb-2">
        <RotateCcw className="h-5 w-5 text-blue-500" />
        <h2 className="text-lg font-bold">Problem Statements</h2>
      </div>
      <p className="text-sm text-gray-500 mb-4">
        Choose a problem to start coding
      </p>

      {/* Styled Table */}
      <Table>
        <TableCaption className="sr-only">Problem List</TableCaption>
        <TableHeader>
          <TableRow className="">
            <TableHead className="w-[80px]">Id</TableHead>
            <TableHead>Problem Title</TableHead>
            <TableHead>Link</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {problems.map((problem, index) => (
            <TableRow key={problem.id} className="border-b">
              <TableCell>{index + 1}</TableCell>
              <TableCell className="font-medium">{problem.title}</TableCell>
              <TableCell>
                <Link
                  href={`/contest/problem/${problem.id}`}
                  className="text-blue-600 hover:underline"
                >
                  Solve
                </Link>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
};

export default ProblemTable;

import {
  authOptions,
  CustomSession,
} from "@/app/api/auth/[...nextauth]/options";
import CreateProblem from "@/components/admin/contest/CreateProblem";
import ProblemCard from "@/components/admin/contest/ProblemCard";
import { fetchProblemsByContestId } from "@/fetch/problem";
import { getServerSession } from "next-auth";
import React from "react";

const page = async ({ params }: { params: Promise<{ id: string }> }) => {
  const { id } = await params;
  const session: CustomSession | null = await getServerSession(authOptions);
  const problems = await fetchProblemsByContestId(id, session?.user?.token!);
  return (
    <>
      <h1>Contest ID: {id}</h1>
      <CreateProblem id={id} token={session?.user?.token!} />
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6 mx-4 w-full">
        {problems.map((problem: any) => (
          <ProblemCard key={problem.id} problem={problem} />
        ))}
      </div>
    </>
  );
};

export default page;

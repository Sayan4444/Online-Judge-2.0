import {
  authOptions,
  CustomSession,
} from "@/app/api/auth/[...nextauth]/options";
import Navbar from "@/components/dashboard/Navbar";
import CodeEditor from "@/components/problem/CodeEditor";
import ProblemDesc from "@/components/problem/ProblemDesc";
import { fetchProblemByProblemId } from "@/fetch/problem";
import { getServerSession } from "next-auth";
import React from "react";

const page = async ({ params }: { params: Promise<{ id: string }> }) => {
  const { id } = await params;
  const session: CustomSession | null = await getServerSession(authOptions);
  if (!session || !session.user) {
    return <div>Please log in to access this page.</div>;
  }
  const problem = await fetchProblemByProblemId(id, session.user.token!);

  return (
    <>
      <Navbar user={session?.user} />
      <div className="flex w-full p-4">
        <ProblemDesc problem={problem!} />
        <CodeEditor user={session?.user} problem={problem!} />
      </div>
    </>
  );
};

export default page;

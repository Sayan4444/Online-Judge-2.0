import {
  authOptions,
  CustomSession,
} from "@/app/api/auth/[...nextauth]/options";
import CreateTestcase from "@/components/admin/problem/CreateTestcase";
import TestcaseCard from "@/components/admin/problem/TestcaseCard";
import { fetchTestCasesByProblemId } from "@/fetch/testcase";
import { getServerSession } from "next-auth";
import React from "react";

const page = async ({ params }: { params: Promise<{ id: string }> }) => {
  const { id } = await params;
  const session: CustomSession | null = await getServerSession(authOptions);
  if (!session || !session.user) {
    return <div>Please log in to access this page.</div>;
  }
  const testcases = await fetchTestCasesByProblemId(id, session.user.token!);
  if (!testcases) {
    return <div>No test cases found for this problem.</div>;
  }
  return (
    <>
      <CreateTestcase problemID={id} token={session?.user?.token!} />
      <div className="flex flex-col gap-4 mt-6 mx-auto px-6 w-full">
        {testcases.map((testcase: any) => (
          <TestcaseCard
            key={testcase.id}
            testcase={testcase}
            token={session?.user?.token!}
          />
        ))}
      </div>
    </>
  );
};

export default page;

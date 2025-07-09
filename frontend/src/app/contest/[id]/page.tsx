import { getServerSession } from "next-auth";
import React from "react";
import {
  authOptions,
  CustomSession,
} from "../../api/auth/[...nextauth]/options";
import ProblemTable from "@/components/contest/ProblemTable";
import { fetchProblemsByContestIdUser } from "@/fetch/problem";
import Navbar from "@/components/dashboard/Navbar";

const contest = async ({ params }: { params: Promise<{ id: string }> }) => {
  const { id } = await params;
  const session: CustomSession | null = await getServerSession(authOptions);
  if (!session || !session.user) {
    return <div>Please log in to access this page.</div>;
  }
  const problems = await fetchProblemsByContestIdUser(id, session.user.token!);

  return (
    <>
      <Navbar user={session?.user} />
      <ProblemTable problems={problems} />
    </>
  );
};

export default contest;

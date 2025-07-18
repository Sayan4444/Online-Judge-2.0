import Navbar from "@/components/dashboard/Navbar";
import { getServerSession } from "next-auth";
import React from "react";
import { authOptions, CustomSession } from "../api/auth/[...nextauth]/options";
import { fetchContests } from "@/fetch/contest";
import ContestTable from "@/components/dashboard/ContestTable";

const dashboard = async () => {
  const session: CustomSession | null = await getServerSession(authOptions);
  if (!session || !session.user) {
    return <div>Please log in to access this page.</div>;
  }
  const contests: Array<ContestType> | [] = await fetchContests();

  return (
    <>
      <Navbar user={session?.user} />
      <ContestTable contests={contests} />
    </>
  );
};

export default dashboard;

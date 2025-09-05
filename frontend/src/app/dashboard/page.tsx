import Navbar from "@/components/dashboard/Navbar";
import { getServerSession } from "next-auth";
import React from "react";
import { authOptions, CustomSession } from "../api/auth/[...nextauth]/options";
import { fetchContests } from "@/fetch/contest";
import ContestCard from "@/components/dashboard/ContestCard";

const dashboard = async () => {
  const session: CustomSession | null = await getServerSession(authOptions);
  if (!session || !session.user) {
    return <div>Please log in to access this page.</div>;
  }
  const contests: Array<ContestType> | [] = await fetchContests();
  const liveContests = contests.filter(
    (contest) => new Date(contest.start_time) > new Date()
  );

  return (
    <>
      <Navbar user={session?.user} />
      <div className="text-3xl font-semibold mx-10 mt-6">Live Contests</div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mx-10 mt-6">
        {contests.length > 0 &&
          contests.map((contest) => (
            <ContestCard key={contest.id} contest={contest} />
          ))}
      </div>
    </>
  );
};

export default dashboard;

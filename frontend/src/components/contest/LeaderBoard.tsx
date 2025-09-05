"use client";
import React, { useEffect, useState } from "react";
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { fetchLeaderboard } from "@/fetch/leaderboard";
import { Trophy } from "lucide-react";

const LeaderBoard = ({
  contestId,
  token,
}: {
  contestId: string;
  token: string;
}) => {
  const [leaderboardData, setLeaderboardData] = useState<
    LeaderboardEntryType[]
  >([]);

  const POLL_INTERVAL = 15 * 60 * 1000; // 15 minutes

  useEffect(() => {
    let isMounted = true;

    const fetchData = async () => {
      try {
        const data = await fetchLeaderboard(contestId, token);
        if (isMounted) {
          setLeaderboardData(data);
        }
      } catch (err) {
        if (isMounted) {
          console.error("Failed to fetch leaderboard:", err);
        }
      }
    };

    fetchData();

    const intervalId = setInterval(fetchData, POLL_INTERVAL);

    return () => {
      isMounted = false;
      clearInterval(intervalId);
    };
  }, [contestId, token]);

  if (!leaderboardData || leaderboardData.length === 0) {
    return (
      <div className="rounded-lg border p-6 shadow">
        <div className="flex items-center gap-2 mb-2">
          <Trophy className="h-5 w-5 text-blue-500" />
          <h2 className="text-lg font-bold">Live Rankings</h2>
        </div>
        <p className="text-sm text-gray-500">No participants yet...</p>
      </div>
    );
  }

  return (
    <div className="rounded-lg border w-full p-6 shadow">
      <div className="flex items-center gap-2 mb-2">
        <Trophy className="h-5 w-5 text-blue-500" />
        <h2 className="text-lg font-bold">Live Rankings</h2>
      </div>
      <p className="text-sm text-gray-500 mb-4">Top performers this contest</p>

      <Table>
        <TableCaption className="sr-only">Leaderboard</TableCaption>
        <TableHeader>
          <TableRow className="">
            <TableHead className="w-[60px]">#</TableHead>
            <TableHead>Username</TableHead>
            <TableHead>Score</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {leaderboardData.map((entry, index) => (
            <TableRow key={entry.user_id} className="border-b">
              <TableCell className="font-medium">{index + 1}</TableCell>
              <TableCell>{entry.username}</TableCell>
              <TableCell>{entry.total_score}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
};

export default LeaderBoard;

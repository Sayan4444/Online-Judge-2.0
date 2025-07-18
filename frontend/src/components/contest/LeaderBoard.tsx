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

  // No data fallback
  if (!leaderboardData || leaderboardData.length === 0) {
    return <>No participants yet...</>;
  }
  return (
    <Table>
      <TableCaption>Leaderboard</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Rank</TableHead>
          <TableHead>Username</TableHead>
          <TableHead>Score</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {leaderboardData.map((entry, index) => (
          <TableRow key={entry.user_id}>
            <TableCell className="font-medium">{index + 1}</TableCell>
            <TableCell>{entry.username}</TableCell>
            <TableCell>{entry.total_score}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
};

export default LeaderBoard;

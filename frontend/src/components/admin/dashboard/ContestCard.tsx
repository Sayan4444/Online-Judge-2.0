"use client";
import React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import axios from "axios";
import { ADMIN_URL } from "@/lib/apiEndpoints";
import {
  CustomSession,
  CustomUser,
} from "@/app/api/auth/[...nextauth]/options";

const ContestCard = (user: CustomUser | null) => {
  const fetchContests = async () => {
    const data = axios.get(`${ADMIN_URL}/contests`, {
      headers: {
        Authorization: user?.token,
      },
    });
    console.log("Contests data:", data);
  };

  return (
    <Card className="w-full max-w-md h-auto cursor-pointer">
      <CardHeader>
        <CardTitle>Contest Title</CardTitle>
        <CardDescription>Short description of the contest.</CardDescription>
      </CardHeader>
      <CardContent>
        <p>Additional details about the contest can go here.</p>
      </CardContent>
      <CardFooter>
        <Button onClick={fetchContests} variant="outline">
          View Contest
        </Button>
      </CardFooter>
    </Card>
  );
};

export default ContestCard;

import LogoutButton from "@/components/auth/LogoutButton";
import Navbar from "@/components/dashboard/Navbar";
import { getServerSession } from "next-auth";
import React from "react";
import { authOptions } from "../api/auth/[...nextauth]/options";

const dashboard = async () => {
  const session = await getServerSession(authOptions);
  if (!session || !session.user) {
    return <div>Please log in to access this page.</div>;
  }

  return (
    <>
      <Navbar user={session?.user} />
    </>
  );
};

export default dashboard;

import {
  authOptions,
  CustomSession,
} from "@/app/api/auth/[...nextauth]/options";
import ContestCard from "@/components/admin/dashboard/ContestCard";
import LogoutButton from "@/components/auth/LogoutButton";
import { getServerSession } from "next-auth";
import React from "react";

async function page() {
  const session: CustomSession | null = await getServerSession(authOptions);

  return (
    <div>
      <p>admin dash</p>
      <LogoutButton />
    </div>
  );
}

export default page;

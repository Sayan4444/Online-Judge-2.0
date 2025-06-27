import { getServerSession } from "next-auth";
import React from "react";
import { authOptions, CustomSession } from "../api/auth/[...nextauth]/options";
import UpdateProfile from "@/components/profile/UpdateProfile";

const page = async () => {
  const session: CustomSession | null = await getServerSession(authOptions);
  if (!session || !session.user) {
    return <div>Please log in to access this page.</div>;
  }
  const user = session.user;
  return (
    <>
      <div className="flex flex-col items-center justify-center h-screen">
        <h1 className="text-4xl font-bold mb-4">Welcome! {user.name}</h1>
        <p className="text-lg">This is the profile page</p>
        <UpdateProfile user={user} />
      </div>
    </>
  );
};

export default page;

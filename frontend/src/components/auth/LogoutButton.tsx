"use client";
import React from "react";
import { Button } from "@/components/ui/button";
import { signOut } from "next-auth/react";

const LogoutButton = () => {
  const handleLogout = () => {
    signOut({
      callbackUrl: "/",
      redirect: true,
    });
  };
  return <Button onClick={handleLogout}>logout</Button>;
};

export default LogoutButton;

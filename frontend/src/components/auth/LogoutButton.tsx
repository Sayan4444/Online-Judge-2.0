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
  return (
    <Button className="rounded-none" variant="secondary" onClick={handleLogout}>
      logout
    </Button>
  );
};

export default LogoutButton;

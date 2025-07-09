"use client";
import React from "react";
import { Button } from "@/components/ui/button";
import { signIn } from "next-auth/react";

const LoginButton = () => {
  const handleLogin = () => {
    signIn("google", {
      callbackUrl: "/dashboard",
      redirect: true,
    });
  };
  return (
    <div>
      <Button onClick={handleLogin}>Login</Button>
    </div>
  );
};

export default LoginButton;

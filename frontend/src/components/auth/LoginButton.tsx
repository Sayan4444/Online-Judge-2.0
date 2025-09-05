"use client";
import React from "react";
import { Button } from "@/components/ui/button";
import { signIn } from "next-auth/react";
import {ArrowRight} from "lucide-react"

const LoginButton = () => {
  const handleLogin = () => {
    signIn("google", {
      callbackUrl: "/dashboard",
      redirect: true,
    });
  };
  return (
    <div>
      <Button variant="outline" onClick={handleLogin}>
        Login
        <ArrowRight className="ml-1" />
      </Button>
    </div>
  );
};

export default LoginButton;

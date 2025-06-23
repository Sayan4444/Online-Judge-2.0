"use client";
import React, { FormEvent, useRef } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { signIn } from "next-auth/react";

const LoginCard = () => {
  const email = useRef<HTMLInputElement | null>(null);
  const password = useRef<HTMLInputElement | null>(null);
  const handleLogin = (e: FormEvent) => {
    e.preventDefault();
    console.log("Login attempt with", {
      email: email,
      password: password,
    });

    signIn("admin-login", {
      redirect: true,
      callbackUrl: "/admin/dashboard",
      email: email?.current?.value,
      password: password?.current?.value,
    });
  };
  return (
    <Card className="w-full max-w-md h-auto">
      <CardHeader>
        <CardTitle>Admin Login</CardTitle>
        <CardDescription>GLUG Admin Dashboard Login</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleLogin}>
          <div className="flex flex-col gap-6 mb-4">
            <div className="grid gap-2">
              <Input
                id="email"
                type="email"
                placeholder="Admin Email"
                ref={email}
                required
              />
            </div>
            <div className="grid gap-2">
              <Input
                id="password"
                type="password"
                placeholder="Password"
                ref={password}
                required
              />
            </div>
          </div>
          <Button type="submit">Submit</Button>
        </form>
      </CardContent>
    </Card>
  );
};

export default LoginCard;

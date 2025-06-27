"use client";
import React, { useRef } from "react";
import { Label } from "../ui/label";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { updateUserProfileName } from "@/fetch/profile";
import { CustomUser } from "@/app/api/auth/[...nextauth]/options";

const UpdateProfile = ({ user }: { user: CustomUser }) => {
  const name = useRef<HTMLInputElement>(null);
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (name.current) {
      console.log(
        "Updating profile for user:",
        user,
        "with name:",
        name.current.value
      );
      const data = await updateUserProfileName({
        name: name.current.value,
        token: user?.token!,
      });
      console.log("Profile updated successfully:", data);
    }
  };
  return (
    <>
      <div className="flex flex-col">
        <form onSubmit={handleSubmit}>
          <Label>Update Profile</Label>
          <Input type="text" ref={name} />

          <Button type="submit" className="mt-4">
            Update
          </Button>
        </form>
      </div>
    </>
  );
};

export default UpdateProfile;

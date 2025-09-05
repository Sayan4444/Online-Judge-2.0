"use client";
import Image from "next/image";
import React from "react";
import ojImage from "../../../public/oj.png";
import { Avatar, AvatarFallback, AvatarImage } from "@radix-ui/react-avatar";
import { CustomUser } from "@/app/api/auth/[...nextauth]/options";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@radix-ui/react-popover";
import LogoutButton from "../auth/LogoutButton";
import { Button } from "../ui/button";
import Link from "next/link";
import { ThemeToggle } from "../theme/ThemeToggle";

const Navbar = ({ user }: { user: CustomUser }) => {
  return (
    <>
      <div className="flex items-center justify-between p-4 border-b-2">
        <Image src={ojImage} alt="Logo" width={50} height={50} />
        <div className="flex items-center gap-4">
          <ThemeToggle />
          <Popover>
            <PopoverTrigger asChild>
              <Avatar className="cursor-pointer">
                <AvatarImage
                  src={user?.image!}
                  alt={user?.name!}
                  width={50}
                  height={50}
                  className="rounded-full"
                />
                <AvatarFallback>CN</AvatarFallback>
              </Avatar>
            </PopoverTrigger>
            <PopoverContent className="p-2 flex flex-col items-start">
              <LogoutButton />
              <Button className="rounded-none mt-2" variant="secondary">
                <Link href="/profile">Profile</Link>
              </Button>
            </PopoverContent>
          </Popover>
        </div>
      </div>
    </>
  );
};

export default Navbar;

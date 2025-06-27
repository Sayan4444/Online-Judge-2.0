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

const Navbar = ({ user }: { user: CustomUser }) => {
  return (
    <>
      <div className="flex items-center justify-between bg-gray-800 p-4">
        <Image src={ojImage} alt="Logo" width={50} height={50} />
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
          <PopoverContent className="p-2 flex flex-col items-start bg-gray-800">
            <LogoutButton />
            <Button className="rounded-none mt-2" variant="secondary">
              <Link href="/profile">Profile</Link>
            </Button>
          </PopoverContent>
        </Popover>
      </div>
    </>
  );
};

export default Navbar;

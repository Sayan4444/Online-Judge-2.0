"use client";
import { CustomUser } from "@/app/api/auth/[...nextauth]/options";
import { Button } from "@/components/ui/button";
import React, { FormEvent, useState } from "react";
import { Calendar } from "@/components/ui/calendar";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { ChevronDownIcon } from "lucide-react";
import axios from "axios";
import { ADMIN_URL } from "@/lib/apiEndpoints";

const CreateContest = ({ user }: { user: CustomUser }) => {
  const [open, setOpen] = useState(false);
  const [date, setDate] = useState<Date | undefined>(undefined);
  const [startTime, setStartTime] = useState<string>("8:30:00");
  const [endTime, setEndTime] = useState<string>("10:30:00");

  const handleCreateContest = async (e: FormEvent) => {
    e.preventDefault();
    const formData = new FormData(e.target as HTMLFormElement);
    const contestTitle = formData.get("contestTitle") as string;
    const contestDescription = formData.get("contestDescription") as string;

    if (
      !contestTitle ||
      !contestDescription ||
      !date ||
      !startTime ||
      !endTime
    ) {
      alert("Please fill in all fields.");
      return;
    }

    const contestStartTime = new Date(date);
    const contestEndTime = new Date(date);
    const [startHours, startMinutes] = startTime.split(":").map(Number);
    const [endHours, endMinutes] = endTime.split(":").map(Number);

    const payload = {
      name: contestTitle,
      description: contestDescription,
      start_time: new Date(contestStartTime.setHours(startHours, startMinutes)),
      end_time: new Date(contestEndTime.setHours(endHours, endMinutes)),
    };

    const data = await axios.post(`${ADMIN_URL}/create-contest`, payload, {
      headers: {
        Authorization: `Bearer ${user.token}`,
      },
    });

    console.log("Contest created successfully:", data);
  };

  return (
    <>
      <div className="flex flex-col items-center justify-center min-h-screen">
        <h1 className="text-2xl font-bold mb-4">Create Contest</h1>
        <p className="mb-4">Welcome, {user.name || "Admin"}!</p>
        <form className="w-full max-w-md" onSubmit={handleCreateContest}>
          <div className="mb-4">
            <Label
              htmlFor="contestTitle"
              className="block text-sm font-medium text-gray-700"
            >
              Contest Title
            </Label>
            <Input type="text" id="contestTitle" name="contestTitle" required />
          </div>
          <div className="mb-4">
            <Label
              htmlFor="contestDescription"
              className="block text-sm font-medium text-gray-700"
            >
              Contest Description
            </Label>
            <textarea
              id="contestDescription"
              name="contestDescription"
              required
              rows={4}
              className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:border-blue-500 focus:ring-blue-500"
            ></textarea>
            <div className="flex flex-col gap-3">
              <Label htmlFor="date-picker" className="px-1">
                Date
              </Label>
              <Popover open={open} onOpenChange={setOpen}>
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    id="date-picker"
                    className="w-32 justify-between font-normal"
                  >
                    {date ? date.toLocaleDateString() : "Select date"}
                    <ChevronDownIcon />
                  </Button>
                </PopoverTrigger>
                <PopoverContent
                  className="w-auto overflow-hidden p-0"
                  align="start"
                >
                  <Calendar
                    mode="single"
                    selected={date}
                    captionLayout="dropdown"
                    onSelect={(date) => {
                      setDate(date);
                      setOpen(false);
                    }}
                  />
                </PopoverContent>
              </Popover>
            </div>
            <div className="flex flex-col gap-3">
              <Label htmlFor="time-picker" className="px-1">
                Start Time
              </Label>
              <Input
                type="time"
                id="time-picker"
                step="1"
                defaultValue="8:30:00"
                onChange={(e) => setStartTime(e.target.value)}
                className="bg-background appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
              />
            </div>
            <div className="flex flex-col gap-3">
              <Label htmlFor="time-picker" className="px-1">
                End Time
              </Label>
              <Input
                type="time"
                id="time-picker"
                step="1"
                defaultValue="10:30:00"
                onChange={(e) => setEndTime(e.target.value)}
                className="bg-background appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
              />
            </div>
          </div>
          <Button type="submit">Create Contest</Button>
        </form>
      </div>
    </>
  );
};

export default CreateContest;

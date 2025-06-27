"use client";
import React, { useState } from "react";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Pencil } from "lucide-react";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Calendar } from "@/components/ui/calendar";
import { updateContest } from "@/fetch/contest";

const EditContest = ({
  contestID,
  token,
}: {
  contestID: string;
  token: string;
}) => {
  const [open, setOpen] = useState(false);
  const [date, setDate] = useState<Date | undefined>(undefined);
  const [startTime, setStartTime] = useState<string>("8:30:00");
  const [endTime, setEndTime] = useState<string>("10:30:00");

  const handleUpdateContest = async (e: React.FormEvent) => {
    e.preventDefault();
    const formData = new FormData(e.target as HTMLFormElement);
    const contestTitle = formData.get("contestTitle") as string;
    const contestDescription = formData.get("contestDescription") as string;
    if (!date) {
      console.log("Please select a date.");
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
    console.log("Payload for updating contest:", payload);
    const data = await updateContest(token, payload, contestID);
    console.log("Contest updated successfully:", data);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Pencil className="h-4 w-4" />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <form onSubmit={handleUpdateContest}>
          <DialogHeader>
            <DialogTitle>Are you absolutely sure?</DialogTitle>
            <DialogDescription>
              This action cannot be undone. This will permanently delete your
              account and remove your data from our servers.
            </DialogDescription>
          </DialogHeader>

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
              <Calendar
                mode="single"
                selected={date}
                captionLayout="dropdown"
                onSelect={(date) => {
                  setDate(date);
                }}
              />
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
          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">Cancel</Button>
            </DialogClose>
            <Button type="submit">Save changes</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default EditContest;

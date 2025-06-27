import {
  authOptions,
  CustomSession,
} from "@/app/api/auth/[...nextauth]/options";
import ContestCard from "@/components/admin/dashboard/ContestCard";
import CreateContest from "@/components/admin/dashboard/CreateContest";
import LogoutButton from "@/components/auth/LogoutButton";
import { fetchContests } from "@/fetch/contest";
import { getServerSession } from "next-auth";

async function page() {
  const session: CustomSession | null = await getServerSession(authOptions);
  const contests: Array<ContestType> | [] = await fetchContests(
    session?.user?.token!
  );
  return (
    <div>
      <div className="flex justify-end p-6">
        <LogoutButton />
      </div>
      <CreateContest user={session?.user!} />
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {contests?.map((contest: any) => (
          <ContestCard
            key={contest.id}
            contest={contest}
            token={session?.user?.token!}
          />
        ))}
      </div>
    </div>
  );
}

export default page;

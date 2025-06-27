import { API_URL } from "@/lib/apiEndpoints";
import axios from "axios";
import { toast } from "sonner";

export const updateUserProfileName = async ({
  name,
  token,
}: {
  name: string;
  token: string;
}) => {
  try {
    const data = await axios.put(
      `${API_URL}/profile`,
      { username: name },
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    toast.success("Profile name updated successfully");
    return data.data;
  } catch (error) {
    console.error("Error updating user profile name:", error);
    return null;
  }
};

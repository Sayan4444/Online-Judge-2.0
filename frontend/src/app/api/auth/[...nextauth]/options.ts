import axios from "axios";
import { Account, AuthOptions, ISODateString } from "next-auth";
import { JWT } from "next-auth/jwt";
import GoogleProvider from "next-auth/providers/google";
import { ADMIN_LOGIN_URL, LOGIN_URL } from "@/lib/apiEndpoints";
import CredentialsProvider from "next-auth/providers/credentials";

// This interface extends the default session object to include a user object
export interface CustomSession {
  user?: CustomUser;
  expires: ISODateString;
}

// This interface extends the default session object to include a user object to store the token without using context or cookies
export interface CustomUser {
  id?: string;
  name?: string | null | undefined;
  email?: string | null | undefined;
  image?: string | null | undefined;
  provider?: string;
  token?: string;
}

export const authOptions: AuthOptions = {
  secret: process.env.NEXTAUTH_SECRET,
  pages: {
    signIn: "/",
  },
  callbacks: {
    async signIn({
      user,
      account,
    }: {
      user: CustomUser;
      account: Account | null;
    }) {
      try {
        if (account?.provider === "admin-login") {
          user.token = user?.token;
          console.log("Admin login successful, user token:", user.token);
          return true;
        }
        const payload = {
          email: user.email!,
          username: user.name!,
          oauth_id: account?.providerAccountId!,
          provider: account?.provider!,
          image: user?.image,
        };
        console.log("Payload for login:", payload);
        const { data } = await axios.post(LOGIN_URL, payload);
        console.log("Response from login:", data);

        user.id = data?.user?.id?.toString();
        user.token = data?.token;
        return true;
      } catch (error) {
        console.error("Sign-in error:", error);
        return false;
      }
    },
    async session({
      session,
      user,
      token,
    }: {
      session: CustomSession;
      user: CustomUser;
      token: JWT;
    }) {
      session.user = token.user as CustomUser;
      return session;
    },
    async jwt({ token, user }) {
      if (user) {
        const customUser = user as CustomUser;
        token.user = {
          ...customUser,
          token:
            typeof customUser.token === "string"
              ? customUser.token
              : (customUser as any).token?.token,
        };
      }
      return token;
    },
  },
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID as string,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET as string,
    }),

    CredentialsProvider({
      id: "admin-login",
      name: "Admin Login",
      credentials: {
        email: { label: "email", type: "email" },
        password: { label: "Password", type: "password" },
      },
      async authorize(credentials: Record<string, string> | undefined) {
        const res = await fetch(ADMIN_LOGIN_URL, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            email: credentials?.email,
            password: credentials?.password,
          }),
        });

        const token = await res.json();
        console.log("Admin login response:", token);

        if (res.ok && token) {
          return {
            id: "1100",
            name: "Admin",
            email: "admin",
            token: token as string,
          };
        }

        return null;
      },
    }),
  ],
};

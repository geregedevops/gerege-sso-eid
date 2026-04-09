import NextAuth from "next-auth";
import { prisma } from "./db";

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [
    {
      id: "gerege-sso",
      name: "e-ID Mongolia",
      type: "oidc",
      issuer: process.env.NEXT_PUBLIC_SSO_URL || "https://sso.gerege.mn",
      clientId: process.env.GEREGE_SSO_CLIENT!,
      clientSecret: process.env.GEREGE_SSO_SECRET!,
      authorization: {
        params: { scope: "openid profile" },
      },
    },
  ],
  callbacks: {
    async jwt({ token, account, profile }) {
      if (account) {
        token.accessToken = account.access_token;
      }
      if (profile) {
        token.sub = profile.sub as string;
        token.name = (profile as any).name || "";
        token.certSerial = (profile as any).cert_serial || "";
        token.givenName = (profile as any).given_name || "";
        token.familyName = (profile as any).family_name || "";

        // Upsert user record
        try {
          await prisma.user.upsert({
            where: { sub: profile.sub as string },
            update: {
              name: (profile as any).name || "",
              givenName: (profile as any).given_name || "",
              familyName: (profile as any).family_name || "",
              certSerial: (profile as any).cert_serial || "",
            },
            create: {
              sub: profile.sub as string,
              name: (profile as any).name || "",
              givenName: (profile as any).given_name || "",
              familyName: (profile as any).family_name || "",
              certSerial: (profile as any).cert_serial || "",
            },
          });
        } catch (e) {
          console.error("User upsert error:", e);
        }
      }
      return token;
    },
    async session({ session, token }) {
      if (session.user) {
        (session.user as any).sub = token.sub;
        (session.user as any).certSerial = token.certSerial;
        (session.user as any).accessToken = token.accessToken;
      }
      return session;
    },
  },
  pages: {
    signIn: "/auth/login",
  },
  trustHost: true,
});

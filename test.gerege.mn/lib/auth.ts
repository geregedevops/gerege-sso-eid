import NextAuth from "next-auth";

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [
    {
      id: "gerege-sso",
      name: "GeregeID",
      type: "oidc",
      issuer: process.env.NEXT_PUBLIC_SSO_URL || "https://sso.gerege.mn",
      clientId: process.env.EID_CLIENT_ID!,
      clientSecret: process.env.EID_CLIENT_SECRET!,
      authorization: {
        params: { scope: "openid profile pos social payment" },
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
        token.tenantId = (profile as any).tenant_id || "";
        token.tenantRole = (profile as any).tenant_role || "";
      }
      return token;
    },
    async session({ session, token }) {
      if (session.user) {
        (session.user as any).sub = token.sub;
        (session.user as any).accessToken = token.accessToken;
        (session.user as any).tenantId = token.tenantId;
        (session.user as any).tenantRole = token.tenantRole;
      }
      return session;
    },
  },
  pages: {
    signIn: "/auth/login",
  },
  trustHost: true,
});

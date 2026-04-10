import NextAuth from "next-auth";

export const { handlers, signIn, signOut, auth } = NextAuth({
  providers: [
    {
      id: "gerege-sso",
      name: "e-ID Mongolia",
      type: "oidc",
      issuer: process.env.NEXT_PUBLIC_SSO_URL || "https://sso.gerege.mn",
      clientId: process.env.GEREGE_SSO_CLIENT!,
      clientSecret: process.env.GEREGE_SSO_SECRET!,
      authorization: { params: { scope: "openid profile" } },
    },
  ],
  callbacks: {
    async jwt({ token, profile }) {
      if (profile) {
        token.sub = profile.sub as string;
        token.certSerial = (profile as any).cert_serial || "";
        token.givenName = (profile as any).given_name || "";
        token.familyName = (profile as any).family_name || "";
        token.name = (profile as any).name || "";
      }
      return token;
    },
    async session({ session, token }) {
      if (session.user) {
        (session.user as any).sub = token.sub;
        (session.user as any).certSerial = token.certSerial;
      }
      return session;
    },
  },
  pages: {
    signIn: "/auth/login",
  },
  trustHost: true,
});

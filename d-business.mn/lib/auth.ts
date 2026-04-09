import NextAuth from "next-auth";
import { query, genId } from "./db";

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

        // Upsert user
        try {
          const sub = profile.sub as string;
          const name = (profile as any).name || "";
          const givenName = (profile as any).given_name || "";
          const familyName = (profile as any).family_name || "";
          const certSerial = (profile as any).cert_serial || "";

          const existing = await query(
            `SELECT id FROM dbiz_users WHERE sub = $1`,
            [sub]
          );

          if (existing.length > 0) {
            await query(
              `UPDATE dbiz_users SET name=$1, "givenName"=$2, "familyName"=$3, "certSerial"=$4, "updatedAt"=now() WHERE sub=$5`,
              [name, givenName, familyName, certSerial, sub]
            );
          } else {
            await query(
              `INSERT INTO dbiz_users (id, sub, name, "givenName", "familyName", "certSerial") VALUES ($1,$2,$3,$4,$5,$6)`,
              [genId(), sub, name, givenName, familyName, certSerial]
            );
          }
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

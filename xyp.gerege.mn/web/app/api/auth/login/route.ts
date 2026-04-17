import { NextResponse } from "next/server";
import {
  SESSION_COOKIE,
  SESSION_MAX_AGE,
  createSessionValue,
  constantTimeEqual,
} from "@/lib/session";

export async function POST(req: Request) {
  const form = await req.formData();
  const password = (form.get("password") as string) || "";
  const expected = process.env.VERIFY_WEB_ADMIN_KEY || "";

  const ok = !!expected && constantTimeEqual(password, expected);
  const target = ok ? "/admin" : "/auth/login?error=1";

  const res = new NextResponse(null, {
    status: 303,
    headers: { Location: target },
  });
  if (ok) {
    const value = await createSessionValue();
    res.cookies.set(SESSION_COOKIE, value, {
      httpOnly: true,
      secure: true,
      sameSite: "lax",
      path: "/",
      maxAge: SESSION_MAX_AGE,
    });
  }
  return res;
}

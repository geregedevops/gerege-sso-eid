import { cookies } from "next/headers";
import { NextResponse } from "next/server";
import { SESSION_COOKIE } from "@/lib/session";

export async function POST(req: Request) {
  cookies().delete(SESSION_COOKIE);
  const url = new URL("/auth/login", req.url);
  return NextResponse.redirect(url, { status: 303 });
}

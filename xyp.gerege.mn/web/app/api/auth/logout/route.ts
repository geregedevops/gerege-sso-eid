import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";
import { SESSION_COOKIE } from "@/lib/session";

export async function POST(req: NextRequest) {
  cookies().delete(SESSION_COOKIE);
  const host = req.headers.get("x-forwarded-host") || req.headers.get("host") || "";
  const proto = req.headers.get("x-forwarded-proto") || "https";
  const url = new URL("/auth/login", `${proto}://${host}`);
  return NextResponse.redirect(url, { status: 303 });
}

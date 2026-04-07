import { NextResponse } from "next/server";
import { auth } from "@/lib/auth";

const SANDBOX_URL = process.env.NEXT_PUBLIC_SANDBOX_URL || "https://sandbox.gerege.mn";

export async function POST(req: Request) {
  const session = await auth();
  const accessToken = (session?.user as any)?.accessToken;

  if (!accessToken) {
    return NextResponse.json({ error: "unauthorized" }, { status: 401 });
  }

  const { method, path, body } = await req.json();

  try {
    const start = Date.now();
    const res = await fetch(`${SANDBOX_URL}${path}`, {
      method: method || "GET",
      headers: {
        Authorization: `Bearer ${accessToken}`,
        "Content-Type": "application/json",
        "X-Sandbox": "true",
      },
      body: body ? JSON.stringify(body) : undefined,
    });

    const duration = Date.now() - start;
    let responseBody;
    try {
      responseBody = await res.json();
    } catch {
      responseBody = { message: await res.text() };
    }

    return NextResponse.json({
      status: res.status,
      statusText: res.statusText,
      body: responseBody,
      duration,
    });
  } catch (e: any) {
    return NextResponse.json({
      status: 502,
      statusText: "Bad Gateway",
      body: { error: `Sandbox unreachable: ${e.message}` },
      duration: 0,
    });
  }
}

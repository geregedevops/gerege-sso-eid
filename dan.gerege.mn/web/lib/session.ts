export const SESSION_COOKIE = "dan_admin_session";
export const SESSION_MAX_AGE = 60 * 60 * 8; // 8 hours

function b64urlEncode(bytes: ArrayBuffer): string {
  const arr = new Uint8Array(bytes);
  let bin = "";
  for (let i = 0; i < arr.length; i++) bin += String.fromCharCode(arr[i]);
  return btoa(bin).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

function b64urlDecode(s: string): Uint8Array {
  const pad = (4 - (s.length % 4)) % 4;
  const b64 = s.replace(/-/g, "+").replace(/_/g, "/") + "=".repeat(pad);
  const bin = atob(b64);
  const bytes = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) bytes[i] = bin.charCodeAt(i);
  return bytes;
}

async function hmacKey(): Promise<CryptoKey> {
  const secret = process.env.DAN_ADMIN_KEY || "";
  return crypto.subtle.importKey(
    "raw",
    new TextEncoder().encode(secret),
    { name: "HMAC", hash: "SHA-256" },
    false,
    ["sign", "verify"],
  );
}

export async function createSessionValue(): Promise<string> {
  const exp = Math.floor(Date.now() / 1000) + SESSION_MAX_AGE;
  const payload = `admin:${exp}`;
  const key = await hmacKey();
  const sig = await crypto.subtle.sign("HMAC", key, new TextEncoder().encode(payload));
  return `${exp}.${b64urlEncode(sig)}`;
}

export async function verifySessionValue(value: string | undefined): Promise<boolean> {
  if (!value) return false;
  const dot = value.indexOf(".");
  if (dot < 0) return false;
  const expStr = value.slice(0, dot);
  const sig = value.slice(dot + 1);
  const exp = parseInt(expStr, 10);
  if (!exp || exp < Math.floor(Date.now() / 1000)) return false;
  try {
    const key = await hmacKey();
    return await crypto.subtle.verify(
      "HMAC",
      key,
      b64urlDecode(sig),
      new TextEncoder().encode(`admin:${exp}`),
    );
  } catch {
    return false;
  }
}

export function constantTimeEqual(a: string, b: string): boolean {
  if (a.length !== b.length) return false;
  let diff = 0;
  for (let i = 0; i < a.length; i++) diff |= a.charCodeAt(i) ^ b.charCodeAt(i);
  return diff === 0;
}

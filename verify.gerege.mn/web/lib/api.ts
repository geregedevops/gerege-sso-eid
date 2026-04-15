const VERIFY_API_URL = process.env.VERIFY_API_URL || "http://verify-api:8446";
const VERIFY_ADMIN_KEY = process.env.VERIFY_WEB_ADMIN_KEY || "";

async function verifyFetch(path: string, options: RequestInit = {}) {
  const res = await fetch(`${VERIFY_API_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${VERIFY_ADMIN_KEY}`,
      ...options.headers,
    },
    cache: "no-store",
  });
  if (!res.ok) {
    const body = await res.text();
    throw new Error(`Verify API error ${res.status}: ${body}`);
  }
  return res.json();
}

export async function listClients() {
  const data = await verifyFetch("/api/clients");
  return data.clients || [];
}

export async function createClient(body: {
  name: string;
  contact_email: string;
}) {
  return verifyFetch("/api/clients", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function deactivateClient(id: string) {
  return verifyFetch(`/api/clients/${id}`, {
    method: "DELETE",
  });
}

export async function getUsage(params?: {
  client_id?: string;
  from?: string;
  to?: string;
}) {
  const qs = new URLSearchParams();
  if (params?.client_id) qs.set("client_id", params.client_id);
  if (params?.from) qs.set("from", params.from);
  if (params?.to) qs.set("to", params.to);
  const query = qs.toString() ? `?${qs.toString()}` : "";
  const data = await verifyFetch(`/api/usage${query}`);
  return data.usage || [];
}

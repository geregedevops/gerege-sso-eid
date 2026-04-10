const SSO_API_URL = process.env.SSO_API_URL || "http://sso:8443";
const DAN_ADMIN_KEY = process.env.DAN_ADMIN_KEY || "";

async function ssoFetch(path: string, options: RequestInit = {}) {
  const res = await fetch(`${SSO_API_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${DAN_ADMIN_KEY}`,
      ...options.headers,
    },
    cache: "no-store",
  });
  if (!res.ok) {
    const body = await res.text();
    throw new Error(`SSO API error ${res.status}: ${body}`);
  }
  return res.json();
}

export async function listDANClients() {
  return ssoFetch("/api/dan/clients");
}

export async function createDANClient(data: {
  name: string;
  callback_urls: string[];
}) {
  return ssoFetch("/api/dan/clients", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function deactivateDANClient(id: string) {
  return ssoFetch(`/api/dan/clients/${id}`, {
    method: "DELETE",
  });
}

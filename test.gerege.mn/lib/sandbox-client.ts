export class SandboxClient {
  private baseURL: string;
  private token: string;

  constructor(token: string) {
    this.baseURL = process.env.NEXT_PUBLIC_SANDBOX_URL || "https://sandbox.gerege.mn";
    this.token = token;
  }

  async request(method: string, path: string, body?: object) {
    const start = Date.now();
    const res = await fetch(`${this.baseURL}${path}`, {
      method,
      headers: {
        Authorization: `Bearer ${this.token}`,
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
      responseBody = null;
    }

    return {
      status: res.status,
      statusText: res.statusText,
      headers: Object.fromEntries(res.headers),
      body: responseBody,
      duration,
    };
  }
}

const API_URL = process.env.API_URL || "https://api.gerege.mn";

export async function initiateSign(
  accessToken: string,
  signerReg: string,
  documentName: string,
  documentBase64: string
) {
  const res = await fetch(`${API_URL}/v1/sign/request`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${accessToken}`,
    },
    body: JSON.stringify({
      signer_reg: signerReg,
      document_name: documentName,
      document: documentBase64,
    }),
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`Sign request failed: ${res.status} ${text}`);
  }

  return res.json() as Promise<{
    session_id: string;
    verification_code: string;
  }>;
}

export async function getSignStatus(accessToken: string, sessionId: string) {
  const res = await fetch(`${API_URL}/v1/sign/${sessionId}/status`, {
    headers: { Authorization: `Bearer ${accessToken}` },
  });

  if (!res.ok) {
    throw new Error(`Status check failed: ${res.status}`);
  }

  return res.json() as Promise<{
    status: string;
    signer_name?: string;
    cert_serial?: string;
    result_url?: string;
  }>;
}

export async function getSignResult(
  accessToken: string,
  sessionId: string
): Promise<{ buffer: Buffer; filename: string }> {
  const res = await fetch(`${API_URL}/v1/sign/${sessionId}/result`, {
    headers: { Authorization: `Bearer ${accessToken}` },
  });

  if (!res.ok) {
    throw new Error(`Result fetch failed: ${res.status}`);
  }

  const disposition = res.headers.get("content-disposition") || "";
  const filenameMatch = disposition.match(/filename="?([^"]+)"?/);
  const filename = filenameMatch ? filenameMatch[1] : "signed.pdf";

  const buffer = Buffer.from(await res.arrayBuffer());
  return { buffer, filename };
}

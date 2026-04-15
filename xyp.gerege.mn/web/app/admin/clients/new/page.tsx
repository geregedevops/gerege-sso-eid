import { redirect } from "next/navigation";
import { cookies } from "next/headers";
import { createClient } from "@/lib/api";

export default function NewClientPage({
  searchParams,
}: {
  searchParams: { created?: string };
}) {
  async function handleCreate(formData: FormData) {
    "use server";

    const name = formData.get("name") as string;
    const contactEmail = (formData.get("contact_email") as string) || "";

    if (!name) {
      throw new Error("name is required");
    }

    const result = await createClient({ name, contact_email: contactEmail });

    // Store result in a cookie (shown once, then cleared)
    cookies().set("new_client", JSON.stringify({
      id: result.client.id,
      secret: result.client_secret,
      name: result.client.name,
    }), { maxAge: 60, httpOnly: true, path: "/admin/clients/new" });

    redirect("/admin/clients/new?created=1");
  }

  // Read from cookie (maxAge=60, auto-expires)
  let newClient: { id: string; secret: string; name: string } | null = null;
  if (searchParams.created) {
    const raw = cookies().get("new_client")?.value;
    if (raw) {
      try {
        newClient = JSON.parse(raw);
      } catch {}
    }
  }

  if (newClient) {
    return (
      <div className="max-w-lg mx-auto px-6 py-10">
        <h1 className="text-2xl font-bold text-white mb-2">Client үүслээ</h1>
        <p className="text-sm text-red-400 mb-6 font-medium">
          Доорх мэдээллийг одоо хадгалаарай. Client secret дахин харагдахгүй!
        </p>

        <div className="bg-surface border border-white/[0.06] rounded-2xl p-6 space-y-4">
          <div>
            <p className="text-xs text-slate-400 mb-1">Нэр</p>
            <p className="text-white font-medium">{newClient.name}</p>
          </div>
          <div>
            <p className="text-xs text-slate-400 mb-1">Client ID</p>
            <p className="font-mono text-sm text-primary break-all select-all">{newClient.id}</p>
          </div>
          <div>
            <p className="text-xs text-slate-400 mb-1">Client Secret</p>
            <p className="font-mono text-sm text-yellow-400 break-all select-all">{newClient.secret}</p>
          </div>
        </div>

        <div className="mt-4 p-4 bg-yellow-500/10 border border-yellow-500/20 rounded-xl">
          <p className="text-xs text-yellow-400">
            API дуудлага хийхдээ HTTP Basic Auth ашиглана:
            <br />
            <code className="font-mono">Authorization: Basic base64(client_id:client_secret)</code>
          </p>
        </div>

        <div className="mt-6">
          <a href="/admin/clients" className="text-sm text-slate-400 hover:text-white">
            &larr; Client жагсаалт руу буцах
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-lg mx-auto px-6 py-10">
      <h1 className="text-2xl font-bold text-white mb-2">Шинэ Client</h1>
      <p className="text-sm text-slate-400 mb-8">
        API ашиглах шинэ client бүртгэх. Client ID болон Secret үүсгэнэ.
      </p>

      <form action={handleCreate} className="space-y-5">
        <div>
          <label className="block text-sm font-medium text-slate-300 mb-1.5">
            Нэр
          </label>
          <input
            name="name"
            required
            placeholder="Миний апп"
            className="w-full px-4 py-2.5 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-slate-300 mb-1.5">
            Холбоо барих имэйл (заавал биш)
          </label>
          <input
            name="contact_email"
            type="email"
            placeholder="admin@company.mn"
            className="w-full px-4 py-2.5 bg-white/[0.05] border border-white/10 rounded-xl text-white text-sm focus:outline-none focus:border-primary"
          />
        </div>

        <button
          type="submit"
          className="w-full py-3 bg-primary hover:bg-primary-light text-white font-bold rounded-xl transition-all"
        >
          Client үүсгэх
        </button>
      </form>

      <div className="mt-6">
        <a href="/admin/clients" className="text-sm text-slate-400 hover:text-white">
          &larr; Client жагсаалт руу буцах
        </a>
      </div>
    </div>
  );
}

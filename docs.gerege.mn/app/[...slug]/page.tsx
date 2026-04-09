import { notFound } from "next/navigation";
import { getPage, getAllSlugs } from "@/content/docs";
import { Sidebar } from "@/components/sidebar";
import { Markdown } from "@/components/markdown";

export function generateStaticParams() {
  return getAllSlugs().map((slug) => ({ slug }));
}

export default async function DocPage({ params }: { params: Promise<{ slug: string[] }> }) {
  const { slug } = await params;
  const page = getPage(slug);
  if (!page) notFound();

  return (
    <div className="flex">
      <Sidebar currentSlug={page.slug} />
      <main className="flex-1 max-w-3xl px-8 py-8">
        <Markdown content={page.content} />
      </main>
    </div>
  );
}

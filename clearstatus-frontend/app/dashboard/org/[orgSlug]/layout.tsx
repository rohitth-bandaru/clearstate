"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { api } from "@/lib/api";
import type { Organization } from "@/lib/types";
import { Sidebar } from "@/components/layout/sidebar";

export default function OrgLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const params = useParams();
  const router = useRouter();
  const orgSlug = params.orgSlug as string;
  const [org, setOrg] = useState<Organization | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api<Organization[] | null>("/api/orgs")
      .then((list) => {
        const arr = Array.isArray(list) ? list : [];
        const found = arr.find((o) => o.slug === orgSlug);
        if (found) setOrg(found);
        else router.replace("/dashboard");
      })
      .catch(() => router.replace("/dashboard"))
      .finally(() => setLoading(false));
  }, [orgSlug, router]);

  if (loading || !org) return <div className="p-8 text-zinc-500">Loading...</div>;

  return (
    <div className="flex min-h-[calc(100vh-3.5rem)]">
      <Sidebar orgSlug={orgSlug} />
      <div className="flex-1 overflow-auto">{children}</div>
    </div>
  );
}

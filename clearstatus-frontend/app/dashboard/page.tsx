"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import type { Organization } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export default function DashboardPage() {
  const [orgs, setOrgs] = useState<Organization[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api<Organization[] | null>("/api/orgs")
      .then((data) => setOrgs(Array.isArray(data) ? data : []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="p-8 text-zinc-500">Loading...</div>;
  if (error) return <div className="p-8 text-red-600">{error}</div>;

  return (
    <div className="p-8 max-w-2xl">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Organizations</CardTitle>
            <CardDescription>Select an organization or create one</CardDescription>
          </div>
          <Button asChild>
            <Link href="/dashboard/new-org">New organization</Link>
          </Button>
        </CardHeader>
        <CardContent>
          {!orgs?.length ? (
            <p className="text-zinc-500 text-sm">No organizations yet. Create one to get started.</p>
          ) : (
            <ul className="space-y-2">
              {(orgs ?? []).map((org) => (
                <li key={org.id}>
                  <Link
                    href={`/dashboard/org/${org.slug}/services`}
                    className="block p-3 rounded-lg border border-zinc-200 hover:bg-zinc-50 transition-colors"
                  >
                    <span className="font-medium">{org.name}</span>
                    <span className="text-zinc-500 text-sm ml-2">/{org.slug}</span>
                  </Link>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

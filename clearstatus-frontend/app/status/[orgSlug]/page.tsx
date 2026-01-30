"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { useStatusPolling } from "@/lib/polling";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

const statusVariant: Record<string, "default" | "secondary" | "destructive" | "outline"> = {
  operational: "default",
  degraded: "secondary",
  partial_outage: "outline",
  major_outage: "destructive",
};

export default function PublicStatusPage() {
  const params = useParams();
  const orgSlug = params.orgSlug as string;
  const { data, error, loading } = useStatusPolling(orgSlug);

  if (loading && !data) {
    return (
      <div className="min-h-screen bg-zinc-50 flex items-center justify-center">
        <p className="text-zinc-500">Loading...</p>
      </div>
    );
  }
  if (error && !data) {
    return (
      <div className="min-h-screen bg-zinc-50 flex items-center justify-center">
        <p className="text-red-600">{error}</p>
      </div>
    );
  }
  if (!data) return null;

  const { org, services = [], incidents = [], timeline = [], summary } = data;

  return (
    <div className="min-h-screen bg-zinc-50">
      <header className="border-b border-zinc-200 bg-white">
        <div className="max-w-3xl mx-auto px-6 py-6">
          <h1 className="text-xl font-semibold text-zinc-900">{org.name}</h1>
          <p className="text-sm text-zinc-500 mt-1">System status</p>
        </div>
      </header>
      <main className="max-w-3xl mx-auto px-6 py-8">
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="text-lg">{summary}</CardTitle>
            <CardDescription>
              Status is updated every 2 seconds. This page refreshes only when data changes.
            </CardDescription>
          </CardHeader>
        </Card>

        <section className="mb-8">
          <h2 className="text-sm font-medium text-zinc-500 mb-3">Services</h2>
          <Card>
            <CardContent className="p-0">
              {!services?.length ? (
                <p className="p-6 text-zinc-500 text-sm">No services configured.</p>
              ) : (
                <ul className="divide-y divide-zinc-100">
                  {(services ?? []).map((svc) => (
                    <li key={svc.id} className="flex items-center justify-between px-6 py-4">
                      <div>
                        <span className="font-medium text-zinc-900">{svc.name}</span>
                        {svc.description && (
                          <p className="text-sm text-zinc-500 mt-0.5">{svc.description}</p>
                        )}
                      </div>
                      <Badge variant={statusVariant[svc.status] || "outline"}>
                        {svc.status.replace("_", " ")}
                      </Badge>
                    </li>
                  ))}
                </ul>
              )}
            </CardContent>
          </Card>
        </section>

        {(incidents ?? []).filter((i) => i.status !== "resolved").length > 0 && (
          <section className="mb-8">
            <h2 className="text-sm font-medium text-zinc-500 mb-3">Active incidents</h2>
            <Card>
              <CardContent className="p-0">
                <ul className="divide-y divide-zinc-100">
                  {(incidents ?? [])
                    .filter((i) => i.status !== "resolved")
                    .map((inc) => (
                      <li key={inc.id} className="px-6 py-4">
                        <div className="flex items-start justify-between gap-4">
                          <div>
                            <span className="font-medium text-zinc-900">{inc.title}</span>
                            <span className="ml-2 text-xs text-zinc-500">({inc.type})</span>
                            <Badge variant="secondary" className="ml-2">
                              {inc.status}
                            </Badge>
                          </div>
                          <span className="text-xs text-zinc-500 whitespace-nowrap">
                            {new Date(inc.created_at).toLocaleString()}
                          </span>
                        </div>
                        {(inc.updates ?? []).length > 0 && (
                          <div className="mt-3 pl-4 border-l-2 border-zinc-200 space-y-2">
                            {(inc.updates ?? []).map((u, i) => (
                              <div key={i}>
                                <p className="text-sm text-zinc-700">{u.message}</p>
                                <p className="text-xs text-zinc-500">
                                  {new Date(u.created_at).toLocaleString()}
                                </p>
                              </div>
                            ))}
                          </div>
                        )}
                      </li>
                    ))}
                </ul>
              </CardContent>
            </Card>
          </section>
        )}

        {(timeline ?? []).length > 0 && (
          <section>
            <h2 className="text-sm font-medium text-zinc-500 mb-3">Timeline</h2>
            <Card>
              <CardContent className="p-0">
                <ul className="divide-y divide-zinc-100">
                  {(timeline ?? []).slice(0, 20).map((entry, i) => (
                    <li key={i} className="px-6 py-3 flex gap-4">
                      <span className="text-xs text-zinc-500 whitespace-nowrap">
                        {new Date(entry.created_at).toLocaleString()}
                      </span>
                      <div className="text-sm text-zinc-700">
                        {entry.type === "incident" && entry.title && (
                          <>Incident: {entry.title}</>
                        )}
                        {entry.type === "incident_update" && entry.message && (
                          <>{entry.message}</>
                        )}
                        {entry.type === "incident_resolved" && entry.title && (
                          <>Resolved: {entry.title}</>
                        )}
                        {entry.type === "service_status" && entry.service_name && entry.status && (
                          <>{entry.service_name}: {entry.status}</>
                        )}
                      </div>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          </section>
        )}

        <Separator className="my-8" />
        <p className="text-center text-sm text-zinc-500">
          Powered by <Link href="/" className="text-zinc-900 underline">ClearStatus</Link>
        </p>
      </main>
    </div>
  );
}

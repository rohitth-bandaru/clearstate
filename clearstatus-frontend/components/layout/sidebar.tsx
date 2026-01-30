"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

const nav = [
  { href: "services", label: "Services" },
  { href: "incidents", label: "Incidents" },
  { href: "teams", label: "Teams" },
];

export function Sidebar({ orgSlug }: { orgSlug: string }) {
  const pathname = usePathname();
  const base = `/dashboard/org/${orgSlug}`;
  return (
    <aside className="w-56 border-r border-zinc-200 bg-zinc-50/50 p-4 flex flex-col">
      <Link href={base} className="text-sm font-medium text-zinc-900 mb-4">
        {orgSlug}
      </Link>
      <nav className="flex flex-col gap-0.5">
        {nav.map(({ href, label }) => {
          const path = `${base}/${href}`;
          const active = pathname === path || pathname.startsWith(path + "/");
          return (
            <Link
              key={href}
              href={path}
              className={cn(
                "px-3 py-2 rounded-md text-sm font-medium transition-colors",
                active
                  ? "bg-zinc-200 text-zinc-900"
                  : "text-zinc-600 hover:bg-zinc-100 hover:text-zinc-900"
              )}
            >
              {label}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}

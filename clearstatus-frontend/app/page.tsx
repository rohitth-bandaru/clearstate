"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { isAuthenticated } from "@/lib/auth";

export default function HomePage() {
  const router = useRouter();
  useEffect(() => {
    if (isAuthenticated()) {
      router.replace("/dashboard");
    } else {
      router.replace("/login");
    }
  }, [router]);
  return (
    <div className="min-h-screen flex items-center justify-center bg-zinc-50">
      <p className="text-zinc-500">Redirecting...</p>
    </div>
  );
}

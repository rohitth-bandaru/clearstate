"use client";

import { useState, useEffect, useRef } from "react";
import equal from "fast-deep-equal";
import { publicFetch } from "./api";
import type { PublicStatusResponse } from "./types";

// Poll interval in ms; set NEXT_PUBLIC_POLL_INTERVAL_MS (default 2000 when unset)
const POLL_INTERVAL_MS = (() => {
  const n = Number(process.env.NEXT_PUBLIC_POLL_INTERVAL_MS);
  return Number.isNaN(n) || n <= 0 ? 2000 : n;
})();

export function useStatusPolling(orgSlug: string | null): {
  data: PublicStatusResponse | null;
  error: string | null;
  loading: boolean;
} {
  const [data, setData] = useState<PublicStatusResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const previousRef = useRef<PublicStatusResponse | null>(null);

  useEffect(() => {
    if (!orgSlug) {
      setData(null);
      setError(null);
      setLoading(false);
      return;
    }
    const slug = orgSlug;
    let cancelled = false;

    async function fetchStatus() {
      try {
        const next = await publicFetch<PublicStatusResponse>(
          `/api/public/orgs/${encodeURIComponent(slug)}/status`
        );
        if (cancelled) return;
        if (!equal(previousRef.current, next)) {
          previousRef.current = next;
          setData(next);
        }
        setError(null);
      } catch (e) {
        if (!cancelled) {
          setError(e instanceof Error ? e.message : "Failed to load status");
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    }

    fetchStatus();
    const interval = setInterval(fetchStatus, POLL_INTERVAL_MS);
    return () => {
      cancelled = true;
      clearInterval(interval);
    };
  }, [orgSlug]);

  return { data, error, loading };
}

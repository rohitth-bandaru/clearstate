"use client";

import { setToken, clearToken, api } from "./api";
import type { User } from "./types";

const USER_KEY = "clearstatus_user";

export function getStoredUser(): User | null {
  if (typeof window === "undefined") return null;
  const raw = localStorage.getItem(USER_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as User;
  } catch {
    return null;
  }
}

export function setStoredUser(user: User | null) {
  if (typeof window === "undefined") return;
  if (user) {
    localStorage.setItem(USER_KEY, JSON.stringify(user));
  } else {
    localStorage.removeItem(USER_KEY);
  }
}

export async function signInWithGoogle(idToken: string): Promise<{ token: string; user: User }> {
  const data = await api<{ token: string; user: User }>("/api/auth/google", {
    method: "POST",
    body: JSON.stringify({ id_token: idToken }),
  });
  setToken(data.token);
  setStoredUser(data.user);
  return data;
}

export function signOut() {
  clearToken();
  setStoredUser(null);
  if (typeof window !== "undefined") {
    window.location.href = "/login";
  }
}

export function isAuthenticated(): boolean {
  if (typeof window === "undefined") return false;
  return !!localStorage.getItem("clearstatus_token");
}

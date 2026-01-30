"use client";

import { useRouter } from "next/navigation";
import { GoogleLogin } from "@react-oauth/google";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { signInWithGoogle } from "@/lib/auth";

const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID || "";

export default function LoginPage() {
  const router = useRouter();

  return (
    <div className="min-h-screen flex items-center justify-center bg-zinc-50">
      <Card className="w-full max-w-sm border-zinc-200 shadow-sm">
        <CardHeader>
          <CardTitle className="text-xl font-semibold">ClearStatus</CardTitle>
          <CardDescription>Sign in to manage your status pages</CardDescription>
        </CardHeader>
        <CardContent>
          {!clientId ? (
            <p className="text-sm text-amber-700 bg-amber-50 border border-amber-200 rounded-md p-3">
              Google sign-in is not configured. Add <code className="bg-amber-100 px-1 rounded">NEXT_PUBLIC_GOOGLE_CLIENT_ID</code> to the root <code className="bg-amber-100 px-1 rounded">.env</code> and rebuild: <code className="bg-amber-100 px-1 rounded text-xs">docker-compose up --build</code>
            </p>
          ) : (
          <GoogleLogin
            onSuccess={async (res) => {
              if (res.credential) {
                try {
                  await signInWithGoogle(res.credential);
                  router.push("/dashboard");
                  router.refresh();
                } catch (e) {
                  console.error(e);
                }
              }
            }}
            onError={() => console.error("Google login failed")}
            useOneTap={false}
            theme="outline"
            size="large"
            text="signin_with"
            shape="rectangular"
          />
          )}
        </CardContent>
      </Card>
    </div>
  );
}

"use client";

import { useRouter } from "next/navigation";
import { GoogleLogin } from "@react-oauth/google";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { signInWithGoogle } from "@/lib/auth";

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
        </CardContent>
      </Card>
    </div>
  );
}

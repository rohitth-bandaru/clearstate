import { redirect } from "next/navigation";

export default async function OrgHomePage({
  params,
}: {
  params: Promise<{ orgSlug: string }>;
}) {
  const { orgSlug } = await params;
  redirect(`/dashboard/org/${orgSlug}/services`);
}

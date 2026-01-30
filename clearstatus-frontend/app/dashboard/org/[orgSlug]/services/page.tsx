"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { api } from "@/lib/api";
import type { Organization, Service, Team } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ServiceStatus, type ServiceStatusType } from "@/lib/types";

const statusVariant: Record<string, "default" | "secondary" | "destructive" | "outline"> = {
  operational: "default",
  degraded: "secondary",
  partial_outage: "outline",
  major_outage: "destructive",
};

export default function ServicesPage() {
  const params = useParams();
  const orgSlug = params.orgSlug as string;
  const [orgId, setOrgId] = useState<string | null>(null);
  const [services, setServices] = useState<Service[]>([]);
  const [teams, setTeams] = useState<Team[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api<Organization[] | null>("/api/orgs")
      .then((list) => {
        const arr = Array.isArray(list) ? list : [];
        const org = arr.find((o) => o.slug === orgSlug);
        if (org) {
          setOrgId(org.id);
          return Promise.all([
            api<Service[] | null>(`/api/orgs/${org.id}/services`),
            api<Team[] | null>(`/api/orgs/${org.id}/teams`),
          ]);
        }
        throw new Error("Org not found");
      })
      .then(([servicesData, teamsData]: [Service[] | null, Team[] | null]) => {
        setServices(Array.isArray(servicesData) ? servicesData : []);
        setTeams(Array.isArray(teamsData) ? teamsData : []);
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [orgSlug]);

  const refetch = () => {
    if (!orgId) return;
    Promise.all([
      api<Service[] | null>(`/api/orgs/${orgId}/services`),
      api<Team[] | null>(`/api/orgs/${orgId}/teams`),
    ]).then(([servicesData, teamsData]: [Service[] | null, Team[] | null]) => {
      setServices(Array.isArray(servicesData) ? servicesData : []);
      if (Array.isArray(teamsData)) setTeams(teamsData);
    }).catch(console.error);
  };

  if (loading) return <div className="p-8 text-zinc-500">Loading...</div>;
  if (error) return <div className="p-8 text-red-600">{error}</div>;
  if (!orgId) return null;

  return (
    <div className="p-8">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Services</CardTitle>
            <CardDescription>Manage services and their status</CardDescription>
          </div>
          <CreateServiceDialog orgId={orgId} teams={teams} onCreated={refetch} />
        </CardHeader>
        <CardContent>
          {!services?.length ? (
            <p className="text-zinc-500 text-sm">No services yet. Create one to get started.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Team</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[100px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(services ?? []).map((svc) => (
                  <TableRow key={svc.id}>
                    <TableCell className="font-medium">{svc.name}</TableCell>
                    <TableCell className="text-zinc-500">
                      {teams.find((t) => t.id === svc.team_id)?.name ?? svc.team_id}
                    </TableCell>
                    <TableCell>
                      <Badge variant={statusVariant[svc.status] || "outline"}>
                        {svc.status.replace("_", " ")}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <EditServiceDialog service={svc} teams={teams} onUpdated={refetch} />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
      <p className="mt-4 text-sm text-zinc-500">
        Public status page:{" "}
        <Link href={`/status/${orgSlug}`} className="text-zinc-900 underline">
          /status/{orgSlug}
        </Link>
      </p>
    </div>
  );
}

function CreateServiceDialog({
  orgId,
  teams,
  onCreated,
}: {
  orgId: string;
  teams: Team[];
  onCreated: () => void;
}) {
  const [open, setOpen] = useState(false);
  const [teamId, setTeamId] = useState("");
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [status, setStatus] = useState<string>(ServiceStatus.operational);
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  const submit = async () => {
    if (!teamId) {
      setErr("Please select a team");
      return;
    }
    setLoading(true);
    setErr(null);
    try {
      await api(`/api/orgs/${orgId}/services`, {
        method: "POST",
        body: JSON.stringify({
          team_id: teamId,
          name: name || "New Service",
          description: description || "",
          status,
        }),
      });
      setOpen(false);
      setTeamId("");
      setName("");
      setDescription("");
      setStatus(ServiceStatus.operational);
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>Add service</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New service</DialogTitle>
          <DialogDescription>Add a service to a team. Slug is generated from the name.</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-4">
          <div>
            <Label>Team</Label>
            <Select value={teamId} onValueChange={setTeamId} required>
              <SelectTrigger>
                <SelectValue placeholder="Select a team" />
              </SelectTrigger>
              <SelectContent>
                {(teams ?? []).map((t) => (
                  <SelectItem key={t.id} value={t.id}>
                    {t.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div>
            <Label>Name</Label>
            <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Website" />
          </div>
          <div>
            <Label>Description</Label>
            <Input value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>
          <div>
            <Label>Status</Label>
            <Select value={status} onValueChange={setStatus}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {Object.values(ServiceStatus).map((s) => (
                  <SelectItem key={s} value={s}>
                    {s.replace("_", " ")}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          {err && <p className="text-sm text-red-600">{err}</p>}
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button onClick={submit} disabled={loading}>
            {loading ? "Creating..." : "Create"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function EditServiceDialog({
  service,
  teams,
  onUpdated,
}: {
  service: Service;
  teams: Team[];
  onUpdated: () => void;
}) {
  const [open, setOpen] = useState(false);
  const [teamId, setTeamId] = useState(service.team_id);
  const [status, setStatus] = useState(service.status);
  const [version, setVersion] = useState(service.version);
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const params = useParams();
  const orgSlug = params.orgSlug as string;
  const [orgId, setOrgId] = useState<string | null>(null);
  useEffect(() => {
    api<Organization[]>("/api/orgs").then((list) => {
      const org = list.find((o) => o.slug === orgSlug);
      if (org) setOrgId(org.id);
    });
  }, [orgSlug]);
  useEffect(() => {
    if (open) {
      setTeamId(service.team_id);
      setStatus(service.status);
      setVersion(service.version);
    }
  }, [open, service.team_id, service.status, service.version]);

  const submit = async () => {
    if (!orgId) return;
    setLoading(true);
    setErr(null);
    try {
      const updated = await api<Service>(`/api/orgs/${orgId}/services/${service.id}`, {
        method: "PUT",
        body: JSON.stringify({ team_id: teamId, status, version }),
      });
      setVersion(updated.version);
      onUpdated();
      setOpen(false);
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Failed. Refetch and retry.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          Edit
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Update status</DialogTitle>
          <DialogDescription>{service.name}</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-4">
          <div>
            <Label>Team</Label>
            <Select value={teamId} onValueChange={setTeamId}>
              <SelectTrigger>
                <SelectValue placeholder="Select team" />
              </SelectTrigger>
              <SelectContent>
                {(teams ?? []).map((t) => (
                  <SelectItem key={t.id} value={t.id}>
                    {t.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div>
            <Label>Status</Label>
            <Select value={status} onValueChange={(v) => setStatus(v as ServiceStatusType)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {Object.values(ServiceStatus).map((s) => (
                  <SelectItem key={s} value={s}>
                    {s.replace("_", " ")}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          {err && <p className="text-sm text-red-600">{err}</p>}
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button onClick={submit} disabled={loading}>
            {loading ? "Saving..." : "Save"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

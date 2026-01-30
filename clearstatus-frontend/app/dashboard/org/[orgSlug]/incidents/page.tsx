"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { api } from "@/lib/api";
import type { Organization, Incident } from "@/lib/types";
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
import { IncidentStatus, type IncidentStatusType } from "@/lib/types";

export default function IncidentsPage() {
  const params = useParams();
  const orgSlug = params.orgSlug as string;
  const [orgId, setOrgId] = useState<string | null>(null);
  const [incidents, setIncidents] = useState<Incident[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api<Organization[] | null>("/api/orgs")
      .then((list) => {
        const arr = Array.isArray(list) ? list : [];
        const org = arr.find((o) => o.slug === orgSlug);
        if (org) {
          setOrgId(org.id);
          return api<Incident[] | null>(`/api/orgs/${org.id}/incidents`);
        }
        throw new Error("Org not found");
      })
      .then((data) => setIncidents(Array.isArray(data) ? data : []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, [orgSlug]);

  const refetch = () => {
    if (!orgId) return;
    api<Incident[] | null>(`/api/orgs/${orgId}/incidents`)
      .then((data) => setIncidents(Array.isArray(data) ? data : []))
      .catch(console.error);
  };

  if (loading) return <div className="p-8 text-zinc-500">Loading...</div>;
  if (error) return <div className="p-8 text-red-600">{error}</div>;
  if (!orgId) return null;

  return (
    <div className="p-8">
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Incidents</CardTitle>
            <CardDescription>Create and manage incidents and maintenance</CardDescription>
          </div>
          <CreateIncidentDialog orgId={orgId} onCreated={refetch} />
        </CardHeader>
        <CardContent>
          {!incidents?.length ? (
            <p className="text-zinc-500 text-sm">No incidents yet.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Title</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="w-[100px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(incidents ?? []).map((inc) => (
                  <TableRow key={inc.id}>
                    <TableCell className="font-medium">{inc.title}</TableCell>
                    <TableCell>{inc.type}</TableCell>
                    <TableCell>
                      <Badge variant={inc.status === "resolved" ? "default" : "secondary"}>
                        {inc.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-zinc-500 text-sm">
                      {new Date(inc.created_at).toLocaleDateString()}
                    </TableCell>
                    <TableCell>
                      <AddUpdateDialog orgId={orgId} incident={inc} onUpdated={refetch} />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

function CreateIncidentDialog({
  orgId,
  onCreated,
}: {
  orgId: string;
  onCreated: () => void;
}) {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [type, setType] = useState<"incident" | "maintenance">("incident");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const [services, setServices] = useState<{ id: string; name: string }[]>([]);
  const [serviceIds, setServiceIds] = useState<string[]>([]);

  useEffect(() => {
    if (open) {
      api<{ id: string; name: string }[] | null>(`/api/orgs/${orgId}/services`).then((list) =>
        setServices(Array.isArray(list) ? list : [])
      );
    }
  }, [open, orgId]);

  const submit = async () => {
    setLoading(true);
    setErr(null);
    try {
      await api(`/api/orgs/${orgId}/incidents`, {
        method: "POST",
        body: JSON.stringify({
          title: title || "New incident",
          type,
          status: "investigating",
          service_ids: serviceIds,
          message: message || "Incident created.",
        }),
      });
      setOpen(false);
      setTitle("");
      setMessage("");
      setServiceIds([]);
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
        <Button>New incident</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New incident</DialogTitle>
          <DialogDescription>Create an incident or scheduled maintenance.</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-4">
          <div>
            <Label>Title</Label>
            <Input value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Outage" />
          </div>
          <div>
            <Label>Type</Label>
            <Select value={type} onValueChange={(v) => setType(v as "incident" | "maintenance")}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="incident">Incident</SelectItem>
                <SelectItem value="maintenance">Maintenance</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div>
            <Label>Initial message</Label>
            <Input
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="We are investigating..."
            />
          </div>
          <div>
            <Label>Affected services (optional)</Label>
            <div className="flex flex-wrap gap-2 mt-1">
              {(services ?? []).map((s) => (
                <Button
                  key={s.id}
                  type="button"
                  variant={serviceIds.includes(s.id) ? "default" : "outline"}
                  size="sm"
                  onClick={() =>
                    setServiceIds((prev) =>
                      prev.includes(s.id) ? prev.filter((id) => id !== s.id) : [...prev, s.id]
                    )
                  }
                >
                  {s.name}
                </Button>
              ))}
            </div>
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

function AddUpdateDialog({
  orgId,
  incident,
  onUpdated,
}: {
  orgId: string;
  incident: Incident;
  onUpdated: () => void;
}) {
  const [open, setOpen] = useState(false);
  const [message, setMessage] = useState("");
  const [status, setStatus] = useState(incident.status);
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  const submit = async () => {
    setLoading(true);
    setErr(null);
    try {
      await api(`/api/orgs/${orgId}/incidents/${incident.id}/updates`, {
        method: "POST",
        body: JSON.stringify({ message: message || "Update", status }),
      });
      setOpen(false);
      setMessage("");
      onUpdated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : "Failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline" size="sm">
          Add update
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add update</DialogTitle>
          <DialogDescription>{incident.title}</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-4">
          <div>
            <Label>Message</Label>
            <Input value={message} onChange={(e) => setMessage(e.target.value)} />
          </div>
          <div>
            <Label>Status</Label>
            <Select value={status} onValueChange={(v) => setStatus(v as IncidentStatusType)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {Object.values(IncidentStatus).map((s) => (
                  <SelectItem key={s} value={s}>
                    {s}
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

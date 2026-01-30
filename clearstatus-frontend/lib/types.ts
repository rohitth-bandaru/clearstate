// DTOs matching backend

export type User = {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
};

export type Organization = {
  id: string;
  name: string;
  slug: string;
  created_at: string;
  updated_at: string;
};

export type Team = {
  id: string;
  org_id: string;
  name: string;
  created_at: string;
  updated_at: string;
};

export const ServiceStatus = {
  operational: "operational",
  degraded: "degraded",
  partial_outage: "partial_outage",
  major_outage: "major_outage",
} as const;
export type ServiceStatusType = (typeof ServiceStatus)[keyof typeof ServiceStatus];

export type Service = {
  id: string;
  org_id: string;
  team_id: string;
  name: string;
  slug: string;
  description: string;
  status: ServiceStatusType;
  sort_order: number;
  version: number;
  created_at: string;
  updated_at: string;
};

export const IncidentStatus = {
  investigating: "investigating",
  identified: "identified",
  monitoring: "monitoring",
  resolved: "resolved",
} as const;
export type IncidentStatusType = (typeof IncidentStatus)[keyof typeof IncidentStatus];

export type Incident = {
  id: string;
  org_id: string;
  title: string;
  status: IncidentStatusType;
  type: "incident" | "maintenance";
  version: number;
  created_at: string;
  updated_at: string;
  resolved_at?: string;
};

export type IncidentWithDetails = Incident & {
  service_ids: string[];
  updates: IncidentUpdate[];
};

export type IncidentUpdate = {
  id: string;
  incident_id: string;
  message: string;
  status: IncidentStatusType;
  created_at: string;
};

// Public status (for polling)
export type PublicStatusResponse = {
  org: { id: string; name: string; slug: string };
  services: PublicService[];
  incidents: PublicIncident[];
  timeline: PublicTimelineEntry[];
  summary: string;
};

export type PublicService = {
  id: string;
  name: string;
  slug: string;
  description: string;
  status: ServiceStatusType;
};

export type PublicIncident = {
  id: string;
  title: string;
  status: IncidentStatusType;
  type: "incident" | "maintenance";
  created_at: string;
  resolved_at?: string;
  service_ids: string[];
  updates: PublicIncidentUpdate[];
};

export type PublicIncidentUpdate = {
  message: string;
  status: IncidentStatusType;
  created_at: string;
};

export type PublicTimelineEntry = {
  type: string;
  service_id?: string;
  service_name?: string;
  status?: string;
  incident_id?: string;
  title?: string;
  message?: string;
  created_at: string;
};

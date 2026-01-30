# Deploy ClearStatus Backend to AWS

This guide covers deploying the Go API to AWS securely, with MySQL on RDS and the app on ECS Fargate behind an ALB.

## Prerequisites

- AWS CLI configured
- Docker (for building image)
- Domain (optional; for HTTPS)

## 1. Database (RDS MySQL)

1. Create a DB subnet group (private subnets).
2. Create a security group `rds-sg` (inbound: 3306 from app security group only).
3. Create RDS MySQL 8 instance:
   - Engine: MySQL 8.0
   - Template: Dev/Test or Production
   - DB instance identifier: `clearstatus-db`
   - Master username/password: set and store in Secrets Manager
   - Instance in VPC, private subnets, no public access
   - Security group: `rds-sg`
4. Note the endpoint and port. DSN format: `user:password@tcp(endpoint:3306)/clearstatus?parseTime=true&multiStatements=true`.
5. Create database: `CREATE DATABASE clearstatus;` (run via bastion or Lambda if no direct access).

## 2. Secrets

Store sensitive config in AWS Secrets Manager or Systems Manager Parameter Store (SSM):

- `clearstatus/jwt-secret` — random string for JWT signing.
- `clearstatus/db-dsn` — full MySQL DSN (or separate host, user, password and build DSN in app).
- `clearstatus/google-client-id` — Google OAuth client ID (same as frontend).

ECS task definition will reference these so the container receives env vars (e.g. via secrets injection).

## 3. ECR Repository

```bash
aws ecr create-repository --repository-name clearstatus-backend
```

Build and push:

```bash
docker build -t clearstatus-backend .
docker tag clearstatus-backend:latest <account-id>.dkr.ecr.<region>.amazonaws.com/clearstatus-backend:latest
aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <account-id>.dkr.ecr.<region>.amazonaws.com
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/clearstatus-backend:latest
```

## 4. ECS Cluster and Task Definition

1. Create ECS cluster: `clearstatus-cluster` (Fargate).
2. Create task execution role with policies: `AmazonECSTaskExecutionRolePolicy`, access to Secrets Manager/SSM for the secrets above.
3. Create task definition:
   - Family: `clearstatus-backend`
   - CPU/Memory: 0.25 vCPU, 512 MB (adjust for load).
   - Container: image = ECR URI, port 8080, env:
     - `PORT=8080`
   - Secrets: map Secrets Manager or SSM keys to env vars `JWT_SECRET`, `DB_DSN`, `GOOGLE_CLIENT_ID`.
   - Logging: CloudWatch log group.
   - Optional: set `MIGRATIONS_DIR` if migrations are at a different path in the image.

Ensure the task runs in private subnets and uses a security group that allows only outbound traffic (and ALB inbound on the listener port).

## 5. Application Load Balancer (ALB)

1. Create ALB in the same VPC (public subnets).
2. Create target group (port 8080, protocol HTTP, VPC, Fargate IP).
3. Create listener (HTTPS 443) with ACM certificate; add HTTP 80 listener that redirects to 443.
4. Security group for ALB: allow 80 and 443 from 0.0.0.0/0; outbound to app security group.

## 6. ECS Service

1. Create ECS service in cluster `clearstatus-cluster`:
   - Task definition: `clearstatus-backend`
   - Desired count: 1 (or more for HA).
   - Load balancer: attach to the ALB target group (container port 8080).
   - VPC: same as ALB and RDS; subnets: private.
   - Security group: allow inbound 8080 from ALB security group only; outbound to RDS and internet (for Google token verification).
2. Ensure health check on target group matches your app (e.g. HTTP 200 on `/api/public/orgs/...` or a dedicated `/health` if you add one).

## 7. Security Checklist

- RDS: no public access; only app security group can reach it.
- JWT secret and DB DSN: only in Secrets Manager/SSM and ECS task secrets; never in code or logs.
- ECS tasks: run in private subnets; only ALB can reach them.
- HTTPS only: redirect HTTP to HTTPS on ALB.
- IAM: minimal permissions for task execution and task role.

## 8. After Deploy

- Frontend: set `NEXT_PUBLIC_API_URL` to `https://your-alb-dns-name.or region.elb.amazonaws.com` or your custom domain.
- If you use a custom domain, add a CNAME to the ALB and attach the certificate to the listener.

## Optional: Health Check

Add a simple `/health` handler that returns 200 so the ALB target group health check can use it without hitting a specific org slug.

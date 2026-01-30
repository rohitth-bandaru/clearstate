# Deploy ClearStatus Frontend to Vercel

This guide covers deploying the Next.js frontend to Vercel securely.

## Prerequisites

- Vercel account
- Backend API URL (e.g. `https://api.yourdomain.com`)
- Google OAuth Client ID (same as backend `GOOGLE_CLIENT_ID`)

## 1. Environment Variables

In the Vercel project (Settings → Environment Variables), add:

| Name | Value | Notes |
|------|--------|------|
| `NEXT_PUBLIC_API_URL` | `https://your-api-url.com` | Backend API base URL (no trailing slash). Use production URL. |
| `NEXT_PUBLIC_GOOGLE_CLIENT_ID` | `xxx.apps.googleusercontent.com` | Google OAuth 2.0 Client ID (Web application). Same as backend. |

- **Security**: `NEXT_PUBLIC_*` vars are exposed to the browser. Do not put secrets here. JWT is stored in `localStorage` after sign-in; only the API URL and Google Client ID are public.
- Use **Production**, **Preview**, and **Development** as needed (e.g. set Preview to a staging API URL).

## 2. Deploy from Git

1. Push your frontend repo to GitHub/GitLab/Bitbucket.
2. In Vercel: **Add New Project** → Import the repository.
3. **Framework Preset**: Next.js (auto-detected).
4. **Root Directory**: Leave empty if the app is at the repo root; otherwise set to the frontend folder.
5. **Build Command**: `npm run build` (default).
6. **Output Directory**: leave default.
7. Add the environment variables above.
8. Deploy.

## 3. Custom Domain (optional)

- In Vercel: Project → Settings → Domains. Add your domain and follow DNS instructions.

## 4. After Deploy

- Open the Vercel URL. You should see the app; **Sign in with Google** will redirect to Google and then to your backend `/api/auth/google`. Ensure the backend is deployed. CORS is not used; all clients are expected to be in the same location.
- Public status page: `https://your-app.vercel.app/status/<org-slug>`. Share this with end users; it polls the backend every 2 seconds and updates only when data changes.

## 5. Security Checklist

- No server-side secrets in `NEXT_PUBLIC_*` (only API URL and Google Client ID).
- JWT is stored in `localStorage`; consider httpOnly cookies in the future if you add a cookie-based auth endpoint.

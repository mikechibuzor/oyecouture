# Oyecouture

Full-custom Oyecouture storefront: Nuxt 3, Go, and PostgreSQL.

## Local development

Requirements: Node.js 20+, Go 1.23+, and Docker with Compose.

```bash
cp .env.example .env
docker compose up -d postgres
```

Run the API:

```bash
cd api
go run ./cmd/server
```

Run the web app in another terminal:

```bash
cd web
npm install
npm run dev
```

The web app is available at `http://localhost:3000`. It calls the API health
endpoint at `http://localhost:8080/healthz`, which also checks PostgreSQL.

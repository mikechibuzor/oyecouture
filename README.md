# Oyecouture

Full-custom Oyecouture storefront: Nuxt 3, Go, and PostgreSQL.

## Local development

Requirements: Node.js 20.19+ or 22.12+, Go 1.23+, and Docker with Compose.

```bash
cp .env.example .env
docker compose up -d postgres
```

Run the API:

```bash
cd api
go run ./cmd/server
```

API startup applies the embedded PostgreSQL migrations and seeds the canonical
catalogue once. After intentionally changing `catalogue/product-data.js`, refresh
and verify the generated seed migrations:

```bash
cd catalogue
npm run generate:seed
npm run check:seed
```

Run the web app in another terminal:

```bash
cd web
npm install
npm run dev
```

The web app is available at `http://localhost:3000`. It calls the API health
endpoint at `http://localhost:8080/healthz`, which also checks PostgreSQL.

Public product reads are available from the Go API:

```bash
curl -sS http://localhost:8080/products
curl -sS 'http://localhost:8080/products?section=casuals-men&q=linen'
curl -sS http://localhost:8080/products/agbada
```

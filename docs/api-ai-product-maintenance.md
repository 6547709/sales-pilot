# Sales-Pilot: AI / MCP / OpenClaw — Product Maintenance API

How to call the **REST JSON API** to maintain **products (solutions)** from external agents, including **JWT lifetime**, **long-lived API keys**, field semantics, and examples.

---

## 1. Base URL

```text
{BASE}/api/v1
```

- Direct backend: e.g. `http://localhost:8080`
- Behind Nginx (recommended): `BASE` is your site origin; paths are `/api/v1/...`

Writes for products / scripts / cases live under **`/api/v1/admin/*`** and require an **admin** principal (JWT or API key).

---

## 2. Authentication

### 2.1 JWT (login)

```http
POST /api/v1/auth/login
Content-Type: application/json

{"username":"admin","password":"your-password"}
```

Use the returned **`access_token`**:

```http
Authorization: Bearer <access_token>
```

- **TTL**: access tokens expire after **24 hours** by default (`auth.TokenTTL` in the backend). After expiry, call `/auth/login` again.

### 2.2 API key (recommended for automation / MCP)

1. Log in to the **admin UI** as an administrator.
2. Open **「API 密钥」** (`/admin/api-keys`), create a key, **copy the full secret once** (it is shown only at creation).
3. On each admin request:

```http
X-API-Key: sp_<64-hex-chars>
```

- **Storage**: server stores **SHA256(full key)** only; the plaintext cannot be recovered.
- **Lifetime**: omit `expires_at` when creating → **no expiry** until revoked. You may set `expires_at` (RFC3339) for a bounded key.
- **Revoke**: `DELETE /api/v1/admin/api-keys/:id` (or use the admin UI). Revoked keys return **401**.

**Managing keys (admin only):**

| Method | Path | Body | Notes |
|--------|------|------|--------|
| GET | `/api/v1/admin/api-keys` | — | List keys (no secrets) |
| POST | `/api/v1/admin/api-keys` | `{"name":"mcp-prod","expires_at":"2030-12-31T23:59:59Z"}` | `expires_at` optional |
| DELETE | `/api/v1/admin/api-keys/:id` | — | Soft revoke (`is_active=false`) |

Creating the **first** API key still requires a normal **JWT login** (or local admin session). After that, scripts can use **only** `X-API-Key`.

### 2.3 Authorization

- Principal must be **`role: admin`** (the API key is bound to the **creating** admin user; if that user is demoted or disabled, the key stops working).

---

## 3. Machine-readable meta (for tools)

```http
GET /api/v1/admin/meta/product-maintain
Authorization: Bearer <jwt>
# or
X-API-Key: <key>
```

Returns JSON: field descriptions, endpoints, and a suggested agent workflow.

---

## 4. Product fields (for LLM filling)

| JSON field | Type | Purpose |
|------------|------|---------|
| `name` | string | Required on create. Customer-facing title. |
| `solution_category_id` | number \| null | Homepage topology solution id from `GET /api/v1/solution-categories`. |
| `category` | string | Display category; overwritten from solution label when `solution_category_id` is set. |
| `manufacturer_name` | string | Vendor / brand. |
| `sales_*` / `presales_*` | strings | Contacts shown on the public detail page. |
| `description` | string | Overview, **Markdown (GFM)**. |
| `highlights` | string[] | Three punchy value bullets. |
| `discovery_questions` | string[] | Open discovery questions. |
| `target_personas` | object | Role → one-line concern (e.g. `CEO`, `CIO`). |
| `trigger_events` | string | Plain text, buying triggers. |
| `competitor_analysis` | string | Markdown competitive write-up. |
| `roi_metrics` | string | ROI / metrics narrative. |
| `is_draft` | boolean | `true` = hidden from public catalog until published. |

**Do not** send nested `solution_category` on `PUT`; send **`solution_category_id`** only.

Public `GET /api/v1/products` **excludes drafts**; use `GET /api/v1/admin/products` for the full list.

---

## 5. curl examples (API key)

```bash
export BASE=http://localhost:8080
export KEY='sp_your_full_key_here'

curl -s "${BASE}/api/v1/admin/products" -H "X-API-Key: ${KEY}"

curl -s -X POST "${BASE}/api/v1/admin/products" \
  -H "X-API-Key: ${KEY}" \
  -H "Content-Type: application/json" \
  -d '{"name":"Example","description":"## Hi","highlights":["a","b","c"],"discovery_questions":[],"target_personas":{},"trigger_events":"","competitor_analysis":"","roi_metrics":"","is_draft":true}'
```

Same headers work for `GET /api/v1/admin/meta/product-maintain`, cases, scripts, etc.

---

## 6. MCP / OpenClaw

There is **no built-in MCP server** in this repo; expose these HTTPS JSON endpoints via your tool’s HTTP client. Recommended:

1. Fetch `GET /admin/meta/product-maintain` at startup (with `X-API-Key` or JWT).
2. Store the API key in **environment variables** or a secret manager, not in git.

---

## 7. Errors

- **401**: missing/invalid JWT or API key, expired key  
- **403**: non-admin user  
- **400**: bad JSON / validation  
- **404**: missing resource (or public access to a draft product)

Health: `GET /api/v1/health`

---

## 8. Code map

| Item | Location |
|------|----------|
| Admin auth (JWT + `X-API-Key`) | `backend/internal/middleware/admin_auth.go` |
| API key CRUD | `backend/internal/handlers/apikeys.go` |
| JWT TTL | `backend/internal/auth/jwt.go` → `TokenTTL` |
| Product model | `backend/internal/models/models.go` |
| Routes | `backend/internal/handlers/handlers.go` |

---

*Keep this file path stable for automation: `docs/api-ai-product-maintenance.md`.*

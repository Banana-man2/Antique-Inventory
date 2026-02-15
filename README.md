# Antique-Inventory

A Go web application for managing an antique gun inventory. Built with Gin (HTTP framework), Ent (ORM), and SQLite.

## Tech Stack

- **Go** - Language
- **Gin** - HTTP routing and HTML templates
- **Ent** - ORM and code generation
- **SQLite** via `modernc.org/sqlite` - Pure Go driver (no CGO/gcc required)

## Project Structure

```
Antique-Inventory/
  main.go              # Entry point, DB setup, route registration
  handlers/
    gun.go             # CRUD handlers (HTML + JSON API)
  templates/
    guns_list.html     # List all guns
    gun_detail.html    # Single gun detail view
    gun_form.html      # Create/edit form
  ent/
    schema/
      gun.go           # Ent schema definition
    ...                # Auto-generated ORM code
```

## Database Schema

Table: `guns`

| Field            | Type    | Notes                    |
|------------------|---------|--------------------------|
| gun_id           | INTEGER | Primary key              |
| gun_name         | TEXT    | Required, max 255 chars  |
| year             | INTEGER | Optional, nullable       |
| condition        | INTEGER | Optional, nullable       |
| description      | TEXT    | Optional, max 255 chars  |
| misc_attachments | TEXT    | Optional, max 255 chars  |
| createdAt        | TIME    | Auto-set on create       |
| updatedAt        | TIME    | Auto-set on create/update|

Database file: `C:\Coding\antique_inventory.db` (auto-created on first run via Ent migration).

## Running

```bash
cd C:\Coding\go_app
go run .
```

Server starts on http://localhost:8080.

## Routes

### HTML Pages

| Method | Path             | Description       |
|--------|------------------|-------------------|
| GET    | /guns            | List all guns     |
| GET    | /guns/new        | New gun form      |
| GET    | /guns/:id        | Gun detail page   |
| GET    | /guns/:id/edit   | Edit gun form     |
| POST   | /guns            | Create gun        |
| POST   | /guns/:id        | Update gun        |

### JSON API

| Method | Path             | Description       |
|--------|------------------|-------------------|
| GET    | /api/guns        | List all guns     |
| GET    | /api/guns/:id    | Get single gun    |
| POST   | /api/guns        | Create gun        |
| PUT    | /api/guns/:id    | Update gun        |
| DELETE | /api/guns/:id    | Delete gun        |

### API Examples

```bash
# List all guns
curl http://localhost:8080/api/guns

# Create a gun
curl -X POST http://localhost:8080/api/guns \
  -H "Content-Type: application/json" \
  -d '{"gun_name": "Colt 1851 Navy", "year": 1851, "condition": 7, "description": "Percussion revolver"}'

# Update a gun
curl -X PUT http://localhost:8080/api/guns/1 \
  -H "Content-Type: application/json" \
  -d '{"gun_name": "Colt 1851 Navy", "year": 1851, "condition": 8}'

# Delete a gun
curl -X DELETE http://localhost:8080/api/guns/1
```

## Regenerating Ent Code

After modifying `ent/schema/gun.go`:

```bash
go generate ./ent
```

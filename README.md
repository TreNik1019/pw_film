# pw_film

Programmierworkshop: REST Schnittstelle in Go

# Programmierworkshop am 19.6.2026

## Namen

- **Niklas Tremmel**
- **Marcus Akay**

## Link zum Git-Repository

[pw_film on GitHub](https://github.com/TreNik1019/pw_film)

## KI-Werkzeuge

### Agenten

- **Antigravity**
- **Gemini‑IDE**

### Chat-URLs

- <https://gemini.google.com>

## Frameworks und Bibliotheken

### REST‑Schnittstelle (Lesen und Neuanlegen)

- **Gin** – Hoch‑performantes HTTP‑Router‑Framework
- **GORM** – ORM für PostgreSQL, automatisierte Migrationen

### Validierung (nur Neuanlegen)

- Eingehende JSON‑Requests werden via `ShouldBindJSON` und strukturelle Tags validiert

### OR‑Mapping (für PostgreSQL)

- Nutzung von GORM‑Tags (`gorm:"primaryKey"`, `gorm:"type:date"`)
- Automatische Erstellung und Aktualisierung von Tabellen über `AutoMigrate`

### Optional: OIDC mit Keycloak

- *Nicht geschafft*

### Einfacher Integrationstest

- **Datenbank-Integrationstests**: Vollständige Tests gegen eine echte PostgreSQL-Datenbank, aufgeteilt in fokussierte Testdateien unter `tests/integration/`:
  - `create_test.go`: Testet das Anlegen (`POST /api/v1/films`).
  - `read_test.go`: Testet das Auslesen (`GET /api/v1/films`, `GET /api/v1/films/:id`, `404 Not Found`).
  - `update_test.go`: Testet das Aktualisieren (`PUT /api/v1/films/:id`).
  - `delete_test.go`: Testet das Löschen (`DELETE /api/v1/films/:id`).
  - `lifecycle_test.go`: Testet den gesamten CRUD-Lebenszyklus in einem Durchlauf.

## Prompts/Requests an KI-Agent/en

- Die vollständige, chronologische Liste aller relevanten(kopien und Serverprobleme ausgeschlossen) heute verwendeten Prompts befindet sich in der Datei [Prompts.md](file:///C:/workspace/pw_film/Prompts.md).

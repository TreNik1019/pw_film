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

### Einfacher Integrationstest

## Prompts/Requests an KI-Agent/en

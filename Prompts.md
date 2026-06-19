# Verlauf der Prompts von heute

Hier sind alle eingegebenen Prompts aus der heutigen Sitzung in einer überarbeiteten und professionellen Form dokumentiert:

## Prompt 1

Beschreibe mir, was mit dem Befehl `go get -u github.com/gin-gonic/gin` im Hintergrund durchgelaufen ist und welche Auswirkungen das auf unser Go-Projekt hat.

---

## Prompt 2

Agiere als erfahrener Go-Backend-Entwickler. Erstelle ein vollständiges, lauffähiges REST-API-Grundgerüst für das Projekt "pw_film" unter Verwendung des Gin Web Frameworks. Das Projekt muss einer sauberen, modularen Go-Standardarchitektur folgen, die API-Routing und Datenbanklogik strikt trennt.

---

## Prompt 3

Erstelle die initialen Verzeichnisstrukturen und Code-Dateien für das REST-API-Grundgerüst (Modulname `pw_film`), getrennt nach Handlern, Modellen und Konfiguration.

---

## Prompt 4

Führe das Setup des Routers und der Handler-Strukturen für die Film-Ressource fort.

---

## Prompt 5

Erweitere das Projekt um eine PostgreSQL-Anbindung. Erstelle dazu eine passende `compose.yml` für Docker, um die Datenbank in einem Container laufen zu lassen.

---

## Prompt 6

Die PostgreSQL-Zugangsdaten (Passwort) sollen aus Sicherheitsgründen in einer externen Textdatei namens `pw.txt` gespeichert werden. Passe die `database.go` so an, dass sie das Passwort dynamisch aus dieser Datei ausliest.

---

## Prompt 7

Setze den vorgeschlagenen Entwurf für die PostgreSQL-Integration und das Auslesen der Passwort-Datei in `database.go` um.

---

## Prompt 8

Optimiere das Einlesen der Passwort-Datei `pw.txt` in `database.go`. Reduziere die Komplexität und sorge für ein einfaches, direktes Auslesen des Passwort-Strings.

---

## Prompt 9

Erstelle im Projekt passende Verzeichnisse für SQL-Initialisierungsskripte und CSV-Seeddaten der PostgreSQL-Datenbank.

---

## Prompt 10

Verifikation des erfolgreichen Serverstarts und der PostgreSQL-Verbindung. Überprüfe die Log-Ausgabe auf Korrektheit und Vollständigkeit.

---

## Prompt 11

Gib mir drei Beispiel-`curl`-Befehle, um die REST-API abzufragen (z.B. alle Filme abrufen, einen Film per ID abrufen).

---

## Prompt 12

Wir möchten den Bruno REST-Client zur Abfrage der API nutzen. Erstelle die nötigen Bruno-Konfigurationsdateien, damit wir diese nur noch importieren müssen.

---

## Prompt 13

Warum benötigen wir die GORM-Methode `TableName() string` im Film-Modell, die den Standard-Plural überschreibt? Erkläre dies kurz.

---

## Prompt 14

Implementiere die neu hinzugekommenen Update- (PUT) und Delete- (DELETE) Operationen im Repository, im Handler und richte die entsprechenden HTTP-Routen im Router ein.

---

## Prompt 15

Generiere die fertigen `.bru`-Dateien für unsere CRUD-Routen, um sie direkt in Bruno importieren und nutzen zu können.

---

## Prompt 16

Behebe den Markdown-Lint-Fehler `MD022/blanks-around-headings` (Headings should be surrounded by blank lines) in der `README.md`.

---

## Prompt 17

Plane das Erstellen von Integrationstests in einem neuen Testverzeichnis, welche wir im Anschluss mit `go test` ausführen können.

---

## Prompt 18

Schreibe passende Integrationstests zur Überprüfung der REST-Schnittstelle gegen die echte PostgreSQL-Datenbank und richte das Testverzeichnis ein.

---

## Prompt 19

Führe das Schreiben der Integrationstests für die HTTP-Endpunkte fort.

---

## Prompt 20

Verwende für die Integrationstests die echte PostgreSQL-Datenbank anstelle von Mock-Datenbanken (Mocks).

---

## Prompt 21

Erkläre, wofür wir die Datei `mock_repository.go` noch benötigen, wenn wir die Tests direkt gegen die echte PostgreSQL-Datenbank laufen lassen.

---

## Prompt 22

Führe die Umstellung des Testsetups auf die echte Datenbank fort.

---

## Prompt 23

Führe das Refactoring und die Bereinigung des Testordners fort, um redundante Mock-Dateien zu entfernen.

---

## Prompt 24

Führe ein vollständiges Code-Review des Projekts durch: Überprüfe auf überflüssigen Code oder ungenutzte Dateien und behebe eventuelle Probleme.

---

## Prompt 25

Fehlerbehebung: Die Integrationstests schlagen beim Insert-Schritt fehl:
`ERROR: duplicate key value violates unique constraint "film_pkey" (SQLSTATE 23505)`
Das manuelle Einfügen von IDs beim Seeding stört das Auto-Increment. Behebe diesen Sequence Drift in PostgreSQL.

---

## Prompt 26

Implementiere den PostgreSQL Sequence-Korrektur-Schritt (`setval`) direkt im Setup-Code der Integrationstests, um den Drift auch dort abzufangen.

---

## Prompt 27

Führe die Integrationstests aus und verifiziere, ob nach dem Sequenz-Fix alle CRUD-Tests erfolgreich durchlaufen.

---

## Prompt 28

Erstelle eine GitHub-Actions-Pipeline (ci.yml), um bei jedem Push oder Pull-Request die Integrationstests automatisch auszuführen.

---

## Prompt 29

Teile die Integrationstests aus der großen Einzeldatei in separate, thematisch fokussierte Testdateien (`create_test.go`, `read_test.go`, etc.) auf, um die Wartbarkeit zu verbessern.

---

## Prompt 30

Überprüfe die GitHub-Actions-CI-Konfigurationsdatei. Stellen wir sicher, dass alle Run-Befehle stimmen und die Tests (inkl. PostgreSQL-Service) dort korrekt ausgeführt anstatt übersprungen werden.

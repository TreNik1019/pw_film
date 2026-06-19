// Package integration_test enthaelt Hilfsfunktionen fuer die Integrationstests.
package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"pw_film/internal/handlers"
	"pw_film/internal/models"
	"pw_film/internal/repository"
	"pw_film/internal/router"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// connectDB oeffnet eine Verbindung zur lokalen Test-Datenbank.
// Das Passwort wird aus pw.txt gelesen, genau wie in InitDB().
func connectDB(t *testing.T) *gorm.DB {
	t.Helper()

	// pw.txt suchen – relativ zum Projekt-Root, von dem go test ausgefuehrt wird
	candidates := []string{
		"internal/postgres/pw.txt",
		"../../internal/postgres/pw.txt",
	}
	var password string
	for _, c := range candidates {
		b, err := os.ReadFile(c)
		if err == nil {
			password = strings.TrimSpace(string(b))
			break
		}
	}
	if password == "" {
		t.Skip("pw.txt nicht gefunden – DB-Integrationstests uebersprungen")
	}

	dsn := fmt.Sprintf(
		"host=localhost user=postgres password=%s dbname=film port=5432 sslmode=disable TimeZone=Europe/Berlin",
		password,
	)

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Datenbankverbindung fehlgeschlagen (DB laeuft?): %v", err)
	}

	// Wichtig: AutoMigrate ausfuehren, damit die Tabellen auch in einer frischen CI-Datenbank existieren.
	if err := db.AutoMigrate(&models.Film{}); err != nil {
		t.Fatalf("AutoMigrate fehlgeschlagen: %v", err)
	}

	// Sequenz synchronisieren – CSV-Seeding fuegt explizite IDs ein ohne die Sequenz zu aktualisieren.
	// Dadurch schlagen INSERT-Statements ohne ID mit einem Duplicate-Key-Fehler fehl.
	if err := db.Exec("SELECT setval(pg_get_serial_sequence('film', 'id'), COALESCE((SELECT MAX(id) FROM film), 0))").Error; err != nil {
		t.Logf("Warnung: Sequenz konnte nicht synchronisiert werden: %v", err)
	}

	return db
}

// setupDBServer startet einen httptest.Server mit dem echten GORM-Repository.
func setupDBServer(t *testing.T) *httptest.Server {
	t.Helper()
	db := connectDB(t)
	repo := repository.NewFilmRepository(db)
	h := handlers.NewFilmHandler(repo)
	r := router.SetupRouter(h)
	return httptest.NewServer(r)
}

// postFilmDB legt einen Film in der Test-DB an.
func postFilmDB(t *testing.T, srv *httptest.Server, film models.Film) models.Film {
	t.Helper()
	body, _ := json.Marshal(film)
	resp, err := http.Post(srv.URL+"/api/v1/films", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST: erwartet 201, bekommen %d", resp.StatusCode)
	}
	var created models.Film
	json.NewDecoder(resp.Body).Decode(&created)
	return created
}

// deleteFilmDB loescht einen Film aus der Test-DB.
func deleteFilmDB(t *testing.T, srv *httptest.Server, id uint) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/films/%d", srv.URL, id), nil)
	client := &http.Client{}
	resp, _ := client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}
}

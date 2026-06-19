// Package integration enthält Integrationstests gegen die echte PostgreSQL-Datenbank.
// Voraussetzung: Die DB läuft (docker compose up -d) und pw.txt ist vorhanden.
//
// Ausführung:
//
//	go test ./tests/integration/... -v
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

// connectDB öffnet eine Verbindung zur lokalen Test-Datenbank.
// Das Passwort wird aus pw.txt gelesen, genau wie in InitDB().
func connectDB(t *testing.T) *gorm.DB {
	t.Helper()

	// pw.txt suchen – relativ zum Projekt-Root, von dem go test ausgeführt wird
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
		t.Skip("pw.txt nicht gefunden – DB-Integrationstests übersprungen")
	}

	dsn := fmt.Sprintf(
		"host=localhost user=postgres password=%s dbname=film port=5432 sslmode=disable TimeZone=Europe/Berlin",
		password,
	)

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Datenbankverbindung fehlgeschlagen (DB läuft?): %v", err)
	}

	// Sequenz synchronisieren – CSV-Seeding fügt explizite IDs ein ohne die Sequenz zu aktualisieren.
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

// ---- Hilfsfunktionen --------------------------------------------------------

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

func deleteFilmDB(t *testing.T, srv *httptest.Server, id uint) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/films/%d", srv.URL, id), nil)
	client := &http.Client{}
	resp, _ := client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}
}

// ---- Tests ------------------------------------------------------------------

// TestDB_CreateFilm prüft, ob ein Film in der echten DB angelegt wird.
func TestDB_CreateFilm(t *testing.T) {
	srv := setupDBServer(t)
	defer srv.Close()

	film := models.Film{
		Titel:     "Testfilm_Create",
		Art:       "BluRay",
		Genre:     "Test",
		Rating:    7.5,
		Verfuegbar: true,
		Preis:     9.99,
	}

	created := postFilmDB(t, srv, film)
	t.Cleanup(func() { deleteFilmDB(t, srv, created.ID) })

	if created.ID == 0 {
		t.Fatal("Kein Film angelegt – ID ist 0")
	}
	if created.Titel != film.Titel {
		t.Fatalf("Titel falsch: erwartet %q, bekommen %q", film.Titel, created.Titel)
	}
}

// TestDB_GetFilmByID prüft, ob ein angelegter Film per ID abgerufen werden kann.
func TestDB_GetFilmByID(t *testing.T) {
	srv := setupDBServer(t)
	defer srv.Close()

	created := postFilmDB(t, srv, models.Film{Titel: "Testfilm_Get", Rating: 8.0})
	t.Cleanup(func() { deleteFilmDB(t, srv, created.ID) })

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/films/%d", srv.URL, created.ID))
	if err != nil {
		t.Fatalf("GET fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d", resp.StatusCode)
	}

	var fetched models.Film
	json.NewDecoder(resp.Body).Decode(&fetched)
	if fetched.Titel != "Testfilm_Get" {
		t.Fatalf("Titel nach GET falsch: %s", fetched.Titel)
	}
}

// TestDB_GetFilmByID_NotFound prüft, ob eine nicht existierende ID 404 liefert.
func TestDB_GetFilmByID_NotFound(t *testing.T) {
	srv := setupDBServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/films/999999")
	if err != nil {
		t.Fatalf("GET fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("erwartet 404, bekommen %d", resp.StatusCode)
	}
}

// TestDB_GetFilms prüft, ob die Liste mindestens den angelegten Testfilm enthält.
func TestDB_GetFilms(t *testing.T) {
	srv := setupDBServer(t)
	defer srv.Close()

	created := postFilmDB(t, srv, models.Film{Titel: "Testfilm_List", Rating: 6.0})
	t.Cleanup(func() { deleteFilmDB(t, srv, created.ID) })

	resp, err := http.Get(srv.URL + "/api/v1/films")
	if err != nil {
		t.Fatalf("GET /films fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d", resp.StatusCode)
	}

	var films []models.Film
	json.NewDecoder(resp.Body).Decode(&films)
	if len(films) == 0 {
		t.Fatal("Film-Liste ist leer – mindestens 1 Film erwartet")
	}

	// prüfen ob unser Testfilm in der Liste ist
	found := false
	for _, f := range films {
		if f.ID == created.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Testfilm mit ID %d nicht in der Liste gefunden", created.ID)
	}
}

// TestDB_UpdateFilm prüft, ob ein Film korrekt aktualisiert wird.
func TestDB_UpdateFilm(t *testing.T) {
	srv := setupDBServer(t)
	defer srv.Close()

	created := postFilmDB(t, srv, models.Film{Titel: "Testfilm_Update", Preis: 12.99})
	t.Cleanup(func() { deleteFilmDB(t, srv, created.ID) })

	created.Preis = 7.99
	body, _ := json.Marshal(created)
	req, _ := http.NewRequest(http.MethodPut,
		fmt.Sprintf("%s/api/v1/films/%d", srv.URL, created.ID),
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("PUT fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("erwartet 200, bekommen %d", resp.StatusCode)
	}

	// Verifizieren per GET
	getResp, _ := http.Get(fmt.Sprintf("%s/api/v1/films/%d", srv.URL, created.ID))
	defer getResp.Body.Close()
	var updated models.Film
	json.NewDecoder(getResp.Body).Decode(&updated)
	if updated.Preis != 7.99 {
		t.Fatalf("Preis nicht aktualisiert: erwartet 7.99, bekommen %f", updated.Preis)
	}
}

// TestDB_DeleteFilm prüft, ob ein Film korrekt gelöscht wird.
func TestDB_DeleteFilm(t *testing.T) {
	srv := setupDBServer(t)
	defer srv.Close()

	created := postFilmDB(t, srv, models.Film{Titel: "Testfilm_Delete", Rating: 5.0})
	// Kein Cleanup nötig – der Test löscht den Film selbst

	req, _ := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/api/v1/films/%d", srv.URL, created.ID), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("DELETE fehlgeschlagen: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("erwartet 204, bekommen %d", resp.StatusCode)
	}

	// Nachher: GET sollte 404 liefern
	getResp, _ := http.Get(fmt.Sprintf("%s/api/v1/films/%d", srv.URL, created.ID))
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusNotFound {
		t.Fatalf("nach DELETE erwartet 404, bekommen %d", getResp.StatusCode)
	}
}

// TestDB_Lifecycle_CRUD prüft den vollständigen CRUD-Lebenszyklus eines Films.
func TestDB_Lifecycle_CRUD(t *testing.T) {
	srv := setupDBServer(t)
	defer srv.Close()

	// CREATE
	film := models.Film{Titel: "Testfilm_CRUD", Rating: 9.0, Preis: 14.99, Verfuegbar: true}
	created := postFilmDB(t, srv, film)
	if created.ID == 0 {
		t.Fatal("Film wurde nicht angelegt")
	}

	// READ
	getResp, _ := http.Get(fmt.Sprintf("%s/api/v1/films/%d", srv.URL, created.ID))
	var fetched models.Film
	json.NewDecoder(getResp.Body).Decode(&fetched)
	getResp.Body.Close()
	if fetched.Titel != "Testfilm_CRUD" {
		t.Fatalf("Titel nach GET falsch: %s", fetched.Titel)
	}

	// UPDATE
	fetched.Preis = 8.49
	body, _ := json.Marshal(fetched)
	req, _ := http.NewRequest(http.MethodPut,
		fmt.Sprintf("%s/api/v1/films/%d", srv.URL, fetched.ID),
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	putResp, _ := client.Do(req)
	putResp.Body.Close()
	if putResp.StatusCode != http.StatusOK {
		t.Fatalf("UPDATE fehlgeschlagen mit Status %d", putResp.StatusCode)
	}

	// Verify UPDATE
	verResp, _ := http.Get(fmt.Sprintf("%s/api/v1/films/%d", srv.URL, fetched.ID))
	var verified models.Film
	json.NewDecoder(verResp.Body).Decode(&verified)
	verResp.Body.Close()
	if verified.Preis != 8.49 {
		t.Fatalf("Preis nach UPDATE falsch: %f", verified.Preis)
	}

	// DELETE
	delReq, _ := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/api/v1/films/%d", srv.URL, fetched.ID), nil)
	delResp, _ := client.Do(delReq)
	delResp.Body.Close()
	if delResp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE fehlgeschlagen mit Status %d", delResp.StatusCode)
	}

	// Verify DELETE
	getAfterDel, _ := http.Get(fmt.Sprintf("%s/api/v1/films/%d", srv.URL, fetched.ID))
	getAfterDel.Body.Close()
	if getAfterDel.StatusCode != http.StatusNotFound {
		t.Fatalf("nach DELETE erwartet 404, bekommen %d", getAfterDel.StatusCode)
	}
}

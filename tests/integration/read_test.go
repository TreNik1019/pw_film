package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pw_film/internal/models"
	"testing"
)

// TestDB_GetFilmByID prueft, ob ein angelegter Film per ID abgerufen werden kann.
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

// TestDB_GetFilmByID_NotFound prueft, ob eine nicht existierende ID 404 liefert.
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

// TestDB_GetFilms prueft, ob die Liste mindestens den angelegten Testfilm enthaelt.
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

	// pruefen ob unser Testfilm in der Liste ist
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

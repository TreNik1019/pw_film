package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pw_film/internal/models"
	"testing"
)

// TestDB_UpdateFilm prueft, ob ein Film korrekt aktualisiert wird.
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

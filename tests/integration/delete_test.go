package integration_test

import (
	"fmt"
	"net/http"
	"pw_film/internal/models"
	"testing"
)

// TestDB_DeleteFilm prueft, ob ein Film korrekt geloescht wird.
func TestDB_DeleteFilm(t *testing.T) {
	srv := setupDBServer(t)
	defer srv.Close()

	created := postFilmDB(t, srv, models.Film{Titel: "Testfilm_Delete", Rating: 5.0})
	// Kein Cleanup noetig – der Test loescht den Film selbst

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

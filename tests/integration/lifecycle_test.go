package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pw_film/internal/models"
	"testing"
)

// TestDB_Lifecycle_CRUD prueft den vollstaendigen CRUD-Lebenszyklus eines Films.
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

package integration_test

import (
	"pw_film/internal/models"
	"testing"
)

// TestDB_CreateFilm prueft, ob ein Film in der echten DB angelegt wird.
func TestDB_CreateFilm(t *testing.T) {
	srv := setupDBServer(t)
	defer srv.Close()

	film := models.Film{
		Titel:      "Testfilm_Create",
		Art:        "BluRay",
		Genre:      "Test",
		Rating:     7.5,
		Verfuegbar: true,
		Preis:      9.99,
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

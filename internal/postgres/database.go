package postgres

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"pw_film/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func getPasswordPath() string {
	candidates := []string{
		"internal/postgres/pw.txt",       // Ausführung vom Root-Verzeichnis
		"../../internal/postgres/pw.txt", // Ausführung von cmd/server
		"pw.txt",                         // Ausführung von internal/postgres
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return "internal/postgres/pw.txt"
}

func InitDB() {
	host := "localhost"
	user := "postgres"
	dbName := "film"
	port := "5432"

	// Passwort robust aus Datei 'pw.txt' lesen (unabhängig vom Ausführungsordner)
	passwordBytes, err := os.ReadFile(getPasswordPath())
	if err != nil {
		log.Fatalf("Fehler beim Lesen der Passwort-Datei 'pw.txt': %v. Bitte stelle sicher, dass die Datei 'pw.txt' existiert.", err)
	}
	password := strings.TrimSpace(string(passwordBytes))

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Berlin",
		host, user, password, dbName, port)

	var errOpen error
	DB, errOpen = gorm.Open(postgres.Open(dsn), &gorm.Config{
		FullSaveAssociations: false,
	})

	if errOpen != nil {
		log.Fatalf("Fehler beim Verbinden mit der bestehenden Postgres-DB: %v", errOpen)
	}

	log.Println("Erfolgreich mit bestehender PostgreSQL-Datenbank verbunden!")

	// Automatische Migration für alle Modelle
	if err := DB.AutoMigrate(&models.Film{}, &models.Regisseur{}, &models.Cover{}); err != nil {
		log.Fatalf("Fehler bei der Datenbank-Migration: %v", err)
	}
	log.Println("Datenbank-Migration abgeschlossen.")

	// Datenbank mit CSV-Daten befüllen, falls leer
	SeedDatabase(DB)
}

func getCSVPath() string {
	candidates := []string{
		"internal/postgres/csv/film.csv",       // Ausführung vom Root-Verzeichnis
		"../../internal/postgres/csv/film.csv", // Ausführung von cmd/server
		"csv/film.csv",                         // Ausführung von internal/postgres
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return "internal/postgres/csv/film.csv"
}

// SeedDatabase liest die CSV-Datei aus und fügt die Filme in die Datenbank ein, falls diese leer ist
func SeedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&models.Film{}).Count(&count)
	if count > 0 {
		return // Bereits befüllt
	}

	csvPath := getCSVPath()
	file, err := os.Open(csvPath)
	if err != nil {
		log.Printf("CSV Seeding übersprungen (Datei '%s' nicht gefunden): %v", csvPath, err)
		return
	}
	defer file.Close()

	// CSV-Reader mit Semikolon als Trennzeichen initialisieren
	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Header-Zeile überspringen
	if _, err := reader.Read(); err != nil {
		log.Printf("Fehler beim Lesen des CSV-Headers: %v", err)
		return
	}

	var insertedCount int
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Fehler beim Lesen einer CSV-Zeile: %v", err)
			continue
		}

		// CSV-Zeilenwerte parsen
		id, _ := strconv.ParseUint(record[0], 10, 64)
		version, _ := strconv.Atoi(record[1])
		titel := record[2]
		art := record[3]
		
		// Datum parsen (Format: YYYY-MM-DD)
		erschDatum, errDate := time.Parse("2006-01-02", record[4])
		if errDate != nil {
			// Alternativ als fallback probieren
			erschDatum = time.Now()
		}

		genre := record[5]
		rating, _ := strconv.ParseFloat(record[6], 64)
		verfuegbar := strings.ToLower(record[7]) == "true"
		preis, _ := strconv.ParseFloat(record[8], 64)
		schlagwoerter := record[9]

		// Erstellungs- und Aktualisierungsdaten parsen (Format: YYYY-MM-DDTHH:MM:SS)
		erzeugt, _ := time.Parse("2006-01-02T15:04:05", record[10])
		aktualisiert, _ := time.Parse("2006-01-02T15:04:05", record[11])

		film := models.Film{
			ID:                uint(id),
			Version:           version,
			Titel:             titel,
			Art:               art,
			Erscheinungsdatum: erschDatum,
			Genre:             genre,
			Rating:            rating,
			Verfuegbar:        verfuegbar,
			Preis:             preis,
			Schlagwoerter:     schlagwoerter,
			Erzeugt:           erzeugt,
			Aktualisiert:      aktualisiert,
		}

		if err := db.Create(&film).Error; err != nil {
			log.Printf("Fehler beim Einfügen des Films '%s': %v", titel, err)
		} else {
			insertedCount++
		}
	}

	if insertedCount > 0 {
		log.Printf("Datenbank erfolgreich mit %d Filmen aus der CSV-Datei befüllt!", insertedCount)
	}
}

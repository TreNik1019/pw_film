package postgres

import (
	"fmt"
	"log"
	"os"
	"strings"

	"pw_film/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	host := "localhost"
	user := "postgres"
	dbName := "pw_film"
	port := "5432"

	// Passwort aus Datei 'pw' im Hauptverzeichnis lesen
	passwordBytes, err := os.ReadFile("pw")
	if err != nil {
		log.Fatalf("Fehler beim Lesen der Passwort-Datei 'pw': %v. Bitte stelle sicher, dass die Datei 'pw' im Hauptverzeichnis existiert.", err)
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
}

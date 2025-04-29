package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Créer les fichiers de migration
	if err := createMigrationFiles(); err != nil {
		log.Fatalf("Erreur lors de la création des fichiers de migration: %v", err)
	}

	// Ouvrir la connexion SQLite
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalf("Erreur lors de l'ouverture de la base de données: %v", err)
	}
	defer db.Close()

	// Créer l'instance du driver SQLite pour les migrations
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatalf("Erreur lors de la création du driver SQLite: %v", err)
	}

	// Initialiser la migration
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"sqlite3", driver)
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la migration: %v", err)
	}

	// Appliquer toutes les migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Erreur lors de l'application des migrations: %v", err)
	}

	fmt.Println("Migrations appliquées avec succès!")
}

func createMigrationFiles() error {
	// Créer le dossier migrations s'il n'existe pas
	if err := os.MkdirAll("migrations", 0755); err != nil {
		return err
	}

	// Créer le fichier de la première migration
	migration1 := `CREATE TABLE "user" ("id" TEXT NOT NULL, "email_address" TEXT NOT NULL, "first_name" TEXT NOT NULL, "last_name" TEXT NOT NULL, PRIMARY KEY ("id"));

CREATE TABLE "article" ("id" TEXT NOT NULL, "user_id" TEXT NOT NULL, "title" TEXT NOT NULL, "content" TEXT NOT NULL, CONSTRAINT "article_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "user"("id") ON DELETE CASCADE ON UPDATE CASCADE, PRIMARY KEY ("id"));
CREATE INDEX "article_user_id_index" ON "article" ("user_id");

CREATE TABLE "clap" ("id" TEXT NOT NULL, "article_id" TEXT NOT NULL, "count" INTEGER NOT NULL, CONSTRAINT "clap_article_id_foreign" FOREIGN KEY("article_id") REFERENCES "article"("id") ON DELETE CASCADE ON UPDATE CASCADE, PRIMARY KEY ("id"));
CREATE INDEX "clap_article_id_index" ON "clap" ("article_id");`

	// Fichier up pour la première migration
	if err := os.WriteFile("migrations/000001_create_initial_tables.up.sql", []byte(migration1), 0644); err != nil {
		return err
	}
	// Fichier down pour la première migration
	if err := os.WriteFile("migrations/000001_create_initial_tables.down.sql", []byte("DROP TABLE IF EXISTS clap; DROP TABLE IF EXISTS article; DROP TABLE IF EXISTS user;"), 0644); err != nil {
		return err
	}

	// Créer le fichier de la deuxième migration
	migration2 := `CREATE TABLE "user_profile_view" ("id" TEXT NOT NULL, "content" JSON NOT NULL, PRIMARY KEY ("id"));`

	// Fichier up pour la deuxième migration
	if err := os.WriteFile("migrations/000002_create_user_profile_view.up.sql", []byte(migration2), 0644); err != nil {
		return err
	}
	// Fichier down pour la deuxième migration
	if err := os.WriteFile("migrations/000002_create_user_profile_view.down.sql", []byte("DROP TABLE IF EXISTS user_profile_view;"), 0644); err != nil {
		return err
	}

	return nil
}

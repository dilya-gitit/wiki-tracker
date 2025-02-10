package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS stats (
			date TEXT,
			lang TEXT,
			edits INTEGER,
			offset INTEGER,
			PRIMARY KEY (date, lang)
		);
		CREATE TABLE IF NOT EXISTS user_lang (
			user_id TEXT PRIMARY KEY,
			lang TEXT
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	return db
}

func TestStoreEditCount(t *testing.T) {
	db = setupTestDB(t)
	storeEditCounts("en", "2025-02-10", 1)

	var count int
	err := db.QueryRow("SELECT edits FROM stats WHERE date = ? AND lang = ?", "2025-02-10", "en").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query edit count: %v", err)
	}

	if count != MAX_BATCH_SIZE {
		t.Errorf("Expected 100 edit, got %d", count)
	}
}

func TestGetEditCount(t *testing.T) {
	db = setupTestDB(t)
	db.Exec("INSERT INTO stats (date, lang, edits, offset) VALUES (?, ?, ?, ?)", "2025-02-10", "en", 3, 1)
	bootstrap()
	count := getEditCount("en", "2025-02-10")
	if count != 3 {
		t.Errorf("Expected 3 edits, got %d", count)
	}
}

func TestSetAndGetUserLanguage(t *testing.T) {
	db = setupTestDB(t)
	setUserLanguage("user123", "fr")

	lang := getUserLanguage("user123")
	if lang != "fr" {
		t.Errorf("Expected language 'fr', got '%s'", lang)
	}
}

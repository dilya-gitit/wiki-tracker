package main

import (
	"database/sql"
	"log"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func storeEditCounts(lang string, date string, offset int) {
	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return
	}

	sql := fmt.Sprintf("INSERT INTO stats (date, lang, offset, edits) VALUES (?, ?, ?, %d) ON CONFLICT(date, lang) DO UPDATE SET edits = edits + %d, offset = %d", MAX_BATCH_SIZE, MAX_BATCH_SIZE, offset)

	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Println("Failed to prepare statement:", err)
		tx.Rollback()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(date, lang, offset)
	if err != nil {
		log.Println("Failed to insert edit count:", err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Failed to commit transaction:", err)
	}
}

func getAllStats() ([]Stat, error) {
	log.Println("Querying all stats")

	rows, err := db.Query("SELECT date, lang, edits, offset FROM stats")
	if err != nil {
		log.Println("Failed to query stats:", err)
		return nil, err
	}
	defer rows.Close()

	var stats []Stat

	for rows.Next() {
		var stat Stat

		err := rows.Scan(&stat.Date, &stat.Lang, &stat.Edits, &stat.Offset)
		if err != nil {
			log.Println("Failed to scan row:", err)
			continue
		}

		// Append the stat object to the slice
		stats = append(stats, stat)
	}

	// Check for errors from iteration
	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
		return nil, err
	}

	log.Println("Retrieved all stats:", stats)
	return stats, nil
}

func getUserLanguage(userID string) string {
	var lang string
	err := db.QueryRow("SELECT lang FROM user_lang WHERE user_id = ?", userID).Scan(&lang)
	if err != nil {
		log.Println("No language preference found for user:", userID, "- Defaulting to English")
		return "en"
	}
	return lang
}

func setUserLanguage(userID, lang string) {
	if lang == "" {
		lang = "en"
	}
	_, err := db.Exec("INSERT INTO user_lang (user_id, lang) VALUES (?, ?) ON CONFLICT(user_id) DO UPDATE SET lang = ?", userID, lang, lang)
	if err != nil {
		log.Println("Failed to set user language:", err)
	}
}

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "wikipedia_bot.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
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
		log.Fatal("Failed to create tables:", err)
	}
}

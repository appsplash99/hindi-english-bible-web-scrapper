// Script #3
// This script will create a SQLite DB file
// and populate it with the verses from the JSON file.

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SupportedLanguages struct {
	English string `json:"english"`
	Hindi   string `json:"hindi"`
}

type Verses struct {
	Verse_id         int    `json:"verse_id"`
	VerseCount       string `json:"verse_count"`
	Hindi_BookName   string `json:"hindi_book_name"`
	English_BookName string `json:"english_book_name"`
	English_Verse    string `json:"english_verse"`
	Hindi_Verse      string `json:"hindi_verse"`
}

func main() {
	// Open the JSON file
	file, err := os.Open("hindi_english_bible__final.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close()

	// Parse the JSON file into a slice of Verses struct
	var verses []Verses
	err = json.NewDecoder(file).Decode(&verses)
	if err != nil {
		fmt.Println("Error parsing JSON file:", err)
		return
	}

	// Create the SQLite DB file
	db, err := sql.Open("sqlite3", "hindiEnglishBible.db")
	if err != nil {
		fmt.Println("Error creating SQLite DB:", err)
		return
	}
	defer db.Close()

	// Create the "verses" table in the DB
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS verses (
			verse_id INTEGER PRIMARY KEY,
			english_book_name VARCHAR(255),
			hindi_book_name VARCHAR(255),
			English_Verse TEXT,
			Hindi_Verse TEXT
		);
	`)
	if err != nil {
		fmt.Println("Error creating verses table:", err)
		return
	}

	// Insert the verses into the DB
	stmt, err := db.Prepare(`
		INSERT INTO verses (
			verse_id,
			english_book_name,
			hindi_book_name,
			English_Verse,
			Hindi_Verse
		) VALUES (?, ?, ?, ?, ?);
	`)
	if err != nil {
		fmt.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	for _, verse := range verses {
		_, err = stmt.Exec(
			verse.VerseCount,
			verse.English_BookName,
			verse.Hindi_BookName,
			verse.English_Verse,
			verse.Hindi_Verse,
		)
		if err != nil {
			fmt.Println("Error inserting verse into DB:", err)
			return
		}
	}

	fmt.Println("Successfully created Bible SQLite DB.")
}

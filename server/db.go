package main

import (
	"database/sql"
	"log"
	"os"
)

func ConnectDB(dbname string) *sql.DB {
	var initDB = false
	_, err := os.Stat(dbname)
	if os.IsNotExist(err) {
		// init db
		initDB = true
	}

	// os.Remove("./appcast.db")
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}

	if initDB {
		log.Println("Initializing database schema...")
		createAccountTable(db)
		createReleaseTable(db)
		createChannelTable(db)
	}
	return db
}

func createAccountTable(db *sql.DB) {
	if _, err := db.Exec(`create table account(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account varchar,
		token varchar
	);`); err != nil {
		log.Fatal(err)
	}
}

func createChannelTable(db *sql.DB) {
	if _, err := db.Exec(`create table channels(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title varchar,
		description text,
		identity varchar,
		token varchar
	);`); err != nil {
		log.Fatal(err)
	}
	/*
		http://localhost:8080/appcast/gotray/4cbd040533a2f43fc6691d773d510cda70f4126a
	*/
	if _, err := db.Exec(`insert into channels(title,description, identity, token) values (?,?,?,?)`, "GoTray", "Desc", "gotray", "4cbd040533a2f43fc6691d773d510cda70f4126a"); err != nil {
		panic(err)
	}
}

func createReleaseTable(db *sql.DB) {
	if _, err := db.Exec(`create table releases(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title varchar,
		desc text,
		releaseNotes text,
		pubDate datetime default current_timestamp,
		filename varchar,
		channel varchar,
		length integer,
		mimetype varchar,
		version varchar,
		shortVersionString varchar,
		dsaSignature varchar,
		token varchar
	);`); err != nil {
		log.Fatal(err)
	}
}

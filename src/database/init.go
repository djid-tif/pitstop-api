package database

import (
	"database/sql"
	"fmt"
	"os"
)

func init() {
	stat, err := os.Stat(databaseFile)
	if err == nil {
		if stat.IsDir() {
			panic("database file could not be a directory !")
		}
		return
	}

	fmt.Println("Création de la base de données ...")
	_, err = os.Create(databaseFile)
	if err != nil {
		panic(fmt.Sprintf("échec de la création de la base de données: %v", err))
	}

	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createUsersTable(db)
	createAwaitTable(db)
	createCategoriesTable(db)
	createTopicsTable(db)
	createPostsTable(db)
	createReportsTable(db)
	createReportMessagesTable(db)
}

func createUsersTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE "users" (
	"id" INTEGER,
	"username" TEXT NOT NULL UNIQUE,
	"email"	TEXT NOT NULL UNIQUE,
	"password" TEXT NOT NULL,
	"registration_date"	INTEGER NOT NULL,
	"rank" INTEGER DEFAULT 1,
	"last_password_refresh" INTEGER NOT NULL,
	"friends" TEXT DEFAULT '[]',
	"friends_sent" TEXT DEFAULT '[]',
	"friends_received" TEXT DEFAULT '[]',
	"settings" TEXT NOT NULL,
	"stats" TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);`)
	if err != nil {
		panic(err)
	}
}
func createAwaitTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE "awaitUsers" (
	"id" INTEGER,
	"username" TEXT NOT NULL UNIQUE,
	"email"	TEXT NOT NULL UNIQUE,
	"password" TEXT NOT NULL,
	"registration_date"	INTEGER NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);`)
	if err != nil {
		panic(err)
	}
}

func createCategoriesTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE "categories" (
	"id" INTEGER,
	"name" TEXT NOT NULL UNIQUE,
	"slug" TEXT NOT NULL UNIQUE,
	"creation_date" INTEGER NOT NULL,
	"topics" TEXT DEFAULT '[]',
	"icon" TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);`)
	if err != nil {
		panic(err)
	}
}

func createTopicsTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE "topics" (
	"id" INTEGER,
	"parent_id" INTEGER NOT NULL,
	"name" TEXT NOT NULL UNIQUE,
	"slug" TEXT NOT NULL UNIQUE,
	"author_id" INTEGER NOT NULL,
	"creation_date" INTEGER NOT NULL,
	"posts" TEXT DEFAULT '[]',
	"pinned" INTEGER DEFAULT 0,
	"locked" INTEGER DEFAULT 0,
	"tags" TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);`)
	if err != nil {
		panic(err)
	}
}

func createPostsTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE "posts" (
	"id" INTEGER,
	"parent_id" INTEGER NOT NULL,
	"content" TEXT NOT NULL,
	"author_id" INTEGER NOT NULL,
	"creation_date" INTEGER NOT NULL,
	"likes" TEXT DEFAULT '[]',
	"dislikes" TEXT DEFAULT '[]',
	PRIMARY KEY("id" AUTOINCREMENT)
);`)
	if err != nil {
		panic(err)
	}
}

func createReportsTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE "reports" (
	"id" INTEGER,
	"type" TEXT NOT NULL,
	"target_id" INTEGER NOT NULL,
	"creation_date" INTEGER NOT NULL,
	"author_id" INTEGER NOT NULL,
	"messages" TEXT DEFAULT '[]',
	PRIMARY KEY("id" AUTOINCREMENT)
);`)
	if err != nil {
		panic(err)
	}
}

func createReportMessagesTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE "report_messages" (
	"id" INTEGER,
	"parent_id" INTEGER NOT NULL,
	"creation_date" INTEGER NOT NULL,
	"author_id" INTEGER NOT NULL,
	"content" TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);`)
	if err != nil {
		panic(err)
	}
}

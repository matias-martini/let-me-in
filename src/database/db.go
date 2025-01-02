package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func CreateDBIfNotExists() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Connect to postgres maintenance database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password)

	// Connect to postgres system database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("failed to connect: %v", err))
	}
	defer db.Close()

	// Ping to ensure connection is valid
	err = db.Ping()
	if err != nil {
		// If postgres database doesn't exist, try template1
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=template1 sslmode=disable",
			host, port, user, password)
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to template1: %v", err))
		}
		err = db.Ping()
		if err != nil {
			panic(fmt.Sprintf("failed to ping template1: %v", err))
		}
	}

	// Check if database exists
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)")
	err = db.QueryRow(query, dbname).Scan(&exists)
	if err != nil {
		panic(fmt.Sprintf("failed to check database existence: %v", err))
	}

	if !exists {
		// Create the database
		createDBQuery := fmt.Sprintf("CREATE DATABASE %s;", dbname)
		_, err := db.Exec(createDBQuery)
		if err != nil {
			panic(fmt.Sprintf("failed to create database: %v", err))
		}
		fmt.Printf("Database '%s' created successfully\n", dbname)
	} else {
		fmt.Printf("Database '%s' already exists\n", dbname)
	}
}

func Init() {
	// Get the database connection details from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Database connection established")

	DB = db
}

func InitTestDB() {
	Init()
	DB = DB.Begin()
}

func ResetTestDB() {
	DB.Rollback()
}


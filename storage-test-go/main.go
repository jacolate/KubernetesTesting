package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

const (
    dbUser         = "root"
    dbPassword     = "rootpassword"
    dbHost         = "mysql"
    dbName         = "testdb"
    dataCheckDelay = 1 * time.Second
    dbTimeout      = 100 * time.Millisecond
)

var (
    downtimeCounter time.Duration
)

func main() {
    // Wait for MySQL to be ready
    time.Sleep(5 * time.Second)

    db := connectToDatabase()
    defer db.Close()

    createDatabase(db)
    db = switchToDatabase(db)
    defer db.Close()

    createTable(db)
    insertInitialData(db)

    verifyDataLoop(db)
}

func connectToDatabase() *sql.DB {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/", dbUser, dbPassword, dbHost)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    return db
}

func createDatabase(db *sql.DB) {
    _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
    if err != nil {
        log.Fatalf("Error creating database: %v", err)
    }
}

func switchToDatabase(db *sql.DB) *sql.DB {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    return db
}

func createTable(db *sql.DB) {
    _, err := db.Exec("CREATE TABLE IF NOT EXISTS test (id INT AUTO_INCREMENT PRIMARY KEY, data VARCHAR(255))")
    if err != nil {
        log.Fatalf("Error creating table: %v", err)
    }
}

func insertInitialData(db *sql.DB) {
    _, err := db.Exec("INSERT INTO test (data) VALUES ('testdata')")
    if err != nil {
        log.Fatalf("Error inserting initial data: %v", err)
    }
    log.Println("Initial data created")
}

func verifyDataLoop(db *sql.DB) {
    for {
        start := time.Now()
        data, err := readDataWithTimeout(db)
        elapsed := time.Since(start)

        if err != nil {
            if err.Error() == "timeout reached while reading data" {
                downtimeCounter += dbTimeout
            } else {
                log.Printf("Error reading data: %v", err)
                downtimeCounter += elapsed
            }
        } else if data != "testdata" {
            log.Fatalf("Data verification failed: expected 'testdata', got '%s'", data)
        } else {
            log.Println("Data verification successful")
        }

        log.Printf("Total downtime: %v", downtimeCounter)
        time.Sleep(dataCheckDelay)
    }
}

func readDataWithTimeout(db *sql.DB) (string, error) {
    var data string
    done := make(chan error, 1)

    go func() {
        err := db.QueryRow("SELECT data FROM test WHERE id = 1").Scan(&data)
        done <- err
    }()

    select {
    case err := <-done:
        return data, err
    case <-time.After(dbTimeout):
        return "", fmt.Errorf("timeout reached while reading data")
    }
}

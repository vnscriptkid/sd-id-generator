package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IdRange struct {
	Service    string `gorm:"primaryKey"`
	LastUsedId int64
}

type IdRangeResponse struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

const rangeSize = 1000

var db *sql.DB

func main() {
	// Initialize database connection
	var err error
	db, err = sql.Open("mysql", "user:user_password@tcp(127.0.0.1:3306)/my_database")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Auto migrate using GORM
	gormDB, err := gorm.Open(mysql.Open("user:user_password@tcp(127.0.0.1:3306)/my_database"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = gormDB.AutoMigrate(&IdRange{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	// Set up Gin router
	r := gin.Default()
	r.GET("/next-id-range/:service", getNextIdRange)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getNextIdRange(c *gin.Context) {
	service := c.Param("service")

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to start transaction: %v", err)})
		return
	}
	defer tx.Rollback()

	// Lock and get the current last_used_id
	var lastUsedId int64
	err = tx.QueryRow("SELECT last_used_id FROM id_ranges WHERE service = ? FOR UPDATE", service).Scan(&lastUsedId)
	if err == sql.ErrNoRows {
		// Insert new row if not exists
		_, err = tx.Exec("INSERT INTO id_ranges (service, last_used_id) VALUES (?, 0)", service)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new service"})
			return
		}
		lastUsedId = 0
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Database error: %v", err)})
		return
	}

	start := lastUsedId + 1
	end := start + rangeSize - 1

	// Update the last_used_id
	_, err = tx.Exec("UPDATE id_ranges SET last_used_id = ? WHERE service = ?", end, service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update range"})
		return
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, IdRangeResponse{Start: start, End: end})
}

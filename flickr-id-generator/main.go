package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var dbs []*sql.DB

func init() {
	// Initialize the two ticket servers (even/odd IDs)
	serverConfigs := []struct {
		dsn      string
		sqlFile  string
		serverID int
	}{
		{"root:root_password@tcp(localhost:3307)/tickets", "even.sql", 1},
		{"root:root_password@tcp(localhost:3308)/tickets", "odd.sql", 2},
	}

	for _, config := range serverConfigs {
		db, err := sql.Open("mysql", config.dsn)
		if err != nil {
			log.Fatalf("Failed to connect to TicketServer%d: %v", config.serverID, err)
		}
		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping TicketServer%d: %v", config.serverID, err)
		}

		// Read and execute the SQL file
		sqlContent, err := ioutil.ReadFile(config.sqlFile)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", config.sqlFile, err)
		}

		commands := strings.Split(string(sqlContent), ";")
		for _, cmd := range commands {
			cmd = strings.TrimSpace(cmd)
			if cmd == "" {
				continue
			}
			_, err := db.Exec(cmd)
			if err != nil {
				log.Fatalf("Failed to execute SQL command on TicketServer%d: %v\nCommand: %s", config.serverID, err, cmd)
			}
		}

		log.Printf("TicketServer%d initialized successfully", config.serverID)
		dbs = append(dbs, db)
	}
}

var serverIndex int = 0

func getServer() *sql.DB {
	// Use round-robin load balancing, random 0 or 1
	serverIndex = (serverIndex + 1) % len(dbs)

	log.Printf("Using TicketServer%d", serverIndex+1)

	return dbs[serverIndex]
}

// GetUniqueID uses round-robin to distribute requests between the two servers
func GetUniqueID() (uint64, error) {
	// Use round-robin load balancing, random 0 or 1
	db := getServer()

	// Execute the query to get a unique ID
	_, err := db.Exec("REPLACE INTO Tickets64 (stub) VALUES ('a')")
	if err != nil {
		return 0, fmt.Errorf("failed to execute REPLACE: %v", err)
	}

	var id uint64
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to get LAST_INSERT_ID: %v", err)
	}

	return id, nil
}

func main() {
	r := gin.Default()

	r.GET("/uniqueid", func(c *gin.Context) {
		id, err := GetUniqueID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// Run the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

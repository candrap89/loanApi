package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/candrap89/loanApi/config"
	"github.com/candrap89/loanApi/handlers"
	"github.com/candrap89/loanApi/queries"
	"github.com/candrap89/loanApi/scheduler"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load database configuration
	cfg, err := config.LoadConfig("service-config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to the database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Database))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize queries and handlers
	userLoanQuery := queries.NewUserLoanQuery(db)
	userLoanHandler := handlers.NewUserLoanHandler(userLoanQuery)

	billingQuery := queries.NewBillingQuery(db)

	// Initialize and start the scheduler
	scheduler := scheduler.NewScheduler(cfg, billingQuery, userLoanQuery)
	scheduler.Start()

	// Define routes
	http.HandleFunc("/user-loan", userLoanHandler.GetUserLoanByCIF)

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

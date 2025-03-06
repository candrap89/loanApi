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
	log.Println("Config loaded successfully:", cfg)

	// Connect to the database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Database))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize queries
	userLoanQuery := queries.NewUserLoanQuery(db)
	billingQuery := queries.NewBillingQuery(db)
	transactionQuery := queries.NewTransactionQuery(db)

	// Initialize scheduler
	scheduler := scheduler.NewScheduler(billingQuery, userLoanQuery, cfg)
	scheduler.Start() // Start the scheduler in the background

	// Initialize handlers
	userLoanHandler := handlers.NewUserLoanHandler(userLoanQuery)
	schedulerHandler := handlers.NewSchedulerHandler(scheduler)
	paymentHandler := handlers.NewPaymentHandler(billingQuery, userLoanQuery, transactionQuery)

	// Define routes
	http.HandleFunc("/user-outstanding", userLoanHandler.GetUserLoanByCIF)
	http.HandleFunc("/delinquents", userLoanHandler.GetDelinquentUsers)
	http.HandleFunc("/trigger-job", schedulerHandler.TriggerJob)
	http.HandleFunc("/payment", paymentHandler.MakePayment)

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

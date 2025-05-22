package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/candrap89/loanApi/config"
	"github.com/candrap89/loanApi/handlers"
	"github.com/candrap89/loanApi/kafka"
	"github.com/candrap89/loanApi/middleware"
	"github.com/candrap89/loanApi/queries"
	"github.com/candrap89/loanApi/scheduler"
	_ "github.com/lib/pq"
)

func main() {
	// Load database configuration
	cfg, err := config.LoadConfig("service-config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Config loaded successfully:", cfg)

	// Connect to the database
	// Connect to the PostgreSQL database
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Database,
	)

	db, err := sql.Open("postgres", connStr)
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

	// Create router and apply middleware to all routes
	router := http.NewServeMux()

	// Create middleware chain
	authMiddleware := middleware.APIKeyMiddleware(cfg.ApiKey)

	// Define routes
	router.Handle("/user-outstanding", authMiddleware(http.HandlerFunc(userLoanHandler.GetUserLoanByCIF)))
	router.Handle("/delinquents", authMiddleware(http.HandlerFunc(userLoanHandler.GetDelinquentUsers)))
	router.Handle("/trigger-job", authMiddleware(http.HandlerFunc(schedulerHandler.TriggerJob)))
	router.Handle("/payment", authMiddleware(http.HandlerFunc(paymentHandler.MakePayment)))
	router.Handle("/user-loan", authMiddleware(http.HandlerFunc(userLoanHandler.CreateUserLoan)))
	router.Handle("/user-vote", authMiddleware(http.HandlerFunc(handlers.GetVoteCountHandler)))

	// Start Kafka consumers
	// object initialization
	consumer := kafka.NewConsumerHandler(userLoanQuery) // Uncomment this line if the function is defined in the kafka package
	go consumer.StartNewProductConsumer()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start the server
	log.Println("Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}

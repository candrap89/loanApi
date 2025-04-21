package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	redis "github.com/candrap89/loanApi/Redis"
	"github.com/candrap89/loanApi/config"
	"github.com/candrap89/loanApi/handlers"
	"github.com/candrap89/loanApi/kafka"
	"github.com/candrap89/loanApi/middleware"
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

	// Initialize Redis
	redisClient := redis.NewRedisClient(redis.RedisConfig{
		Hosts:    cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Connect to the database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Database))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// define routes
	routeConfig := []middleware.RouteAPIKey{
		{
			Route:     "/user-outstanding",
			APIKey:    cfg.APIKey.OutstandingKey,
			RateLimit: cfg.RateLimit.Outstanding,
		},
		{
			Route:     "/delinquents",
			APIKey:    cfg.APIKey.DelinguentsKey,
			RateLimit: cfg.RateLimit.Delinguents,
		},
		{
			Route:     "/payment",
			APIKey:    cfg.APIKey.PaymentKey,
			RateLimit: cfg.RateLimit.Payment,
		},
	}

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

	// Define routes
	//router.Handle("/user-outstanding", authMiddleware(http.HandlerFunc(userLoanHandler.GetUserLoanByCIF)))
	//router.Handle("/delinquents", authMiddleware(http.HandlerFunc(userLoanHandler.GetDelinquentUsers)))
	router.Handle("/trigger-job", http.HandlerFunc(schedulerHandler.TriggerJob))
	//router.Handle("/payment", authMiddleware(http.HandlerFunc(paymentHandler.MakePayment)))
	router.Handle("/user-loan", http.HandlerFunc(userLoanHandler.CreateUserLoan))
	router.Handle("/user-vote", http.HandlerFunc(handlers.GetVoteCountHandler))

	routerProtected := http.NewServeMux()
	routerProtected.HandleFunc("/user-outstanding", userLoanHandler.GetUserLoanByCIF)
	routerProtected.HandleFunc("/delinquents", userLoanHandler.GetDelinquentUsers)
	routerProtected.HandleFunc("/payment", paymentHandler.MakePayment)

	router.Handle("/", middleware.APIKeyMiddleware(routeConfig, redisClient, userLoanQuery)(routerProtected))

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

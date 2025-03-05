package scheduler

import (
	"log"
	"time"

	"github.com/candrap89/loanApi/config"
	"github.com/candrap89/loanApi/models"
	"github.com/candrap89/loanApi/queries"
)

type Scheduler struct {
	Config        *config.Config
	BillingQuery  *queries.BillingQuery
	UserLoanQuery *queries.UserLoanQuery // Add UserLoanQuery to fetch users
}

func NewScheduler(cfg *config.Config, billingQuery *queries.BillingQuery, userLoanQuery *queries.UserLoanQuery) *Scheduler {
	return &Scheduler{
		Config:        cfg,
		BillingQuery:  billingQuery,
		UserLoanQuery: userLoanQuery,
	}
}

func (s *Scheduler) Start() {
	go func() {
		for {
			s.runTask()
			s.sleepUntilNextRun()
		}
	}()
}

func (s *Scheduler) runTask() {
	log.Println("Running scheduler task...")

	// Fetch all users from the user_loan table
	users, err := s.UserLoanQuery.GetAllUsers()
	if err != nil {
		log.Printf("Failed to fetch users: %v", err)
		return
	}

	// Generate billing records for each user
	for _, user := range users {
		billing := models.Billing{
			IDUser:          user.ID,
			BillAmount:      user.Loan * (user.Interest / 100), // Example: Bill amount = loan * interest rate
			PaidStatus:      false,
			LastUpdatedAt:   time.Now(),
			LoanOutstanding: user.LoanOutstanding,
			Week:            int(time.Now().Weekday()),
		}

		// Insert the billing record into the database
		err := s.BillingQuery.InsertBilling(billing)
		if err != nil {
			log.Printf("Failed to insert billing record for user %d: %v", user.ID, err)
		} else {
			log.Printf("Billing record inserted successfully for user %d", user.ID)
		}
	}
}

func (s *Scheduler) sleepUntilNextRun() {
	interval := s.Config.Scheduler.Interval
	nextRun := time.Now()

	switch interval {
	case "minute":
		nextRun = nextRun.Add(time.Minute)
	case "hour":
		nextRun = nextRun.Add(time.Hour)
	case "day":
		nextRun = nextRun.AddDate(0, 0, 1)
	case "week":
		nextRun = nextRun.AddDate(0, 0, 7)
	default:
		log.Fatalf("Invalid scheduler interval: %s", interval)
	}

	if s.Config.Scheduler.Time != "" {
		// Parse the specific time for daily/weekly schedules
		layout := "15:04"
		scheduledTime, err := time.Parse(layout, s.Config.Scheduler.Time)
		if err != nil {
			log.Fatalf("Failed to parse scheduler time: %v", err)
		}

		nextRun = time.Date(
			nextRun.Year(),
			nextRun.Month(),
			nextRun.Day(),
			scheduledTime.Hour(),
			scheduledTime.Minute(),
			0, 0, nextRun.Location(),
		)
	}

	sleepDuration := time.Until(nextRun)
	log.Printf("Next run in %v", sleepDuration)
	time.Sleep(sleepDuration)
}

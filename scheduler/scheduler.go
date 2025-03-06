package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/candrap89/loanApi/config"
	"github.com/candrap89/loanApi/models"
	"github.com/candrap89/loanApi/queries"
)

type Scheduler struct {
	BillingQuery  *queries.BillingQuery
	UserLoanQuery *queries.UserLoanQuery // Add UserLoanQuery
	Config        *config.Config
}

func NewScheduler(billingQuery *queries.BillingQuery, userLoanQuery *queries.UserLoanQuery, cfg *config.Config) *Scheduler {
	return &Scheduler{
		BillingQuery:  billingQuery,
		UserLoanQuery: userLoanQuery,
		Config:        cfg,
	}
}

// the logic for generating and inserting billing records
func (s *Scheduler) RunJob() error {
	log.Println("Running scheduler job...")

	// Fetch all deliquent users
	deliquentUsers, err := s.BillingQuery.GetDeliquentUsers() // Use UserLoanQuery
	if err != nil {
		return fmt.Errorf("failed to fetch deliquent user: %v", err)
	}

	// Fetch all users from the user_loan table
	users, err := s.UserLoanQuery.GetAllUsers() // Use UserLoanQuery
	if err != nil {
		return fmt.Errorf("failed to fetch users: %v", err)
	}

	for _, user := range users {
		billAmount := (user.Loan + (user.Loan * (user.Interest / 100))) / 50

		// Calculate the total bill amount (sum of all unpaid bills for the user)
		totalBillAmount, err := s.BillingQuery.GetTotalUnpaidBillAmount(user.ID)
		if err != nil {
			return fmt.Errorf("failed to get total unpaid bill amount for user %d: %v", user.ID, err)
		}
		totalBillAmount += billAmount // Add the current bill amount to the total

		// Fetch the latest week value for the user
		week, err := s.BillingQuery.GetLatestWeek(user.ID)
		if err != nil {
			return fmt.Errorf("failed to get latest week for user %d: %v", user.ID, err)
		}
		isDeliquent := contains(deliquentUsers, user.ID)

		billing := models.Billing{
			IDUser:          user.ID,
			BillAmount:      billAmount,
			PaidStatus:      false,
			LastUpdatedAt:   time.Now(),
			LoanOutstanding: user.LoanOutstanding,
			Week:            week,
			TotalBillAmount: totalBillAmount,
		}

		// Update the user to deliquent if they have outstanding loan
		err = s.UserLoanQuery.UpdateUserTodeliquent(isDeliquent, user.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch users: %v", err)
		}

		// Insert the billing record into the database
		err = s.BillingQuery.InsertBilling(billing)
		if err != nil {
			return fmt.Errorf("failed to insert billing record for user %d: %v", user.ID, err)
		}
	}

	log.Println("Job completed successfully")
	return nil
}

func contains(deliquentUsers []int, i int) bool {
	for _, id := range deliquentUsers {
		if id == i {
			return true
		}
	}
	return false
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	go func() {
		for {
			if err := s.RunJob(); err != nil {
				log.Printf("Error running job: %v", err)
			}
			s.sleepUntilNextRun()
		}
	}()
}

// sleepUntilNextRun calculates the duration until the next scheduled run based on the config
func (s *Scheduler) sleepUntilNextRun() {
	now := time.Now()
	var nextRun time.Time

	// Parse the scheduled time from the config
	scheduledTime, err := time.Parse("15:04", s.Config.Scheduler.Time)
	if err != nil {
		log.Fatalf("Failed to parse scheduler time: %v", err)
	}

	switch s.Config.Scheduler.Interval {
	case "minute":
		// Run every minute
		nextRun = now.Add(time.Minute)
	case "hour":
		// Run every hour
		nextRun = now.Add(time.Hour)
	case "day":
		// Run daily at the specified time
		nextRun = time.Date(now.Year(), now.Month(), now.Day(), scheduledTime.Hour(), scheduledTime.Minute(), 0, 0, now.Location())
		if now.After(nextRun) {
			nextRun = nextRun.AddDate(0, 0, 1) // Move to the next day if the time has already passed
		}
	case "week":
		// Run weekly at the specified time
		nextRun = time.Date(now.Year(), now.Month(), now.Day(), scheduledTime.Hour(), scheduledTime.Minute(), 0, 0, now.Location())
		if now.After(nextRun) {
			nextRun = nextRun.AddDate(0, 0, 7) // Move to the next week if the time has already passed
		}
	default:
		log.Fatalf("Invalid scheduler interval: %s", s.Config.Scheduler.Interval)
	}

	sleepDuration := time.Until(nextRun)
	log.Printf("Next run in %v", sleepDuration)
	time.Sleep(sleepDuration)
}

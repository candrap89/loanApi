package handlers

import (
	"fmt"
	"net/http"

	"github.com/candrap89/loanApi/queries"
	"github.com/candrap89/loanApi/scheduler"
)

type SchedulerHandler struct {
	BillingQuery *queries.BillingQuery
	Scheduler    *scheduler.Scheduler
}

func NewSchedulerHandler(scheduler *scheduler.Scheduler) *SchedulerHandler {
	return &SchedulerHandler{Scheduler: scheduler}
}

// TriggerJob manually triggers the scheduler job
func (h *SchedulerHandler) TriggerJob(w http.ResponseWriter, r *http.Request) {
	if err := h.Scheduler.RunJob(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to trigger job: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Job triggered successfully"))
}

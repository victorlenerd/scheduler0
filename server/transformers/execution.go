package transformers

import (
	"encoding/json"
	"scheduler0/server/managers/execution"
	"time"
)

// Execution this transformer is used for execution entity type
type Execution struct {
	ID            int64     `json:"id" sql:",pk:notnull"`
	UUID          string    `json:"uuid"`
	JobID         int64     `json:"job_id"`
	JobUUID       string    `json:"job_uuid"`
	StatusCode    string    `json:"status_code"`
	Timeout       uint64    `json:"timeout"`
	Response      string    `json:"response"`
	ExecutionTime uint64    `json:"execution_time"`
	TimeAdded     time.Time `json:"time_added"`
	TimeExecuted  time.Time `json:"time_executed"`
	DateCreated   time.Time `json:"date_created"`
}

// PaginatedExecution this holds meta information for pagination
type PaginatedExecution struct {
	Total  int         `json:"total"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
	Data   []Execution `json:"executions"`
}

// ToJSON returns JSON representation of transformer
func (executionTransformer *Execution) ToJSON() ([]byte, error) {
	data, err := json.Marshal(executionTransformer)
	if err != nil {
		return data, err
	}
	return data, nil
}

// FromJSON extracts content of JSON object into transformer
func (executionTransformer *Execution) FromJSON(body []byte) error {
	if err := json.Unmarshal(body, &executionTransformer); err != nil {
		return err
	}
	return nil
}

// ToManager converts content of transformer into manager
func (executionTransformer *Execution) ToManager() (execution.Manager, error) {
	executionManager := execution.Manager{
		ID:            executionTransformer.ID,
		UUID:          executionTransformer.UUID,
		JobID:         executionTransformer.JobID,
		JobUUID:       executionTransformer.JobUUID,
		StatusCode:    executionTransformer.StatusCode,
		TimeAdded:     executionTransformer.TimeAdded,
		ExecutionTime: executionTransformer.ExecutionTime,
		DateCreated:   executionTransformer.DateCreated,
	}

	return executionManager, nil
}

// FromManager extract content of manager into transformer
func (executionTransformer *Execution) FromManager(executionManager execution.Manager) {
	executionTransformer.ID = executionManager.ID
	executionTransformer.UUID = executionManager.UUID
	executionTransformer.JobID = executionManager.JobID
	executionTransformer.JobUUID = executionManager.JobUUID
	executionTransformer.StatusCode = executionManager.StatusCode
	executionTransformer.DateCreated = executionManager.DateCreated
	executionTransformer.ExecutionTime = executionManager.ExecutionTime
	executionTransformer.TimeAdded = executionManager.TimeAdded
}

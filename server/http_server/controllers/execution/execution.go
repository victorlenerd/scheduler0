package execution

import (
	"net/http"
	"scheduler0/server/http_server/controllers"
	"scheduler0/server/service"
	"scheduler0/utils"
	"strconv"
)

// Controller handlers all incoming http requests for /executions
type Controller controllers.Controller

// List request handler that returns paginated executions result set
func (executionController *Controller) List(w http.ResponseWriter, r *http.Request) {
	executionService := service.ExecutionService{DBConnection: executionController.DBConnection, Ctx: r.Context()}

	offset := 0
	limit := 50

	jobID, err := utils.ValidateQueryString("jobUUID", r)
	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	limitParam, err := utils.ValidateQueryString("limit", r)
	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	offsetParam, err := utils.ValidateQueryString("offset", r)
	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	offset, err = strconv.Atoi(offsetParam)
	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	limit, err = strconv.Atoi(limitParam)
	if err != nil {
		utils.SendJSON(w, err.Error(), false, http.StatusBadRequest, nil)
		return
	}

	executions, getError := executionService.GetAllExecutionsByJobUUID(jobID, offset, limit)

	if getError != nil {
		utils.SendJSON(w, getError.Message, false, http.StatusBadRequest, nil)
	} else {
		utils.SendJSON(w, executions, true, http.StatusOK, nil)
	}
}

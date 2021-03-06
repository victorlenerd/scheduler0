package fixtures

import (
	"errors"
	"github.com/bxcodec/faker/v3"
	"github.com/go-pg/pg"
	"scheduler0/server/managers/job"
	fixtures2 "scheduler0/server/managers/job/fixtures"
	"scheduler0/server/managers/project"
	fixtures3 "scheduler0/server/managers/project/fixtures"
	"scheduler0/utils"
)

// CreateJobFixture creates a project and job for testing
func CreateJobFixture(dbConnection *pg.DB) *job.Manager {
	projectFixture := fixtures3.ProjectFixture{}
	err := faker.FakeData(&projectFixture)
	utils.CheckErr(err)

	projectManager := project.ProjectManager{
		Name:        projectFixture.Name,
		Description: projectFixture.Description,
	}
	_, projectManagerError := projectManager.CreateOne(dbConnection)
	if projectManagerError != nil {
		utils.Error(projectManagerError.Message)
	}

	jobFixture := fixtures2.JobFixture{}
	jobTransformers := jobFixture.CreateNJobTransformers(1)
	jobManager, toManagerError := jobTransformers[0].ToManager()
	if toManagerError != nil {
		utils.Error(toManagerError)
	}

	jobManager.ProjectUUID = projectManager.UUID
	jobManager.ID = projectManager.ID

	_, createJobManagerError := jobManager.CreateOne(dbConnection)
	if createJobManagerError != nil {
		utils.Error(createJobManagerError.Message)
		panic(errors.New(createJobManagerError.Message))
	}

	return &jobManager
}

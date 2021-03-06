package fixtures

import (
	"github.com/bxcodec/faker"
	"scheduler0/server/transformers"
	"scheduler0/utils"
)

// JobFixture job fixture for testing
type JobFixture struct {
	UUID        string `faker:"uuid_hyphenated"`
	ProjectUUID string `faker:"uuid_hyphenated"`
	Description string `faker:"sentence"`
	CronSpec    string
	Data        string `faker:"username"`
	Timezone    string `faker:"timezone"`
	CallbackUrl string `faker:"ipv6"`
	StartDate   string `faker:"timestamp"`
	EndDate     string `faker:"timestamp"`
}

// CreateNJobTransformers create n number of job transformer fixtures for testing
func (jobFixture *JobFixture) CreateNJobTransformers(n int) []transformers.Job {
	jobTransformers := []transformers.Job{}

	for i := 0; i < n; i++ {
		err := faker.FakeData(&jobFixture)
		utils.CheckErr(err)

		jobTransformers = append(jobTransformers, transformers.Job{
			UUID:        jobFixture.UUID,
			ProjectUUID: "",
			Spec:        "* * * * 1",
			Data:        jobFixture.Data,
			CallbackUrl: jobFixture.CallbackUrl,
		})
	}

	return jobTransformers
}

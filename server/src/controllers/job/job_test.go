package job_test

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/victorlenerd/scheduler0/server/src/controllers/job"
	jobTestFixtures "github.com/victorlenerd/scheduler0/server/src/controllers/job/fixtures"
	jobManagerTestFixtures "github.com/victorlenerd/scheduler0/server/src/managers/job/fixtures"
	projectTestFixtures "github.com/victorlenerd/scheduler0/server/src/managers/project/fixtures"
	"github.com/victorlenerd/scheduler0/server/src/transformers"
	"github.com/victorlenerd/scheduler0/server/src/utils"
	"github.com/victorlenerd/scheduler0/server/tests"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var _ = Describe("Job Controller", func() {

	pool := tests.GetTestPool()

	BeforeEach(func() {
		tests.Teardown()
		tests.Prepare()
	})

	Context("TestJobController_CreateOne", func() {

		It("Respond with status 400 if request body does not contain required values", func () {
			jobController := job.JobController{ Pool: pool}
			jobFixture := jobManagerTestFixtures.JobFixture{}
			jobTransformers := jobFixture.CreateNJobTransformers(1)
			jobByte, err := jobTransformers[0].ToJson()
			utils.CheckErr(err)
			jobStr := string(jobByte)

			req, err := http.NewRequest("POST", "/jobs", strings.NewReader(jobStr))
			if err != nil {
				utils.Error("Cannot create http request")
			}

			w := httptest.NewRecorder()
			jobController.CreateJob(w, req)

			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

		It("Respond with status 201 if request body is valid", func () {
			projectManager := projectTestFixtures.CreateProjectManagerFixture()
			projectManager.CreateOne(pool)

			jobFixture := jobManagerTestFixtures.JobFixture{}
			jobTransformers := jobFixture.CreateNJobTransformers(1)
			jobTransformers[0].ProjectUUID = projectManager.UUID
			jobByte, err := jobTransformers[0].ToJson()

			if err != nil {
				utils.Error(fmt.Sprintf("Cannot create job %v", err))
			}

			jobStr := string(jobByte)
			req, err := http.NewRequest("POST", "/jobs", strings.NewReader(jobStr))
			if err != nil {
				utils.Error(fmt.Sprintf("Cannot create job %v", err))
			}

			w := httptest.NewRecorder()

			controller := job.JobController{ Pool: pool}
			controller.CreateJob(w, req)
			body, err := ioutil.ReadAll(w.Body)

			if err != nil {
				utils.Error("Could not read response body %v", err)
			}

			var response map[string]interface{}

			if err = json.Unmarshal(body, &response); err != nil {
				utils.Error(fmt.Sprintf("Could unmarsha json response %v", err))
			}

			if len(response) < 1 {
				utils.Error("Response payload is empty")
			}

			utils.Info(response)

			Expect(w.Code).To(Equal(http.StatusCreated))
		})

	})

	Context("TestJobController_GetAll", func() {
		It("Respond with status 200 and return all created jobs", func() {
			projectTransformers := projectTestFixtures.CreateProjectTransformerFixture()
			projectManager := projectTransformers.ToManager()
			projectManager.CreateOne(pool)
			n := 5

			jobFixture := jobManagerTestFixtures.JobFixture{}
			jobTransformers := jobFixture.CreateNJobTransformers(n)

			for i := 0; i < n; i++ {
				jobManager, err := jobTransformers[i].ToManager()
				if err != nil {
					utils.CheckErr(err)
				}
				jobManager.ProjectUUID = projectManager.UUID
				jobManager.CreateOne(pool)
			}

			req, err := http.NewRequest("GET", "/jobs?offset=0&limit=10&projectUUID="+projectManager.UUID, nil)
			if err != nil {
				utils.Error(fmt.Sprintf("Cannot create http request %v", err))
			}

			w := httptest.NewRecorder()
			controller := job.JobController{ Pool: pool}
			controller.ListJobs(w, req)

			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})


	Context("TestJobController_UpdateOne", func() {

		It("Respond with status 400 if update attempts to change cron spec", func () {
			_, jobManager := jobTestFixtures.CreateJobAndProjectManagerFixture(pool)
			jobTransformer := transformers.Job{}
			jobTransformer.FromManager(jobManager)
			jobTransformer.CronSpec = "* * 3 * *"
			jobByte, err := jobTransformer.ToJson()
			utils.CheckErr(err)
			jobStr := string(jobByte)
			req, err := http.NewRequest("PUT", "/jobs/"+jobTransformer.UUID, strings.NewReader(jobStr))
			if err != nil {
				utils.Error(fmt.Sprintf("Cannot create http request %v", err))
			}

			w := httptest.NewRecorder()
			controller := job.JobController{ Pool: pool}

			controller.UpdateJob(w, req)
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})


		It("Respond with status 200 if update body is valid", func () {
			_, jobManager := jobTestFixtures.CreateJobAndProjectManagerFixture(pool)
			jobTransformer := transformers.Job{}
			jobTransformer.FromManager(jobManager)
			jobTransformer.Description = "some job description"
			jobByte, err := jobTransformer.ToJson()
			jobStr := string(jobByte)

			req, err := http.NewRequest("PUT", "/jobs/"+jobTransformer.UUID, strings.NewReader(jobStr))

			if err != nil {
				utils.Error(fmt.Sprintf("Cannot create http request %v", err))
			}

			w := httptest.NewRecorder()
			controller := job.JobController{ Pool: pool}
			router := mux.NewRouter()
			router.HandleFunc("/jobs/{uuid}", controller.UpdateJob)
			router.ServeHTTP(w, req)

			_, err = ioutil.ReadAll(w.Body)
			if err != nil {
				utils.Error(fmt.Sprintf("Cannot create http request %v", err))
			}

			Expect(w.Code).To(Equal(http.StatusOK))
		})

	})

	It("TestJobController_DeleteOne", func () {
		_, jobManager := jobTestFixtures.CreateJobAndProjectManagerFixture(pool)

		req, err := http.NewRequest("DELETE", "/jobs/"+jobManager.UUID, nil)
		if err != nil {
			utils.Error(fmt.Sprintf("cannot create request to delete job %v", err))
		}

		w := httptest.NewRecorder()
		controller := job.JobController{ Pool: pool}

		router := mux.NewRouter()
		router.HandleFunc("/jobs/{uuid}", controller.DeleteJob)
		router.ServeHTTP(w, req)

		if err != nil {
			utils.Error("Cannot create http request %v", err)
		}

		Expect(w.Code).To(Equal(http.StatusNoContent))
	})
})

func TestJob_Controller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Job Controller Suite")
}
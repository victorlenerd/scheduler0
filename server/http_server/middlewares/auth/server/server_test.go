package server_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"scheduler0/server/db"
	"scheduler0/server/http_server/middlewares/auth"
	"scheduler0/server/http_server/middlewares/auth/server"
	"scheduler0/server/managers/credential"
	"scheduler0/server/managers/credential/fixtures"
	"scheduler0/server/service"
	"scheduler0/utils"
	"testing"
)

var _ = Describe("Server Side Auth Test", func() {

	BeforeEach(func() {
		db.Teardown()
		db.Prepare()
	})

	It("Should identify request from server side", func() {
		req, err := http.NewRequest("POST", "/", nil)
		Expect(err).To(BeNil())

		dbConnection := db.GetTestDBConnection()

		credentialService := service.Credential{
			DBConnection: dbConnection,
		}

		credentialFixture := fixtures.CredentialFixture{}
		credentialTransformers := credentialFixture.CreateNCredentialTransformer(1)
		credentialTransformer := credentialTransformers[0]

		credentialTransformer.Platform = credential.ServerPlatform

		_, createError := credentialService.CreateNewCredential(credentialTransformer)
		if createError != nil {
			utils.Error(fmt.Sprintf("Error: %v", createError.Message))
		}

		req.Header.Set(auth.APIKeyHeader, credentialTransformer.ApiKey)
		req.Header.Set(auth.SecretKeyHeader, credentialTransformer.ApiSecret)

		Expect(server.IsServerClient(req)).To(BeTrue())
	})

	It("Should not identify request from non server side", func() {
		req, err := http.NewRequest("POST", "/", nil)
		Expect(err).To(BeNil())

		dbConnection := db.GetTestDBConnection()

		credentialService := service.Credential{
			DBConnection: dbConnection,
		}

		credentialFixture := fixtures.CredentialFixture{}
		credentialTransformers := credentialFixture.CreateNCredentialTransformer(1)

		_, createError := credentialService.CreateNewCredential(credentialTransformers[0])
		if createError != nil {
			utils.Warn(fmt.Sprintf("Error: %v", createError.Message))
		}

		credentialTransformer := credentialTransformers[0]

		req.Header.Set(auth.APIKeyHeader, credentialTransformer.ApiKey)
		req.Header.Set(auth.AndroidPackageIDHeader, credentialTransformer.AndroidPackageNameRestriction)
		Expect(server.IsServerClient(req)).ToNot(BeTrue())

		req.Header.Set(auth.APIKeyHeader, credentialTransformer.ApiKey)
		req.Header.Set(auth.IOSBundleHeader, credentialTransformer.IOSBundleIDRestriction)
		Expect(server.IsServerClient(req)).ToNot(BeTrue())

		req.Header.Set(auth.APIKeyHeader, credentialTransformer.ApiKey)
		Expect(server.IsServerClient(req)).ToNot(BeTrue())
	})

	It("Should identify authorized request from server clients", func() {
		req, err := http.NewRequest("POST", "/", nil)
		Expect(err).To(BeNil())

		dbConnection := db.GetTestDBConnection()

		credentialService := service.Credential{
			DBConnection: dbConnection,
		}

		credentialFixture := fixtures.CredentialFixture{}
		credentialTransformers := credentialFixture.CreateNCredentialTransformer(1)
		credentialTransformer := credentialTransformers[0]

		credentialTransformer.Platform = credential.ServerPlatform

		credentialManagerUUID, createError := credentialService.CreateNewCredential(credentialTransformer)
		if createError != nil {
			utils.Error(fmt.Sprintf("Error: %v", createError.Message))
		}

		updatedCredentialTransformer, getError := credentialService.FindOneCredentialByUUID(credentialManagerUUID)
		if getError != nil {
			utils.Error(fmt.Sprintf("Error: %v", getError.Error()))
		}

		if updatedCredentialTransformer == nil {
			Expect(updatedCredentialTransformer).ToNot(BeNil())
			return
		}

		req.Header.Set(auth.APIKeyHeader, updatedCredentialTransformer.ApiKey)
		req.Header.Set(auth.SecretKeyHeader, updatedCredentialTransformer.ApiSecret)

		Expect(server.IsAuthorizedServerClient(req, dbConnection)).To(BeTrue())
	})
})

func TestServerSideAuth_Middleware(t *testing.T) {
	utils.SetTestScheduler0Configurations()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Side Auth Test")
}

package common_test

import (
	"api-gw/internal/tenant/registration"
	"api-gw/pkg/client"
	"api-gw/pkg/common"
	"api-gw/pkg/config"
	"api-gw/pkg/envoy"
	"context"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/golang/mock/gomock"
	"github.com/labstack/gommon/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	reg_svc_mock "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/mocks/pkg/registration-service/global"
	reg_svc "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/registration-service/global"
	apinexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/api.nexus.vmware.com/v1"
	common_nexus_vmware_com "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/common.nexus.vmware.com/v1"
	confignexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/config.nexus.vmware.com/v1"
	runtimenexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/runtime.nexus.vmware.com/v1"
	tenant_config_v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantconfig.nexus.vmware.com/v1"
	v1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/tenantruntime.nexus.vmware.com/v1"
	userv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/user.nexus.vmware.com/v1"
	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Common Function tests", func() {
	AfterSuite(func() {
		envoy.XDSListener.Close()
	})

	It("Should test convertSKUtoLicense", func() {
		config.SKUConfig = &config.SKUMap{
			SKU: map[string][]string{
				"advance": {
					"testadvance1",
					"testadvance2",
				},
			},
		}
		//Verify if sku to license is parsed correctly
		license := common.ConvertProductIDtoLicense("testadvance1")
		Expect(license).To(Equal("advance"))

		license = common.ConvertProductIDtoLicense("testadvance")
		Expect(license).To(Equal(""))
	})

	It("Should test Setting CSP environment variables and verifyPermissions", func() {
		err := common.SetCSPVariables()
		Expect(err).NotTo(BeNil())

		client.NexusClient = nexus_client.NewFakeClient()
		_, err = common.GetConfigNode(client.NexusClient, "default")
		Expect(err).To(HaveOccurred())

		_, err = common.GetRuntimeNode(client.NexusClient, "default")
		Expect(err).To(HaveOccurred())

		_, err = common.CheckTenantIfExists(client.NexusClient, "test")
		Expect(err).To(HaveOccurred())

		err = common.DeleteUserObject(client.NexusClient, "test")
		Expect(err).To(HaveOccurred())

		err = os.Setenv("CSP_PERMISSION_NAME", "external/8d190d7b-ebb4-4fc9-b4e9-fb4b14148e50/staging-1e")
		Expect(err).To(BeNil())

		err = common.SetCSPVariables()
		Expect(err).To(BeNil())
		Expect(common.CSP_SERVICE_ID).To(Equal("8d190d7b-ebb4-4fc9-b4e9-fb4b14148e50"))
		Expect(common.CSP_SERVICE_NAME).To(Equal("staging-1e"))

		common.SetCSPPermissionOrg()

		token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2ODQzMTk4NTEsImV4cCI6MTcxNTg1NTg1MSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJwZXJtcyI6WyJleHRlcm5hbC84ZDE5MGQ3Yi1lYmI0LTRmYzktYjRlOS1mYjRiMTQxNDhlNTAvc3RhZ2luZy0xZTp1c2VyIiwiY3NwOm1lbWJlciJdfQ.kmxSP_s9ylwnHZzBJGPa7nr0bjG6rYvgDjqgaszftuw"
		claims, _, _ := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
		hasAccess := common.VerifyPermissions(token, claims.Claims, common.Permissions)
		Expect(hasAccess).To(BeTrue())

		token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2ODQzMTk4NTEsImV4cCI6MTcxNTg1NTg1MSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSIsIkdpdmVuTmFtZSI6IkpvaG5ueSIsIlN1cm5hbWUiOiJSb2NrZXQiLCJFbWFpbCI6Impyb2NrZXRAZXhhbXBsZS5jb20iLCJwZXJtcyI6WyJjc3A6bWVtYmVyIiwib3JnIl19.FLKgJdkAOmAS-phE35UY--5fH1_xuwIVzTtzrgLgbJc"
		claims, _, _ = new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
		hasAccess = common.VerifyPermissions(token, claims.Claims, common.Permissions)
		Expect(hasAccess).To(BeFalse())

		err = os.Setenv("CSP_SERVICE_OWNER_TOKEN", "test")
		Expect(err).To(BeNil())
		Expect(common.GetCSPServiceOwnerToken()).To(Equal("test"))

		sdurl := common.GenerateServiceDefinitionURL("http://localhost", "test")
		Expect(sdurl).To(Equal("http://localhost/csp/gateway/slc/api/v2/orgs/test/services"))

	})

	It("should test create cookie method", func() {
		cookie := common.CreateCookie("test", "value", time.Time{})
		Expect(cookie.Name).To(Equal("test"))
		Expect(cookie.Value).To(Equal("value"))
	})

	It("User related common methods", func() {
		client.NexusClient = nexus_client.NewFakeClient()
		_, err := client.NexusClient.Api().CreateNexusByName(context.TODO(), &apinexusv1.Nexus{
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
		})
		Expect(err).NotTo(HaveOccurred())

		_, err = common.GetConfigNode(client.NexusClient, "default")
		Expect(err).NotTo(BeNil())

		_, err = client.NexusClient.Config().CreateConfigByName(context.TODO(), &confignexusv1.Config{
			ObjectMeta: metav1.ObjectMeta{
				Name: "943ea6107388dc0d02a4c4d861295cd2ce24d551",
				Labels: map[string]string{
					common.DISPLAY_NAME: "default",
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		_, err = client.NexusClient.Runtime().CreateRuntimeByName(context.TODO(), &runtimenexusv1.Runtime{
			ObjectMeta: metav1.ObjectMeta{
				Name: "e817339e4e7bf29fa47ca62dd272b44282d271b8",
				Labels: map[string]string{
					common.DISPLAY_NAME: "default",
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		// check if tenant creates with valid SKU
		err = common.CreateTenantIfNotExists(client.NexusClient, "test", "advance")
		Expect(err).To(BeNil())

		// check if tenant creates fails with invalid SKU
		err = common.CreateTenantIfNotExists(client.NexusClient, "test2", "advanced")
		Expect(err).ToNot(BeNil())

		runtimeNode, err := common.GetRuntimeNode(client.NexusClient, "default")
		Expect(err).To(BeNil())

		_, err = runtimeNode.AddTenant(context.TODO(), &v1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: v1.TenantSpec{
				Namespace: "test",
			},
			Status: v1.TenantNexusStatus{
				AppStatus: v1.TenantStatus{
					InstalledApplications: common_nexus_vmware_com.ApplicationStatus{
						NexusApps: map[string]common_nexus_vmware_com.NexusApp{
							"nexus-tenant-runtime": {
								OamApp: common_nexus_vmware_com.OamApp{
									Components: map[string]common_nexus_vmware_com.ComponentDefinition{
										"test.nexus-tenant-runtime": {
											Name:   "test.nexus-tenant-runtime.tsm-tenant-runtime",
											Sync:   "OutOfSync",
											Health: "running",
										},
									},
								},
							},
							"tsm-tenant": {
								OamApp: common_nexus_vmware_com.OamApp{
									Components: map[string]common_nexus_vmware_com.ComponentDefinition{
										"test.tsm-tenant-runtime": {
											Name:   "test.tsm-tenant-runtime",
											Sync:   "Synced",
											Health: "running",
										},
									},
								},
							},
						},
					},
				},
			},
		})
		Expect(err).To(BeNil())

		_, found, err := common.CheckTenantRuntimeIfExists(client.NexusClient, "test")
		Expect(found).To(Equal(true))
		Expect(err).To(BeNil())

		tenantS, err := runtimeNode.GetTenant(context.TODO(), "test")
		Expect(err).To(BeNil())
		status, message := common.GetTenantStatus(tenantS.Status.AppStatus)
		Expect(string(status)).To(Equal(string(common.CREATING)))
		common.AddTenantState("test", common.TenantState{
			Status:        status,
			Message:       message,
			CreationStart: time.Now().Format(time.RFC3339Nano),
		})
		stateObj, ok := common.GetTenantState("test")
		Expect(ok).To(BeTrue())
		Expect(stateObj.Status).To(Equal(status))

		httpStatus, _ := common.GetServableTenantStatus("test")
		Expect(httpStatus).To(Equal(503))

		tenantS.SetAppStatus(context.TODO(), &v1.TenantStatus{
			InstalledApplications: common_nexus_vmware_com.ApplicationStatus{
				NexusApps: map[string]common_nexus_vmware_com.NexusApp{
					"nexus-tenant-runtime": {
						OamApp: common_nexus_vmware_com.OamApp{
							Components: map[string]common_nexus_vmware_com.ComponentDefinition{
								"test.nexus-tenant-runtime": {
									Name:   "test.nexus-tenant-runtime.tsm-tenant-runtime",
									Sync:   "Synced",
									Health: "running",
								},
							},
						},
					},
					"tsm-tenant": {
						OamApp: common_nexus_vmware_com.OamApp{
							Components: map[string]common_nexus_vmware_com.ComponentDefinition{
								"test.tsm-tenant-runtime": {
									Name:   "test.tsm-tenant-runtime",
									Sync:   "Synced",
									Health: "running",
								},
							},
						},
					},
				},
			},
		})

		tenantS, err = runtimeNode.GetTenant(context.TODO(), "test")
		status, message = common.GetTenantStatus(tenantS.Status.AppStatus)
		Expect(string(status)).To(Equal(string(common.CREATED)))
		common.AddTenantState("test", common.TenantState{
			Status:        status,
			Message:       message,
			CreationStart: time.Now().Format(time.RFC3339Nano),
		})

		stateObj, ok = common.GetTenantState("test")
		Expect(ok).To(BeTrue())
		Expect(stateObj.Status).To(Equal(status))

		httpStatus, _ = common.GetServableTenantStatus("test")
		Expect(httpStatus).To(Equal(200))

		err = common.CreateUser(client.NexusClient, "test", userv1.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: "user1",
			},
			Spec: userv1.UserSpec{
				Username: "user1",
				Password: "password",
				TenantId: "test",
			},
		})
		Expect(err).To(BeNil())

		err = common.CreateUser(client.NexusClient, "test2", userv1.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: "user2",
			},
			Spec: userv1.UserSpec{
				Username: "user2",
				Password: "password",
				TenantId: "test2",
			},
		})
		Expect(err).NotTo(BeNil())
		Expect(common.UserMap["user1"].Username).To(Equal("user1"))

		configObj, _ := common.GetConfigNode(client.NexusClient, "default")
		_, err = configObj.AddUser(context.Background(), &userv1.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: "newuser",
			},
			Spec: userv1.UserSpec{
				Username: "newuser",
				Password: "password",
				TenantId: "test",
			},
		})
		Expect(err).To(BeNil())

		_, found = common.GetUser("newuser")
		Expect(found).To(BeFalse())

		common.InitAdminDatamodelCache()
		_, found = common.GetUser("newuser")
		Expect(found).To(BeTrue())

		common.AddUser("testuser", userv1.UserSpec{
			Username: "testuser",
			Password: "password",
			TenantId: "test",
		})
		_, found = common.GetUser("testuser")
		Expect(found).To(BeTrue())

		common.DeleteUser("testuser")
		_, found = common.GetUser("testuser")
		Expect(found).NotTo(BeTrue())

		token := "user1:password"
		username := common.GetUserNameFromToken(token)
		Expect(username).To(Equal("user1"))

	})

	It("should test config load", func() {
		configObj, err := config.LoadConfig("../../test/config/api-gw-config.yaml")
		Expect(err).To(BeNil())
		Expect(configObj.EnableNexusRuntime).To(Equal(true))

		_, err = config.LoadSKUConfig("../../test/config/sku-configmap.yaml")
		Expect(err).To(BeNil())

		// Error condition
		_, err = config.LoadConfig("test/config/api-gw-config.yaml")
		Expect(err).NotTo(BeNil())

		_, err = config.LoadSKUConfig("test/config/sku-configmap.yaml")
		Expect(err).NotTo(BeNil())

		_, err = config.LoadStaticUrlsConfig("test/config/sku-configmap.yaml")
		Expect(err).NotTo(BeNil())

		_, err = config.LoadStaticUrlsConfig("../../test/config/staticUrls.yaml")
		Expect(err).To(BeNil())

	})

	It("should check commond methods for status and displayName", func() {
		servUnJson := common.GetServiceUnavailableJson("Error in tenant", "STATE_CREATION", "STATE_IN_PROGRESS", "2022-03-17T12:52:20.998Z")
		Expect(servUnJson).To(Equal(map[string]interface{}{
			"featureFlag": "firstTimeExperience",
			"error":       "Error in tenant",
			"details": map[string]string{
				"state":         "STATE_CREATION",
				"status":        "STATE_IN_PROGRESS",
				"message":       "Error in tenant",
				"creationStart": "2022-03-17T12:52:20.998Z",
			},
		}))

		common.AddTenantDisplayName("test", "testing")
		name, ok := common.GetTenantDisplayName("test")
		Expect(ok).To(BeTrue())
		Expect(name).To(Equal("testing"))

		common.DeleteTenantDisplayName("test")
		_, ok = common.GetTenantDisplayName("test")
		Expect(ok).To(BeFalse())

		config.GlobalStaticRouteConfig = &config.GlobalStaticRoutes{
			Suffix: []string{"js", "css", "png"},
			Prefix: []string{"/home", "/allspark-static"},
		}

		envoy.Init(nil, nil, nil, logrus.Level(log.Level()))
		snap, err := envoy.GenerateNewSnapshot(nil, nil, nil, nil)
		Expect(snap).NotTo(BeNil())
		Expect(err).To(BeNil())

		ctrl := gomock.NewController(GinkgoT())
		regClient := reg_svc_mock.NewMockGlobalRegistrationClient(ctrl)

		gomock.InOrder(
			regClient.EXPECT().RegisterTenant(gomock.Any(), gomock.Any()).Return(&reg_svc.TenantResponse{
				Code: 0,
			}, nil),
		)

		gomock.InOrder(
			regClient.EXPECT().UnregisterTenant(gomock.Any(), gomock.Any()).Return(&reg_svc.TenantResponse{
				Code: 0,
			}, nil),
		)

		err = registration.AddTenantToSystem(tenant_config_v1.Tenant{
			ObjectMeta: metav1.ObjectMeta{
				Name: "8088123",
				Labels: map[string]string{
					common.DISPLAY_NAME: "test",
				},
			},
			Spec: tenant_config_v1.TenantSpec{
				Name: "test",
				Skus: []string{"advance"},
			},
		}, regClient)
		Expect(err).NotTo(HaveOccurred())

		err = common.RegisterTenant(regClient, "test", reg_svc.TenantRequest_License(common.AvailableSkus["advance"]))
		Expect(err).NotTo(HaveOccurred())

		err = common.UnregisterTenant(regClient, "test", reg_svc.TenantRequest_License(common.AvailableSkus["advance"]))
		Expect(err).NotTo(HaveOccurred())

		gomock.InOrder(
			regClient.EXPECT().RegisterTenant(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("could not create tenant")),
		)

		gomock.InOrder(
			regClient.EXPECT().UnregisterTenant(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("could not create tenant")),
		)

		err = common.RegisterTenant(regClient, "test", reg_svc.TenantRequest_License(common.AvailableSkus["advance"]))
		Expect(err).To(HaveOccurred())

		err = common.UnregisterTenant(regClient, "test", reg_svc.TenantRequest_License(common.AvailableSkus["advance"]))
		Expect(err).To(HaveOccurred())

	})
})

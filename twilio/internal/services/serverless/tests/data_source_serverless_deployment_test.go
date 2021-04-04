package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/RJPearson94/terraform-provider-twilio/twilio/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var deploymentDataSourceName = "twilio_serverless_deployment"

func TestAccDataSourceTwilioServerlessDeployment_basic(t *testing.T) {
	stateDataSourceName := fmt.Sprintf("data.%s.deployment", deploymentDataSourceName)
	uniqueName := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTwilioServerlessDeployment_basic(uniqueName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(stateDataSourceName, "id"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "account_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "service_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "environment_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "build_sid"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "date_created"),
					resource.TestCheckNoResourceAttr(stateDataSourceName, "date_updated"),
					resource.TestCheckResourceAttrSet(stateDataSourceName, "url"),
				),
			},
		},
	})
}

func TestAccDataSourceTwilioServerlessDeployment_invalidServiceSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTwilioServerlessDeployment_invalidServiceSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of service_sid to match regular expression "\^ZS\[0-9a-fA-F\]\{32\}\$", got service_sid`),
			},
		},
	})
}

func TestAccDataSourceTwilioServerlessDeployment_invalidEnvironmentSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTwilioServerlessDeployment_invalidEnvironmentSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of environment_sid to match regular expression "\^ZE\[0-9a-fA-F\]\{32\}\$", got environment_sid`),
			},
		},
	})
}

func TestAccDataSourceTwilioServerlessDeployment_invalidSid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.PreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceTwilioServerlessDeployment_invalidSid(),
				ExpectError: regexp.MustCompile(`(?s)expected value of sid to match regular expression "\^ZD\[0-9a-fA-F\]\{32\}\$", got sid`),
			},
		},
	})
}

func testAccDataSourceTwilioServerlessDeployment_basic(uniqueName string) string {
	return fmt.Sprintf(`
resource "twilio_serverless_service" "service" {
  unique_name   = "service-%[1]s"
  friendly_name = "test"
}

resource "twilio_serverless_function" "function" {
  service_sid       = twilio_serverless_service.service.sid
  friendly_name     = "test"
  content           = <<EOF
exports.handler = function (context, event, callback) {
	callback(null, "Hello World");
};
EOF
  content_type      = "application/javascript"
  content_file_name = "helloWorld.js"
  path              = "/test-function"
  visibility        = "private"
}

resource "twilio_serverless_build" "build" {
  service_sid = twilio_serverless_service.service.sid
  function_version {
    sid = twilio_serverless_function.function.latest_version_sid
  }
  polling {
    enabled = true
  }
}

resource "twilio_serverless_environment" "environment" {
  service_sid = twilio_serverless_service.service.sid
  unique_name = "%[1]s"
}

resource "twilio_serverless_deployment" "deployment" {
  service_sid     = twilio_serverless_service.service.sid
  environment_sid = twilio_serverless_environment.environment.sid
  build_sid       = twilio_serverless_build.build.sid
}

data "twilio_serverless_deployment" "deployment" {
  service_sid     = twilio_serverless_deployment.deployment.service_sid
  environment_sid = twilio_serverless_deployment.deployment.environment_sid
  sid             = twilio_serverless_deployment.deployment.sid
}
`, uniqueName)
}

func testAccDataSourceTwilioServerlessDeployment_invalidServiceSid() string {
	return `
data "twilio_serverless_deployment" "deployment" {
  service_sid     = "service_sid"
  environment_sid = "ZEaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  sid             = "ZDaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`
}

func testAccDataSourceTwilioServerlessDeployment_invalidEnvironmentSid() string {
	return `
data "twilio_serverless_deployment" "deployment" {
  service_sid     = "ZSaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  environment_sid = "environment_sid"
  sid             = "ZDaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`
}

func testAccDataSourceTwilioServerlessDeployment_invalidSid() string {
	return `
data "twilio_serverless_deployment" "deployment" {
  service_sid     = "ZSaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  environment_sid = "ZEaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
  sid             = "sid"
}
`
}

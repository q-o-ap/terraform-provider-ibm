/* IBM Confidential
*  Object Code Only Source Materials
*  5747-SM3
*  (c) Copyright IBM Corp. 2017,2021
*
*  The source code for this program is not published or otherwise divested
*  of its trade secrets, irrespective of what has been deposited with the
*  U.S. Copyright Office. */

package ibm

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/IBM-Cloud/bluemix-go/models"
)

func TestAccIBMCertificateManagerOrder_Import(t *testing.T) {
	var conf models.CertificateInfo
	orderName := fmt.Sprintf("tf-acc-test1-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	cmsName := fmt.Sprintf("tf-acc-test1-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCertificateManagerOrderDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCertificateManagerOrder_basic(cmsName, orderName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCMOrderExists("ibm_certificate_manager_order.cert", conf),
					resource.TestCheckResourceAttr("ibm_certificate_manager_order.cert", "name", orderName),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_certificate_manager_order.cert",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"dns_provider_instance_crn", "renew_certificate"},
			},
		},
	})
}
func TestAccIBMCertificateManagerOrder_Basic(t *testing.T) {
	var conf models.CertificateInfo
	orderName := fmt.Sprintf("tf-acc-test1-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	updatedName := fmt.Sprintf("tf-acc-test1-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	cmsName := fmt.Sprintf("tf-acc-test1-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMCertificateManagerOrderDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMCertificateManagerOrder_basic(cmsName, orderName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCMOrderExists("ibm_certificate_manager_order.cert", conf),
					resource.TestCheckResourceAttr("ibm_certificate_manager_order.cert", "name", orderName),
					resource.TestCheckResourceAttr("ibm_certificate_manager_order.cert", "auto_renew_enabled", "false"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMCertificateManagerOrder_Update(cmsName, updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMCMOrderExists("ibm_certificate_manager_order.cert", conf),
					resource.TestCheckResourceAttr("ibm_certificate_manager_order.cert", "name", updatedName),
					resource.TestCheckResourceAttr("ibm_certificate_manager_order.cert", "auto_renew_enabled", "true"),
					resource.TestCheckResourceAttr("ibm_certificate_manager_order.cert", "renew_certificate", "true"),
				),
			},
		},
	})
}

func testAccCheckIBMCertificateManagerOrderDestroy(s *terraform.State) error {
	time.Sleep(100 * time.Second)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_certificate_manager_order" {
			continue
		}
		certID := rs.Primary.ID
		cmClient, err := testAccProvider.Meta().(ClientSession).CertificateManagerAPI()
		if err != nil {
			return err
		}
		certAPI := cmClient.Certificate()
		_, err = certAPI.GetCertData(certID)

		if err != nil && !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "412") {
			return fmt.Errorf("Error checking if instance (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}
	return nil

}

func testAccCheckIBMCertificateManagerOrder_basic(cmsName, orderName string) string {
	return fmt.Sprintf(`
	resource "ibm_resource_instance" "cm" {
		name                = "%s"
		location            = "us-south"
		service             = "cloudcerts"
		plan                = "free"
	}
	data "ibm_resource_group" "web_group" {
		name = "default"
	}
	data "ibm_cis" "instance" {
		name              = "Terraform-Test-CIS"
		resource_group_id = data.ibm_resource_group.web_group.id
	}
	data "ibm_cis_domain" "web_domain" {
		cis_id = data.ibm_cis.instance.id
		domain = "cis-test-domain.com"
	}
	resource "ibm_certificate_manager_order" "cert" {
		certificate_manager_instance_id = ibm_resource_instance.cm.id
		name                            = "%s"
		description                     = "test description"
		domains                         = ["cis-test-domain.com"]
		rotate_keys                     = false
		domain_validation_method        = "dns-01"
		dns_provider_instance_crn       = data.ibm_cis.instance.id
		auto_renew_enabled 				= false
	  }
	  
	  `, cmsName, orderName)
}
func testAccCheckIBMCertificateManagerOrder_Update(cmsName, updatedName string) string {
	return fmt.Sprintf(`
	resource "ibm_resource_instance" "cm" {
		name                = "%s"
		location            = "us-south"
		service             = "cloudcerts"
		plan                = "free"
	}
	data "ibm_resource_group" "web_group" {
		name = "default"
	}
	data "ibm_cis" "instance" {
		name              = "Terraform-Test-CIS"
		resource_group_id = data.ibm_resource_group.web_group.id
	}
	data "ibm_cis_domain" "web_domain" {
		cis_id = data.ibm_cis.instance.id
		domain = "cis-test-domain.com"
	}
	resource "ibm_certificate_manager_order" "cert" {
		certificate_manager_instance_id = ibm_resource_instance.cm.id
		name                            = "%s"
		description                     = "test description"
		domains                         = ["cis-test-domain.com"]
		rotate_keys                     = false
		domain_validation_method        = "dns-01"
		dns_provider_instance_crn       = data.ibm_cis.instance.id
		auto_renew_enabled 				= true
		renew_certificate = true
	  }
	  
	  `, cmsName, updatedName)
}

func testAccCheckIBMCMOrderExists(n string, obj models.CertificateInfo) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		cmClient, err := testAccProvider.Meta().(ClientSession).CertificateManagerAPI()
		if err != nil {
			return err
		}
		certID := rs.Primary.ID

		crt, err := cmClient.Certificate().GetMetaData(certID)
		if err != nil {
			return err
		}

		obj = crt
		return nil
	}
}

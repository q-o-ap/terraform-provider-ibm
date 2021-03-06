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
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccIBMSecurityGroup_basic(t *testing.T) {
	var sg datatypes.Network_SecurityGroup

	name1 := fmt.Sprintf("terraformsguat_create_step_name_%d", acctest.RandIntRange(10, 100))
	desc1 := fmt.Sprintf("terraformsguat_create_step_desc_%d", acctest.RandIntRange(10, 100))
	name2 := fmt.Sprintf("terraformsguat_create_step_name_%d", acctest.RandIntRange(10, 100))
	desc2 := fmt.Sprintf("terraformsguat_create_step_desc_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMSecurityGroupConfig(name1, desc1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMSecurityGroupExists("ibm_security_group.testacc_security_group", &sg),
					resource.TestCheckResourceAttr(
						"ibm_security_group.testacc_security_group", "name", name1),
					resource.TestCheckResourceAttr(
						"ibm_security_group.testacc_security_group", "description", desc1),
				),
			},

			{
				Config: testAccCheckIBMSecurityGroupConfig(name2, desc2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMSecurityGroupExists("ibm_security_group.testacc_security_group", &sg),
					resource.TestCheckResourceAttr(
						"ibm_security_group.testacc_security_group", "name", name2),
					resource.TestCheckResourceAttr(
						"ibm_security_group.testacc_security_group", "description", desc2),
				),
			},
		},
	})
}

func testAccCheckIBMSecurityGroupDestroy(s *terraform.State) error {
	service := services.GetNetworkSecurityGroupService(testAccProvider.Meta().(ClientSession).SoftLayerSession())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_security_group" {
			continue
		}

		sgID, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the key
		_, err := service.Id(sgID).GetObject()

		if err == nil {
			return fmt.Errorf("Security Group %d still exists", sgID)
		}
	}

	return nil
}

func testAccCheckIBMSecurityGroupExists(n string, sg *datatypes.Network_SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		sgID, _ := strconv.Atoi(rs.Primary.ID)

		service := services.GetNetworkSecurityGroupService(testAccProvider.Meta().(ClientSession).SoftLayerSession())
		foundSG, err := service.Id(sgID).GetObject()

		if err != nil {
			return err
		}

		if strconv.Itoa(int(*foundSG.Id)) != rs.Primary.ID {
			return fmt.Errorf("Record %d not found", sgID)
		}

		*sg = foundSG

		return nil
	}
}

func testAccCheckIBMSecurityGroupConfig(name, description string) string {
	return fmt.Sprintf(`
resource "ibm_security_group" "testacc_security_group" {
    name = "%s"
    description = "%s"
}`, name, description)

}

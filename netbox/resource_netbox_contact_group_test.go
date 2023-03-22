package netbox

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/fbreckle/go-netbox/netbox/client"
	"github.com/fbreckle/go-netbox/netbox/client/tenancy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetboxContactGroup_basic(t *testing.T) {

	testSlug := "t_grp_basic"
	testName := testAccGetTestName(testSlug)
	randomSlug := testAccGetTestName(testSlug)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name = "%s"
  slug = "%s"
}`, testName, randomSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", randomSlug),
				),
			},
			{
				ResourceName:      "netbox_contact_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxContactGroup_defaultSlug(t *testing.T) {

	testSlug := "contact_defSlug"
	testName := testAccGetTestName(testSlug)
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name = "%s"
}`, testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", testName),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", getSlug(testName)),
				),
			},
		},
	})
}

func init() {
	resource.AddTestSweepers("netbox_contact_group", &resource.Sweeper{
		Name:         "netbox_contact_group",
		Dependencies: []string{},
		F: func(region string) error {
			m, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("Error getting client: %s", err)
			}
			api := m.(*client.NetBoxAPI)
			params := tenancy.NewTenancyContactGroupsListParams()
			res, err := api.Tenancy.TenancyContactGroupsList(params, nil)
			if err != nil {
				return err
			}
			for _, contact := range res.GetPayload().Results {
				if strings.HasPrefix(*contact.Name, testPrefix) {
					deleteParams := tenancy.NewTenancyContactGroupsDeleteParams().WithID(contact.ID)
					_, err := api.Tenancy.TenancyContactGroupsDelete(deleteParams, nil)
					if err != nil {
						return err
					}
					log.Print("[DEBUG] Deleted a contact group")
				}
			}
			return nil
		},
	})
}

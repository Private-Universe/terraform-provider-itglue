package itglue

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	itglueRest "github.com/Private-Universe/itglue"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccITGlueOrderBasic(t *testing.T) {
	coffeeID := "1"
	quantity := "2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckITGlueOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckITGlueFlexibleAssetConfigBasic(coffeeID, quantity),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckITGlueFlexibleAssetExists("itglue_flexible_asset.new"),
				),
			},
		},
	})
}

func testAccCheckITGlueOrderDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*itglueRest.ITGAPI)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "itglue_flexible_asset" {
			continue
		}

		assetID := rs.Primary.ID
		s := strings.Split(assetID, "-")
		id, err := strconv.Atoi(s[1])
		_, err = c.DeleteFlexibleAsset(id)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckITGlueFlexibleAssetConfigBasic(coffeeID, quantity string) string {
	return fmt.Sprintf(`
	resource "itglue_flexible_asset" "new" {
		items {
			coffee {
				id = %s
			}
    		quantity = %s
  		}
	}
	`, coffeeID, quantity)
}

func testAccCheckITGlueFlexibleAssetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OrderID set")
		}

		return nil
	}
}

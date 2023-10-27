// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timestreamwrite_test

import (
	// TIP: ==== IMPORTS ====
	// This is a common set of imports but not customized to your code since
	// your code hasn't been written yet. Make sure you, your IDE, or
	// goimports -w <file> fixes these imports.
	//
	// The provider linter wants your imports to be in two groups: first,
	// standard library (i.e., "fmt" or "strings"), second, everything else.
	//
	// Also, AWS Go SDK v2 may handle nested structures differently than v1,
	// using the services/timestreamwrite/types package. If so, you'll
	// need to import types and reference the nested types, e.g., as
	// types.<Type Name>.
	"fmt"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"

	"github.com/hashicorp/terraform-provider-aws/names"
)

// TIP: File Structure. The basic outline for all test files should be as
// follows. Improve this data source's maintainability by following this
// outline.
//
// 1. Package declaration (add "_test" since this is a test file)
// 2. Imports
// 3. Unit tests
// 4. Basic test
// 5. Disappears test
// 6. All the other tests
// 7. Helper functions (exists, destroy, check, etc.)
// 8. Functions that return Terraform configurations

// TIP: ==== ACCEPTANCE TESTS ====
// This is an example of a basic acceptance test. This should test as much of
// standard functionality of the data source as possible, and test importing, if
// applicable. We prefix its name with "TestAcc", the service, and the
// data source name.
//
// Acceptance test access AWS and cost money to run.
func TestAccTimestreamWriteDatabaseDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	resourceName := "aws_timestreamwrite_database.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_timestreamwrite_database.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.TimestreamWriteEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDatabaseDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "creation_time", resourceName, "creation_time"),
					resource.TestCheckResourceAttrPair(dataSourceName, "database_name", resourceName, "database_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "kms_key_id", resourceName, "kms_key_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "last_updated_time", resourceName, "last_updated_time"),
					resource.TestCheckResourceAttrPair(dataSourceName, "table_count", resourceName, "table_count"),
				),
			},
		},
	})
}

func testAccDatabaseDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description = %[1]q

  policy = <<POLICY
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Effect": "Allow",
     "Principal": {
       "AWS": "*"
     },
     "Action": "kms:*",
     "Resource": "*"
   }
 ]
}
POLICY
}

resource "aws_timestreamwrite_database" "test" {
  database_name = %[1]q
  kms_key_id    = aws_kms_key.test.arn
}

data "aws_timestreamwrite_database" "test" {
  database_name = aws_timestreamwrite_database.name
}`, rName)
}

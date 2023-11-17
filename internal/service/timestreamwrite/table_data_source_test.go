// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timestreamwrite_test

// **PLEASE DELETE THIS AND ALL TIP COMMENTS BEFORE SUBMITTING A PR FOR REVIEW!**
//
// TIP: ==== INTRODUCTION ====
// Thank you for trying the skaff tool!
//
// You have opted to include these helpful comments. They all include "TIP:"
// to help you find and remove them when you're done with them.
//
// While some aspects of this file are customized to your input, the
// scaffold tool does *not* look at the AWS API and ensure it has correct
// function, structure, and variable names. It makes guesses based on
// commonalities. You will need to make significant adjustments.
//
// In other words, as generated, this is a rough outline of the work you will
// need to do. If something doesn't make sense for your situation, get rid of
// it.

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
	"github.com/YakDriver/regexache"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite"
	"github.com/aws/aws-sdk-go-v2/service/timestreamwrite/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
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

// TIP: ==== UNIT TESTS ====
// This is an example of a unit test. Its name is not prefixed with
// "TestAcc" like an acceptance test.
//
// Unlike acceptance tests, unit tests do not access AWS and are focused on a
// function (or method). Because of this, they are quick and cheap to run.
//
// In designing a data source's implementation, isolate complex bits from AWS bits
// so that they can be tested through a unit test. We encourage more unit tests
// in the provider.
//
// Cut and dry functions using well-used patterns, like typical flatteners and
// expanders, don't need unit testing. However, if they are complex or
// intricate, they should be unit tested.
//func TestTableExampleUnitTest(t *testing.T) {
//	t.Parallel()
//
//	testCases := []struct {
//		TestName string
//		Input    string
//		Expected string
//		Error    bool
//	}{
//		{
//			TestName: "empty",
//			Input:    "",
//			Expected: "",
//			Error:    true,
//		},
//		{
//			TestName: "descriptive name",
//			Input:    "some input",
//			Expected: "some output",
//			Error:    false,
//		},
//		{
//			TestName: "another descriptive name",
//			Input:    "more input",
//			Expected: "more output",
//			Error:    false,
//		},
//	}
//
//	for _, testCase := range testCases {
//		testCase := testCase
//		t.Run(testCase.TestName, func(t *testing.T) {
//			t.Parallel()
//			got, err := tftimestreamwrite.FunctionFromDataSource(testCase.Input)
//
//			if err != nil && !testCase.Error {
//				t.Errorf("got error (%s), expected no error", err)
//			}
//
//			if err == nil && testCase.Error {
//				t.Errorf("got (%s) and no error, expected error", got)
//			}
//
//			if got != testCase.Expected {
//				t.Errorf("got %s, expected %s", got, testCase.Expected)
//			}
//		})
//	}
//}

// TIP: ==== ACCEPTANCE TESTS ====
// This is an example of a basic acceptance test. This should test as much of
// standard functionality of the data source as possible, and test importing, if
// applicable. We prefix its name with "TestAcc", the service, and the
// data source name.
//
// Acceptance test access AWS and cost money to run.
func TestAccTimestreamWriteTableDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var table types.Table

	databaseName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	tableName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	dataSourceName := "data.aws_timestreamwrite_table.test"
	resourceName := "aws_timestreamwrite_table.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, timestreamwrite.ServiceID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.TimestreamWriteEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTableDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTableDataSourceConfig_basic(databaseName, tableName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableExists(ctx, dataSourceName, &table),
					acctest.MatchResourceAttrRegionalARN(dataSourceName, "arn", "timestreamwrite", regexache.MustCompile(`table:+.`)),
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "creation_time", resourceName, "creation_time"),
					resource.TestCheckResourceAttrPair(dataSourceName, "database_name", resourceName, "database_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "last_updated_time", resourceName, "last_updated_time"),
					resource.TestCheckResourceAttrPair(dataSourceName, "table_name", resourceName, "table_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "table_status", resourceName, "table_status"),
				),
			},
		},
	})
}

func testAccTableDataSourceConfig_basic(databaseName string, tableName string) string {
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

resource "aws_timestreamwrite_table" "test" {
  database_name = aws_timestreamwrite_database.test.name
  table_name    = %[2]q
}

data "aws_timestreamwrite_table" "test" {
  database_name = aws_timestreamwrite_database.test.name
  table_name    = aws_timestreamwrite_table.test.name
}
`, databaseName, tableName)
}

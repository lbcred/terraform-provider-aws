---
subcategory: "Timestream Write"
layout: "aws"
page_title: "AWS: aws_timestreamwrite_database"
description: |-
  Terraform data source that provides details for an AWS Timestream Write Database.
---
<!---
TIP: A few guiding principles for writing documentation:
1. Use simple language while avoiding jargon and figures of speech.
2. Focus on brevity and clarity to keep a reader's attention.
3. Use active voice and present tense whenever you can.
4. Document your feature as it exists now; do not mention the future or past if you can help it.
5. Use accessible and inclusive language.
--->

# Data Source: aws_timestreamwrite_database

Terraform data source that provides details for an AWS Timestream Write Database.

## Example Usage

### Basic Usage

```terraform
data "aws_timestreamwrite_database" "example" {
  database_name = "example-database"
}
```

## Argument Reference

The following arguments are required:

* `database_name` - (Required) Concise argument description. Do not begin the description with "An", "The", "Defines", "Indicates", or "Specifies," as these are verbose. In other words, "Indicates the amount of storage," can be rewritten as "Amount of storage," without losing any information.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `arn` - ARN of the Database. Do not begin the description with "An", "The", "Defines", "Indicates", or "Specifies," as these are verbose. In other words, "Indicates the amount of storage," can be rewritten as "Amount of storage," without losing any information.
* `creation_time` - 
* `database_name` - 
* `kms_key_id` - 
* `last_updated_time` - 
* `table_count` - 

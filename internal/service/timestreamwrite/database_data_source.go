// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timestreamwrite

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
	// awstypes.<Type Name>.
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// TIP: ==== FILE STRUCTURE ====
// All data sources should follow this basic outline. Improve this data source's
// maintainability by sticking to it.
//
// 1. Package declaration
// 2. Imports
// 3. Main data source struct with schema method
// 4. Read method
// 5. Other functions (flatteners, expanders, waiters, finders, etc.)

// Function annotations are used for datasource registration to the Provider. DO NOT EDIT.
// @FrameworkDataSource("aws_timestreamwrite_database")
func newDataSourceDatabase(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceDatabase{}, nil
}

const (
	DSNameDatabase = "Database Data Source"
)

type dataSourceDatabase struct {
	framework.DataSourceWithConfigure
}

func (d *dataSourceDatabase) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) { // nosemgrep:ci.meta-in-func-name
	resp.TypeName = "aws_timestreamwrite_database"
}

func (d *dataSourceDatabase) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"arn": framework.ARNAttributeComputedOnly(),
			"creation_time": schema.Int64Attribute{
				Computed: true,
			},
			"database_name": schema.StringAttribute{
				Required: true,
			},
			"kms_key_id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated_time": schema.Int64Attribute{
				Computed: true,
			},
			"table_count": schema.Int64Attribute{
				Computed: true,
			},
			"tags": tftags.TagsAttributeComputedOnly(),
		},
	}
}

func (d *dataSourceDatabase) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// 1. Get a client connection to the relevant service
	// 2. Fetch the config
	// 3. Get information about a resource from AWS
	// 4. Set the ID, arguments, and attributes
	// 5. Set the tags
	// 6. Set the state
	conn := d.Meta().TimestreamWriteClient(ctx)

	var data dataSourceDatabaseData
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := findDatabaseByName(ctx, conn, data.DatabaseName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.TimestreamWrite, create.ErrActionReading, DSNameDatabase, data.DatabaseName.String(), err),
			err.Error(),
		)
		return
	}

	creationTime := out.CreationTime.Unix()
	lastUpdatedTime := out.LastUpdatedTime.Unix()

	data.ARN = flex.StringToFramework(ctx, out.Arn)
	data.CreationTime = flex.Int64ToFramework(ctx, &creationTime)
	data.DatabaseName = flex.StringToFramework(ctx, out.DatabaseName)
	data.KmsKeyId = flex.StringToFramework(ctx, out.KmsKeyId)
	data.LastUpdatedTime = flex.Int64ToFramework(ctx, &lastUpdatedTime)
	data.TableCount = flex.Int64ToFramework(ctx, &out.TableCount)

	// TIP: -- 5. Set the tags

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type dataSourceDatabaseData struct {
	ARN             types.String `tfsdk:"arn"`
	CreationTime    types.Int64  `tfsdk:"creation_time"`
	DatabaseName    types.String `tfsdk:"database_name"`
	KmsKeyId        types.String `tfsdk:"kms_key_id"`
	LastUpdatedTime types.Int64  `tfsdk:"last_updated_time"`
	TableCount      types.Int64  `tfsdk:"table_count"`
}

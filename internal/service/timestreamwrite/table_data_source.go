// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timestreamwrite

import (
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
// @FrameworkDataSource(name="Table")
func newDataSourceTable(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceTable{}, nil
}

const (
	DSNameTable = "Table Data Source"
)

type dataSourceTable struct {
	framework.DataSourceWithConfigure
}

func (d *dataSourceTable) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) { // nosemgrep:ci.meta-in-func-name
	resp.TypeName = "aws_timestreamwrite_table"
}

func (d *dataSourceTable) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"arn": framework.ARNAttributeComputedOnly(),
			"creation_time": schema.Int64Attribute{
				Computed: true,
			},
			"database_name": schema.StringAttribute{
				Required: true,
			},
			"last_updated_time": schema.Int64Attribute{
				Computed: true,
			},
			"table_name": schema.StringAttribute{
				Required: true,
			},
			"table_status": schema.StringAttribute{
				Computed: true,
			},
			"tags": tftags.TagsAttributeComputedOnly(),
		},
		Blocks: map[string]schema.Block{
			"magnetic_store_write_properties": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable_magnetic_store_writes": schema.BoolAttribute{
							Computed: true,
						},
					},
					Blocks: map[string]schema.Block{
						"magnetic_store_rejected_data_location": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Blocks: map[string]schema.Block{
									"s3_configuration": schema.SetNestedBlock{
										NestedObject: schema.NestedBlockObject{
											Attributes: map[string]schema.Attribute{
												"bucket_name": schema.StringAttribute{
													Computed: true,
												},
												"encryption_option": schema.StringAttribute{
													Computed: true,
												},
												"kms_key_id": schema.StringAttribute{
													Computed: true,
												},
												"object_key_prefix": schema.StringAttribute{
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"retention_properties": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"magnetic_store_retention_period_in_days": schema.Int64Attribute{
							Computed: true,
						},
						"memory_store_retention_period_in_hours": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
			"schema": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"composite_partition_key": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									// TIP: Attributes that are required on a corresponding resource will be
									// computed on the data source (unless required as part of the search criteria).
									"enforcement_in_record": schema.StringAttribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
									"type": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *dataSourceTable) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	conn := d.Meta().TimestreamWriteClient(ctx)

	var data dataSourceTableData
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tableName := data.TableName
	databaseName := data.DatabaseName

	out, err := findTableByTwoPartKey(ctx, conn, databaseName.String(), tableName.String())
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.TimestreamWrite, create.ErrActionReading, DSNameTable, databaseName.String(), err),
			err.Error(),
		)
		return
	}

	creationTime := out.CreationTime.Unix()
	lastUpdatedTime := out.LastUpdatedTime.Unix()
	tableStatus := string(out.TableStatus)

	data.ARN = flex.StringToFramework(ctx, out.Arn)
	data.TableName = flex.StringToFramework(ctx, out.TableName)
	data.DatabaseName = flex.StringToFramework(ctx, out.DatabaseName)
	data.LastUpdatedTime = flex.Int64ToFramework(ctx, &lastUpdatedTime)
	data.CreationTime = flex.Int64ToFramework(ctx, &creationTime)
	data.TableStatus = flex.StringToFramework(ctx, &tableStatus)

	// TIP: -- 4. Set the ID, arguments, and attributes
	//
	// For simple data types (i.e., schema.StringAttribute, schema.BoolAttribute,
	// schema.Int64Attribute, and schema.Float64Attribue), simply setting the
	// appropriate data struct field is sufficient. The flex package implements
	// helpers for converting between Go and Plugin-Framework types seamlessly. No
	// error or nil checking is necessary.
	//
	// However, there are some situations where more handling is needed such as
	// complex data types (e.g., schema.ListAttribute, schema.SetAttribute). In
	// these cases the flatten function may have a diagnostics return value, which
	// should be appended to resp.Diagnostics.

	FlattenedMagneticStoreWriteProperties := flattenMagneticStoreWriteProperties(out.MagneticStoreWriteProperties)
	MagneticStoreWriteProperties := FlattenedMagneticStoreWriteProperties[0].(magneticStoreWriteProperties)

	data.MagneticStoreWriteProperties = MagneticStoreWriteProperties

	FlattenedRetentionProperties := flattenRetentionProperties(out.RetentionProperties)
	RetentionProperties := FlattenedRetentionProperties[0].(retentionProperties)

	data.RetentionProperties = RetentionProperties

	//
	//// TIP: Setting a complex type.
	//complexArgument, diag := flattenComplexArgument(ctx, out.ComplexArgument)
	//resp.Diagnostics.Append(diag...)
	//data.ComplexArgument = complexArgument

	// TIP: -- 5. Set the tags

	// TIP: -- 6. Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// TIP: ==== DATA STRUCTURES ====
// With Terraform Plugin-Framework configurations are deserialized into
// Go types, providing type safety without the need for type assertions.
// These structs should match the schema definition exactly, and the `tfsdk`
// tag value should match the attribute name.
//
// Nested objects are represented in their own data struct. These will
// also have a corresponding attribute type mapping for use inside flex
// functions.
//
// See more:
// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/accessing-values
type dataSourceTableData struct {
	ARN                          types.String                 `tfsdk:"arn"`
	CreationTime                 types.Int64                  `tfsdk:"creation_time"`
	DatabaseName                 types.String                 `tfsdk:"database_name"`
	LastUpdatedTime              types.Int64                  `tfsdk:"last_updated_time"`
	MagneticStoreWriteProperties magneticStoreWriteProperties `tfsdk:"magnetic_store_write_properties"`
	RetentionProperties          retentionProperties          `tfsdk:"retention_properties"`
	TableSchema                  tableSchema                  `tfsdk:"schema"`
	TableName                    types.String                 `tfsdk:"table_name"`
	TableStatus                  types.String                 `tfsdk:"table_status"`
}

type magneticStoreWriteProperties struct {
	EnableMagneticStoreWrites         types.Bool                        `tfsdk:"enable_magnetic_store_writes"`
	MagneticStoreRejectedDataLocation magneticStoreRejectedDataLocation `tfsdk:"magnetic_store_rejected_data_location"`
}

type magneticStoreRejectedDataLocation struct {
	S3Configuration s3Configuration `tfsdk:"s3_configuration"`
}

type s3Configuration struct {
	BucketName       types.String `tfsdk:"bucket_name"`
	EncryptionOption types.String `tfsdk:"encryption_option"`
	KmsKeyId         types.String `tfsdk:"kms_key_id"`
	ObjectKeyPrefix  types.String `tfsdk:"object_key_prefix"`
}

type retentionProperties struct {
	MagneticStoreRetentionPeriodInDays types.Int64 `tfsdk:"magnetic_store_retention_period_in_days"`
	MemoryStoreRetentionPeriodInHours  types.Int64 `tfsdk:"memory_store_retention_period_in_hours"`
}

type tableSchema struct {
	CompositePartitionKey types.List `tfsdk:"composite_partition_key"`
}

type compositePartitionKey struct {
	EnforcementInRecord types.String `tfsdk:"enforcement_in_record"`
	Name                types.String `tfsdk:"name"`
	Type                types.String `tfsdk:"type"`
}

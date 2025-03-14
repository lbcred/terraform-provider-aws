// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package flex

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
)

// StringFromFramework converts a Framework String value to a string pointer.
// A null String is converted to a nil string pointer.
func StringFromFramework(ctx context.Context, v basetypes.StringValuable) *string {
	var output *string

	panicOnError(Expand(ctx, v, &output))

	return output
}

// StringFromFramework converts a single Framework String value to a string pointer slice.
// A null String is converted to a nil slice.
func StringSliceFromFramework(ctx context.Context, v basetypes.StringValuable) []*string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	return []*string{StringFromFramework(ctx, v)}
}

// StringValueToFramework converts a string value to a Framework String value.
// An empty string is converted to a null String.
func StringValueToFramework[T ~string](ctx context.Context, v T) types.String {
	if v == "" {
		return types.StringNull()
	}

	var output types.String

	panicOnError(Flatten(ctx, v, &output))

	return output
}

// StringValueToFrameworkLegacy converts a string value to a Framework String value.
// An empty string is left as an empty String.
func StringValueToFrameworkLegacy[T ~string](_ context.Context, v T) types.String {
	return types.StringValue(string(v))
}

// StringToFramework converts a string pointer to a Framework String value.
// A nil string pointer is converted to a null String.
func StringToFramework(ctx context.Context, v *string) types.String {
	var output types.String

	panicOnError(Flatten(ctx, v, &output))

	return output
}

// StringToFrameworkLegacy converts a string pointer to a Framework String value.
// A nil string pointer is converted to an empty String.
func StringToFrameworkLegacy(_ context.Context, v *string) types.String {
	return types.StringValue(aws.ToString(v))
}

// StringToFrameworkARN converts a string pointer to a Framework custom ARN value.
// A nil string pointer is converted to a null ARN.
// If diags is nil, any errors cause a panic.
func StringToFrameworkARN(ctx context.Context, v *string, diags *diag.Diagnostics) fwtypes.ARN {
	var output fwtypes.ARN

	if diags == nil {
		panicOnError(Flatten(ctx, v, &output))
	} else {
		diags.Append(Flatten(ctx, v, &output)...)
	}

	return output
}

func StringFromFrameworkLegacy(_ context.Context, v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	s := v.ValueString()
	if s == "" {
		return nil
	}

	return aws.String(s)
}

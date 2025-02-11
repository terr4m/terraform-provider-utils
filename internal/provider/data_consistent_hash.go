package provider

import (
	"context"
	"fmt"

	"github.com/terr4m/terraform-provider-utils/internal/hash"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ConsistentHashDataSource{}
var _ datasource.DataSourceWithValidateConfig = &ConsistentHashDataSource{}

// NewConsistentHashDataSource creates a new consistent hash data source.
func NewConsistentHashDataSource() datasource.DataSource {
	return &ConsistentHashDataSource{}
}

// ConsistentHashDataSource defines the data source implementation.
type ConsistentHashDataSource struct {
}

// ConsistentHashDataSourceModel describes the data source data model.
type ConsistentHashDataSourceModel struct {
	Members           types.Set     `tfsdk:"members"`
	Keys              types.Set     `tfsdk:"keys"`
	PartitionCount    types.Int64   `tfsdk:"partition_count"`
	ReplicationFactor types.Int64   `tfsdk:"replication_factor"`
	Load              types.Float64 `tfsdk:"load"`
	Mapping           types.Map     `tfsdk:"mapping"`
}

// Metadata returns the data source metadata.
func (d *ConsistentHashDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_consistent_hash", req.ProviderTypeName)
}

// Schema returns the data source schema.
func (d *ConsistentHashDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Utils consistent hash TF data source.",
		Attributes: map[string]schema.Attribute{
			"members": schema.SetAttribute{
				MarkdownDescription: "The members of the consistent hash.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"keys": schema.SetAttribute{
				MarkdownDescription: "The keys of the consistent hash.",
				ElementType:         types.StringType,
				Required:            true,
			},
			"partition_count": schema.Int64Attribute{
				MarkdownDescription: "The number of partitions to use for hashing.",
				Optional:            true,
				Computed:            true,
			},
			"replication_factor": schema.Int64Attribute{
				MarkdownDescription: "The number of replicas to use for hashing.",
				Optional:            true,
				Computed:            true,
			},
			"load": schema.Float64Attribute{
				MarkdownDescription: "The load factor to use for hashing.",
				Optional:            true,
				Computed:            true,
			},
			"mapping": schema.MapAttribute{
				MarkdownDescription: "The mapping of keys to members.",
				ElementType:         types.SetType{ElemType: types.StringType},
				Computed:            true,
			},
		},
	}
}

// ValidateConfig validates the data source configuration.
func (d *ConsistentHashDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data ConsistentHashDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	if !data.PartitionCount.IsNull() && data.PartitionCount.ValueInt64() < int64(len(data.Members.Elements())) {
		resp.Diagnostics.AddAttributeError(path.Root("partition_count"), "Partition count invalid.", "partition count must be greater than or equal to the number of members")
	}
}

// Read reads the data source.
func (d *ConsistentHashDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConsistentHashDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	var members []string
	if resp.Diagnostics.Append(data.Members.ElementsAs(ctx, &members, false)...); resp.Diagnostics.HasError() {
		return
	}

	var keys []string
	if resp.Diagnostics.Append(data.Keys.ElementsAs(ctx, &keys, false)...); resp.Diagnostics.HasError() {
		return
	}

	var partitionCount int
	if !data.PartitionCount.IsNull() {
		partitionCount = int(data.PartitionCount.ValueInt64())
	} else {
		partitionCount = hash.DefaultPartitionCount
	}

	var replicationFactor int
	if !data.ReplicationFactor.IsNull() {
		replicationFactor = int(data.ReplicationFactor.ValueInt64())
	} else {
		replicationFactor = hash.DefaultReplicationFactor
	}

	var load float64
	if !data.Load.IsNull() {
		load = data.Load.ValueFloat64()
	} else {
		load = hash.DefaultLoad
	}

	ch := hash.NewConsistentHash(members, partitionCount, replicationFactor, load)
	mapping, err := ch.CalculateMapping(keys)
	if err != nil {
		resp.Diagnostics.AddError("Invalid consistent hash configuration", err.Error())
		return
	}

	tfMapping, diags := types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, mapping)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	data.PartitionCount = types.Int64Value(int64(partitionCount))
	data.ReplicationFactor = types.Int64Value(int64(replicationFactor))
	data.Load = types.Float64Value(load)
	data.Mapping = tfMapping

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

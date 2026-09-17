package provider

import (
	"context"
	"fmt"

	"github.com/terr4m/terraform-provider-utils/internal/hash"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ConsistentHashDataSource{}

// NewConsistentHashDataSource creates a new consistent hash data source.
func NewConsistentHashDataSource() datasource.DataSource {
	return &ConsistentHashDataSource{}
}

// ConsistentHashDataSource defines the data source implementation.
type ConsistentHashDataSource struct{}

// ConsistentHashDataSourceModel describes the data source data model.
type ConsistentHashDataSourceModel struct {
	Members types.Set     `tfsdk:"members"`
	Keys    types.Set     `tfsdk:"keys"`
	Mapping *MappingModel `tfsdk:"mapping"`
}

type MappingModel struct {
	Keys    types.Map `tfsdk:"keys"`
	Members types.Map `tfsdk:"members"`
}

// Metadata returns the data source metadata.
func (d *ConsistentHashDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_consistent_hash", req.ProviderTypeName)
}

// Schema returns the data source schema.
func (d *ConsistentHashDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source to configure a hash ring that provides a consistent hashing function which simultaneously achieves both uniformity and consistency.",
		Attributes: map[string]schema.Attribute{
			"members": schema.SetAttribute{
				MarkdownDescription: "Members to configure the hash ring for.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"keys": schema.SetAttribute{
				MarkdownDescription: "Keys to hash with the hash ring.",
				ElementType:         types.StringType,
				Required:            true,
			},
			"mapping": schema.SingleNestedAttribute{
				MarkdownDescription: "Mapping between the keys and the members of the hash ring.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"keys": schema.MapAttribute{
						MarkdownDescription: "Mapping of keys to members.",
						ElementType:         types.StringType,
						Computed:            true,
					},
					"members": schema.MapAttribute{
						MarkdownDescription: "Mapping of members to keys.",
						ElementType:         types.SetType{ElemType: types.StringType},
						Computed:            true,
					},
				},
			},
		},
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

	ring, err := hash.NewRing(hash.WithMembers(members))
	if err != nil {
		resp.Diagnostics.AddError("Could not configure hash ring.", err.Error())
		return
	}

	keyMap, err := ring.Calculate(keys)
	if err != nil {
		resp.Diagnostics.AddError("Unable to calculate hash mapping.", err.Error())
		return
	}

	memberMap := make(map[string][]string, len(members))
	for _, member := range members {
		memberMap[member] = []string{}
	}

	for k, v := range keyMap {
		memberMap[v] = append(memberMap[v], k)
	}

	km, diags := types.MapValueFrom(ctx, types.StringType, keyMap)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	mm, diags := types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, memberMap)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	data.Mapping = &MappingModel{
		Keys:    km,
		Members: mm,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure UtilsProvider satisfies various provider interfaces.
var _ provider.Provider = &UtilsProvider{}
var _ provider.ProviderWithFunctions = &UtilsProvider{}
var _ provider.ProviderWithEphemeralResources = &UtilsProvider{}

// New returns a new provider implementation.
func New(version, commit string) func() provider.Provider {
	return func() provider.Provider {
		return &UtilsProvider{
			version: version,
			commit:  commit,
		}
	}
}

// UtilsProviderData is the data available to the resource and data sources.
type UtilsProviderData struct {
	provider *UtilsProvider
	Model    *UtilsProviderModel
}

// UtilsProviderModel describes the provider data model.
type UtilsProviderModel struct {
}

// UtilsProvider defines the provider implementation.
type UtilsProvider struct {
	version string
	commit  string
}

func (p *UtilsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "utils"
	resp.Version = p.version
}

func (p *UtilsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Utils provider.",
	}
}

// Configure configures the provider.
func (p *UtilsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if req.ClientCapabilities.DeferralAllowed && !req.Config.Raw.IsFullyKnown() {
		resp.Deferred = &provider.Deferred{
			Reason: provider.DeferredReasonProviderConfigUnknown,
		}
	}

	// Load the provider config
	model := &UtilsProviderModel{}
	diags := req.Config.Get(ctx, model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configure provider data
	providerData := &UtilsProviderData{
		provider: p,
		Model:    model,
	}

	resp.DataSourceData = providerData
	resp.EphemeralResourceData = providerData
	resp.ResourceData = providerData
}

func (p *UtilsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewConsistentHashDataSource,
	}
}

func (p *UtilsProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *UtilsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *UtilsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

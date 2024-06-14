package provider

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Khan/genqlient/graphql"
)

var (
	envVarName         = "CALA_API_ENDPOINT"
	errMissingEndpoint = "Required endpoint could not be found. Please set the endpoint using an input variable in the provider configuration block or by using the `" + envVarName + "` environment variable."
)

var _ provider.Provider = &CalaProvider{}

type CalaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type CalaProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *CalaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cala"
	resp.Version = p.version
}

func (p *CalaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The endpoint for cala server.",
				Required:            true,
			},
		},
	}
}

func (p *CalaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CalaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := ""

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	// If a endpoint wasn't set in the provider configuration block, try and fetch it
	// from the environment variable.
	if endpoint == "" {
		endpoint = os.Getenv(envVarName)
	}

	// If we still don't have a endpoint at this point, we return an error.
	if endpoint == "" {
		resp.Diagnostics.AddError("Missing Endpoint", errMissingEndpoint)
		return
	}

	httpClient := http.Client{
		Transport: &authedTransport{
			endpoint: endpoint,
			wrapped:  http.DefaultTransport,
		},
	}

	client := graphql.NewClient(endpoint, &httpClient)

	resp.DataSourceData = &client
	resp.ResourceData = &client
}

func (p *CalaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAccountResource,
	}
}

func (p *CalaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CalaProvider{
			version: version,
		}
	}
}

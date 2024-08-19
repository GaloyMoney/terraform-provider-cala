package provider

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &BigQueryIntegrationResource{}
var _ resource.ResourceWithImportState = &BigQueryIntegrationResource{}

func NewBigQueryIntegrationResource() resource.Resource {
	return &BigQueryIntegrationResource{}
}

type BigQueryIntegrationResource struct {
	client *graphql.Client
}

type BigQueryIntegrationResourceModel struct {
	BigQueryIntegrationId     types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	ServiceAccountCredsBase64 types.String `tfsdk:"service_account_creds_base64"`
	ProjectId                 types.String `tfsdk:"project_id"`
	DatasetId                 types.String `tfsdk:"dataset_id"`
}

func (r *BigQueryIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_big_query_integration"
}

func (r *BigQueryIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cala BigQuery Integration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the integration.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the integration.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the integration.",
				Optional:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Gcp Project Id",
				Required:            true,
			},
			"dataset_id": schema.StringAttribute{
				MarkdownDescription: "Gcp Biq Query Dataset Id",
				Required:            true,
			},
			"service_account_creds_base64": schema.StringAttribute{
				MarkdownDescription: "The GCP service account creds",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (r *BigQueryIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*graphql.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *BigQueryIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *BigQueryIntegrationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := BigQueryIntegrationCreateInput{
		IntegrationId:             data.BigQueryIntegrationId.ValueString(),
		Name:                      data.Name.ValueString(),
		Description:               data.Description.ValueStringPointer(),
		GcpProjectId:              data.ProjectId.ValueString(),
		GcpDatasetId:              data.DatasetId.ValueString(),
		ServiceAccountCredsBase64: data.ServiceAccountCredsBase64.ValueString(),
	}

	response, err := bigQueryIntegrationCreate(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create integration, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an integration")

	integration := response.BigQuery.IntegrationCreate.Integration

	data.BigQueryIntegrationId = types.StringValue(integration.IntegrationId)
	data.Name = types.StringValue(integration.Name)
	data.Description = types.StringPointerValue(integration.Description)
	data.ProjectId = types.StringValue(integration.GcpProjectId)
	data.DatasetId = types.StringValue(integration.GcpDatasetId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BigQueryIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *BigQueryIntegrationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	response, err := bigQueryIntegrationGet(ctx, *r.client, data.BigQueryIntegrationId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read integration, got error: %s", err))
		return
	}

	integration := response.BigQuery.Integration

	data.BigQueryIntegrationId = types.StringValue(integration.IntegrationId)
	data.Name = types.StringValue(integration.Name)
	data.Description = types.StringPointerValue(integration.Description)
	data.ProjectId = types.StringValue(integration.GcpProjectId)
	data.DatasetId = types.StringValue(integration.GcpDatasetId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BigQueryIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *BigQueryIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *BigQueryIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}

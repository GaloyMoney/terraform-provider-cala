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

var _ resource.Resource = &BitfinexIntegrationResource{}
var _ resource.ResourceWithImportState = &BitfinexIntegrationResource{}

func NewBitfinexIntegrationResource() resource.Resource {
	return &BitfinexIntegrationResource{}
}

type BitfinexIntegrationResource struct {
	client *graphql.Client
}

type BitfinexIntegrationResourceModel struct {
	BitfinexIntegrationId types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	JournalId             types.String `tfsdk:"journal_id"`
	Key                   types.String `tfsdk:"key"`
	Secret                types.String `tfsdk:"secret"`
	OmnibusAccountSetId   types.String `tfsdk:"omnibus_account_set_id"`
}

func (r *BitfinexIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bitfinex_integration"
}

func (r *BitfinexIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cala Bitfinex Integration.",
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
			"journal_id": schema.StringAttribute{
				MarkdownDescription: "journal_id",
				Required:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The bitfinex API key",
				Required:            true,
				Sensitive:           true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "The bitfinex API secret",
				Required:            true,
				Sensitive:           true,
			},
			"omnibus_account_set_id": schema.StringAttribute{
				MarkdownDescription: "The AccountSet id for the omnibus AccountSet",
				Computed:            true,
			},
		},
	}
}

func (r *BitfinexIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BitfinexIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *BitfinexIntegrationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := BfxIntegrationCreateInput{
		IntegrationId: data.BitfinexIntegrationId.ValueString(),
		Name:          data.Name.ValueString(),
		Description:   data.Description.ValueStringPointer(),
		JournalId:     data.JournalId.ValueString(),
		Key:           data.Key.ValueString(),
		Secret:        data.Secret.ValueString(),
	}

	response, err := bfxIntegrationCreate(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create integration, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an integration")

	integration := response.Bitfinex.IntegrationCreate.Integration

	data.BitfinexIntegrationId = types.StringValue(integration.IntegrationId)
	data.Name = types.StringValue(integration.Name)
	data.Description = types.StringPointerValue(integration.Description)
	data.OmnibusAccountSetId = types.StringValue(integration.OmnibusAccountSetId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BitfinexIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *BitfinexIntegrationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	response, err := bfxIntegrationGet(ctx, *r.client, data.BitfinexIntegrationId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read integration, got error: %s", err))
		return
	}

	integration := response.Bitfinex.Integration

	data.BitfinexIntegrationId = types.StringValue(integration.IntegrationId)
	data.Name = types.StringValue(integration.Name)
	data.Description = types.StringPointerValue(integration.Description)
	data.OmnibusAccountSetId = types.StringValue(integration.OmnibusAccountSetId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BitfinexIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *BitfinexIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *BitfinexIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}

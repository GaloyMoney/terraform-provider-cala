package provider

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &AccountSetResource{}
var _ resource.ResourceWithImportState = &AccountSetResource{}

func NewAccountSetResource() resource.Resource {
	return &AccountSetResource{}
}

type AccountSetResource struct {
	client *graphql.Client
}

type AccountSetResourceModel struct {
	AccountSetId      types.String `tfsdk:"id"`
	JournalId         types.String `tfsdk:"journal_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	NormalBalanceType types.String `tfsdk:"normal_balance_type"`
}

func (r *AccountSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_set"
}

func (r *AccountSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cala account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the account.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the account.",
				Required:            true,
			},
			"journal_id": schema.StringAttribute{
				MarkdownDescription: "ID of the journal.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the account.",
				Optional:            true,
			},
			"normal_balance_type": schema.StringAttribute{
				MarkdownDescription: "normalBalanceType",
				Default:             stringdefault.StaticString("CREDIT"),
				Computed:            true,
			},
		},
	}
}

func (r *AccountSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AccountSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AccountSetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	normalBalanceType, err := toDebitOrCredit(data.NormalBalanceType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Normal Balance Type", fmt.Sprintf("Unable to convert normal_balance_type to DebitOrCredit: %s", err))
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Invalid Status", fmt.Sprintf("Unable to convert status to Status: %s", err))
		return
	}

	input := AccountSetCreateInput{
		AccountSetId:      data.AccountSetId.ValueString(),
		JournalId:         data.JournalId.ValueString(),
		Name:              data.Name.ValueString(),
		Description:       data.Description.ValueStringPointer(),
		NormalBalanceType: normalBalanceType,
	}

	response, err := accountSetCreate(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create accountSet, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an accountSet")

	account := response.AccountSetCreate.AccountSet

	data.AccountSetId = types.StringValue(account.AccountSetId)
	data.JournalId = types.StringValue(account.JournalId)
	data.Name = types.StringValue(account.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AccountSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *AccountSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *AccountSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *AccountSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}

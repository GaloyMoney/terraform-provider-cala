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

var _ resource.Resource = &AccountResource{}
var _ resource.ResourceWithImportState = &AccountResource{}

func NewAccountResource() resource.Resource {
	return &AccountResource{}
}

type AccountResource struct {
	client *graphql.Client
}

type AccountResourceModel struct {
	AccountId       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Code     types.String `tfsdk:"code"`
	NormalBalanceType types.String `tfsdk:"normal_balance_type"`
	Status   types.String `tfsdk:"status"`
}

func (r *AccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (r *AccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"code": schema.StringAttribute{
				MarkdownDescription: "code",
				Required:            true,
			},
			"normal_balance_type": schema.StringAttribute{
				MarkdownDescription: "normalBalanceType",
				Default:             stringdefault.StaticString("CREDIT"),
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "status",
				Default:             stringdefault.StaticString("ACTIVE"),
				Computed:            true,
			},
		},
	}
}

func (r *AccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AccountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := AccountCreateInput{
		AccountId: data.AccountId.ValueString(),
		Name:      data.Name.ValueString(),
		Code: data.Code.ValueString(),
		NormalBalanceType: DebitOrCreditCredit,
		Status: StatusActive,
	}

	response, err := accountCreate(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create account, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an account")

	account := response.AccountCreate.Account

	data.AccountId = types.StringValue(account.AccountId)
	data.Name = types.StringValue(account.Name)
	data.Code = types.StringValue(account.Code)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *AccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *AccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *AccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}

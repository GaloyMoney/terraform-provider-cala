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
	AccountId         types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Code              types.String `tfsdk:"code"`
	NormalBalanceType types.String `tfsdk:"normal_balance_type"`
	Status            types.String `tfsdk:"status"`
	ExternalId        types.String `tfsdk:"external_id"`
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
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the account.",
				Optional:            true,
			},
			"code": schema.StringAttribute{
				MarkdownDescription: "code",
				Required:            true,
			},
			"normal_balance_type": schema.StringAttribute{
				MarkdownDescription: "normalBalanceType",
				Optional:            true,
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "status",
				Default:             stringdefault.StaticString("ACTIVE"),
				Computed:            true,
			},
			"external_id": schema.StringAttribute{
				MarkdownDescription: "externalId",
				Optional:            true,
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

	normalBalanceType, err := toDebitOrCredit(data.NormalBalanceType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Normal Balance Type", fmt.Sprintf("Unable to convert normal_balance_type to DebitOrCredit: %s", err))
		return
	}

	status, err := toStatus(data.Status.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Status", fmt.Sprintf("Unable to convert status to Status: %s", err))
		return
	}

	input := AccountCreateInput{
		AccountId:         data.AccountId.ValueString(),
		Name:              data.Name.ValueString(),
		Description:       data.Description.ValueStringPointer(),
		Code:              data.Code.ValueString(),
		NormalBalanceType: normalBalanceType,
		Status:            status,
		ExternalId:        data.ExternalId.ValueStringPointer(),
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
	data.Description = types.StringPointerValue(account.Description)
	data.ExternalId = types.StringPointerValue(account.ExternalId)
	data.NormalBalanceType = types.StringValue(string(account.NormalBalanceType))
	data.Status = types.StringValue(string(account.Status))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AccountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	response, err := accountGet(ctx, *r.client, data.AccountId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read account, got error: %s", err))
		return
	}

	account := response.Account

	data.AccountId = types.StringValue(account.AccountId)
	data.Description = types.StringPointerValue(account.Description)
	data.Name = types.StringValue(account.Name)
	data.Code = types.StringValue(account.Code)
	data.NormalBalanceType = types.StringValue(string(account.NormalBalanceType))
	data.Status = types.StringValue(string(account.Status))
	data.ExternalId = types.StringPointerValue(account.ExternalId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *AccountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	normalBalanceType, err := toDebitOrCredit(data.NormalBalanceType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Normal Balance Type", fmt.Sprintf("Unable to convert normal_balance_type to DebitOrCredit: %s", err))
		return
	}

	// Prepare the input for the update mutation, only updating the name field
	input := AccountUpdateInput{
		Name:              data.Name.ValueStringPointer(),
		Description:       data.Description.ValueStringPointer(),
		Code:              data.Code.ValueStringPointer(),
		NormalBalanceType: &normalBalanceType,
		ExternalId:        data.ExternalId.ValueStringPointer(),
	}

	// Call the update mutation
	response, err := accountUpdate(ctx, *r.client, data.AccountId.ValueString(), input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update account, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated an account")

	// Update the state with the new data
	account := response.AccountUpdate.Account

	data.AccountId = types.StringValue(account.AccountId)
	data.Name = types.StringValue(account.Name)
	data.Code = types.StringValue(account.Code)
	data.Description = types.StringPointerValue(account.Description)
	data.ExternalId = types.StringPointerValue(account.ExternalId)
	data.NormalBalanceType = types.StringValue(string(account.NormalBalanceType))
	data.Status = types.StringValue(string(account.Status))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *AccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}

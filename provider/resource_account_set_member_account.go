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

var _ resource.Resource = &AccountSetMemberAccountResource{}
var _ resource.ResourceWithImportState = &AccountSetMemberAccountResource{}

func NewAccountSetMemberAccountResource() resource.Resource {
	return &AccountSetMemberAccountResource{}
}

type AccountSetMemberAccountResource struct {
	client *graphql.Client
}

type AccountSetMemberAccountResourceModel struct {
	AccountSetMemberId types.String `tfsdk:"id"`
	AccountSetId       types.String `tfsdk:"account_set_id"`
	MemberAccountId    types.String `tfsdk:"member_account_id"`
}

func (r *AccountSetMemberAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_set_member_account"
}

func (r *AccountSetMemberAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Represents the membership of an account in an account set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the account.",
				Computed:            true,
			},
			"account_set_id": schema.StringAttribute{
				MarkdownDescription: "Id of the AccountSet",
				Required:            true,
			},
			"member_account_id": schema.StringAttribute{
				MarkdownDescription: "Id of the member AccountSet",
				Required:            true,
			},
		},
	}
}

func (r *AccountSetMemberAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AccountSetMemberAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AccountSetMemberAccountResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := accountSetMemberAccountCreate(ctx, *r.client, data.AccountSetId.ValueString(), data.MemberAccountId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create account set member, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Added an account to an account set")

	data.AccountSetMemberId = types.StringValue(fmt.Sprintf("account_set/%s/account/%s", data.AccountSetId.ValueString(), data.MemberAccountId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AccountSetMemberAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AccountSetMemberAccountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	response, err := accountGet(ctx, *r.client, data.MemberAccountId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read account, got error: %s", err))
		return
	}

	found := false
	for _, n := range response.GetAccount().Sets.Nodes {
		if n.AccountSetId == data.AccountSetId.ValueString() {
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("Not Found", "The account set ID does not match any nodes.")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AccountSetMemberAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *AccountSetMemberAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AccountSetMemberAccountResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	_, err := accountSetMemberAccountRemove(ctx, *r.client, data.AccountSetId.ValueString(), data.MemberAccountId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete account set member, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Removed an account from an account set")
}

func (r *AccountSetMemberAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}

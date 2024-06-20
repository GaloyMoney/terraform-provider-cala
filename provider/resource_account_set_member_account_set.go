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

var _ resource.Resource = &AccountSetMemberAccountSetResource{}
var _ resource.ResourceWithImportState = &AccountSetMemberAccountSetResource{}

func NewAccountSetMemberAccountSetResource() resource.Resource {
	return &AccountSetMemberAccountSetResource{}
}

type AccountSetMemberAccountSetResource struct {
	client *graphql.Client
}

type AccountSetMemberAccountSetResourceModel struct {
	AccountSetMemberId types.String `tfsdk:"id"`
	AccountSetId       types.String `tfsdk:"account_set_id"`
	MemberAccountSetId types.String `tfsdk:"member_account_set_id"`
}

func (r *AccountSetMemberAccountSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_set_member_account_set"
}

func (r *AccountSetMemberAccountSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Represents the membership of an account set in another account set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set.",
				Computed:            true,
			},
			"account_set_id": schema.StringAttribute{
				MarkdownDescription: "Id of the AccountSet",
				Required:            true,
			},
			"member_account_set_id": schema.StringAttribute{
				MarkdownDescription: "Id of the member AccountSet",
				Required: true,
			},
		},
	}
}

func (r *AccountSetMemberAccountSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// create
func (r *AccountSetMemberAccountSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AccountSetMemberAccountSetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := accountSetMemberAccountSetCreate(ctx, *r.client, data.AccountSetId.ValueString(), data.MemberAccountSetId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create account set member, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Added an account set to an account set")

  data.AccountSetMemberId = types.StringValue(fmt.Sprintf("account_set/%s/member_account_set/%s", data.AccountSetId.ValueString(), data.MemberAccountSetId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AccountSetMemberAccountSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *AccountSetMemberAccountSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *AccountSetMemberAccountSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AccountSetMemberAccountSetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	_, err := accountSetMemberAccountSetRemove(ctx, *r.client, data.AccountSetId.ValueString(), data.MemberAccountSetId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete account set member, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Removed an account set from an account set")

}

func (r *AccountSetMemberAccountSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}

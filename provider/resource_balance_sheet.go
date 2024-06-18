package provider

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &BalanceSheetResource{}
var _ resource.ResourceWithImportState = &BalanceSheetResource{}

func NewBalanceSheetResource() resource.Resource {
	return &BalanceSheetResource{}
}

type BalanceSheetResource struct {
	client *graphql.Client
}

type BalanceSheetResourceModel struct {
	JournalId               types.String `tfsdk:"journal_id"`
	AssetsAccountSetId      types.String `tfsdk:"assets_account_set_id"`
	LiabilitiesAccountSetId types.String `tfsdk:"liabilities_account_set_id"`
	Schedule1AccountSetId   types.String `tfsdk:"schedule1_account_set_id"`
	Schedule2AccountSetId   types.String `tfsdk:"schedule2_account_set_id"`
	Schedule3AccountSetId   types.String `tfsdk:"schedule3_account_set_id"`
	Schedule4AccountSetId   types.String `tfsdk:"schedule4_account_set_id"`
	Schedule5AccountSetId   types.String `tfsdk:"schedule5_account_set_id"`
	Schedule6AccountSetId   types.String `tfsdk:"schedule6_account_set_id"`
	Schedule7AccountSetId   types.String `tfsdk:"schedule7_account_set_id"`
	Schedule8AccountSetId   types.String `tfsdk:"schedule8_account_set_id"`
	Schedule9AccountSetId   types.String `tfsdk:"schedule9_account_set_id"`
	Schedule10AccountSetId  types.String `tfsdk:"schedule10_account_set_id"`
	Schedule11AccountSetId  types.String `tfsdk:"schedule11_account_set_id"`
	Schedule12AccountSetId  types.String `tfsdk:"schedule12_account_set_id"`
}

func (r *BalanceSheetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_balance_sheet"
}

func (r *BalanceSheetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cala balance sheet.",
		Attributes: map[string]schema.Attribute{
			"journal_id": schema.StringAttribute{
				MarkdownDescription: "ID of the journal associated with the balance sheet.",
				Required:            true,
			},
			"assets_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for assets.",
				Computed:            true,
			},
			"liabilities_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for liabilities.",
				Computed:            true,
			},
			"schedule1_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 1.",
				Computed:            true,
			},
			"schedule2_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 2.",
				Computed:            true,
			},
			"schedule3_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 3.",
				Computed:            true,
			},
			"schedule4_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 4.",
				Computed:            true,
			},
			"schedule5_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 5.",
				Computed:            true,
			},
			"schedule6_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 6.",
				Computed:            true,
			},
			"schedule7_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 7.",
				Computed:            true,
			},
			"schedule8_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 8.",
				Computed:            true,
			},
			"schedule9_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 9.",
				Computed:            true,
			},
			"schedule10_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 10.",
				Computed:            true,
			},
			"schedule11_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 11.",
				Computed:            true,
			},
			"schedule12_account_set_id": schema.StringAttribute{
				MarkdownDescription: "ID of the account set for schedule 12.",
				Computed:            true,
			},
		},
	}
}

func (r *BalanceSheetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BalanceSheetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *BalanceSheetResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	response, err := balanceSheetGet(ctx, *r.client, data.JournalId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Failed to get balance sheet", err.Error())
		return
	}

	balanceSheet := response.BalanceSheet.ByJournalId

	data.AssetsAccountSetId = types.StringValue(balanceSheet.Assets.AccountSetId)
	data.LiabilitiesAccountSetId = types.StringValue(balanceSheet.Liabilities.AccountSetId)
	data.Schedule1AccountSetId = types.StringValue(balanceSheet.Schedule1.AccountSetId)
	data.Schedule2AccountSetId = types.StringValue(balanceSheet.Schedule2.AccountSetId)
	data.Schedule3AccountSetId = types.StringValue(balanceSheet.Schedule3.AccountSetId)
	data.Schedule4AccountSetId = types.StringValue(balanceSheet.Schedule4.AccountSetId)
	data.Schedule5AccountSetId = types.StringValue(balanceSheet.Schedule5.AccountSetId)
	data.Schedule6AccountSetId = types.StringValue(balanceSheet.Schedule6.AccountSetId)
	data.Schedule7AccountSetId = types.StringValue(balanceSheet.Schedule7.AccountSetId)
	data.Schedule8AccountSetId = types.StringValue(balanceSheet.Schedule8.AccountSetId)
	data.Schedule9AccountSetId = types.StringValue(balanceSheet.Schedule9.AccountSetId)
	data.Schedule10AccountSetId = types.StringValue(balanceSheet.Schedule10.AccountSetId)
	data.Schedule11AccountSetId = types.StringValue(balanceSheet.Schedule11.AccountSetId)
	data.Schedule12AccountSetId = types.StringValue(balanceSheet.Schedule12.AccountSetId)
}

func (r *BalanceSheetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *BalanceSheetResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := BalanceSheetCreateInput{
		JournalId: data.JournalId.ValueString(),
	}

	// Create the balance sheet
	response, err := balanceSheetCreate(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Failed to create balance sheet", err.Error())
		return
	}

	balanceSheet := response.BalanceSheet.Create.BalanceSheet

	data.AssetsAccountSetId = types.StringValue(balanceSheet.Assets.AccountSetId)
	data.LiabilitiesAccountSetId = types.StringValue(balanceSheet.Liabilities.AccountSetId)
	data.Schedule1AccountSetId = types.StringValue(balanceSheet.Schedule1.AccountSetId)
	data.Schedule2AccountSetId = types.StringValue(balanceSheet.Schedule2.AccountSetId)
	data.Schedule3AccountSetId = types.StringValue(balanceSheet.Schedule3.AccountSetId)
	data.Schedule4AccountSetId = types.StringValue(balanceSheet.Schedule4.AccountSetId)
	data.Schedule5AccountSetId = types.StringValue(balanceSheet.Schedule5.AccountSetId)
	data.Schedule6AccountSetId = types.StringValue(balanceSheet.Schedule6.AccountSetId)
	data.Schedule7AccountSetId = types.StringValue(balanceSheet.Schedule7.AccountSetId)
	data.Schedule8AccountSetId = types.StringValue(balanceSheet.Schedule8.AccountSetId)
	data.Schedule9AccountSetId = types.StringValue(balanceSheet.Schedule9.AccountSetId)
	data.Schedule10AccountSetId = types.StringValue(balanceSheet.Schedule10.AccountSetId)
	data.Schedule11AccountSetId = types.StringValue(balanceSheet.Schedule11.AccountSetId)
	data.Schedule12AccountSetId = types.StringValue(balanceSheet.Schedule12.AccountSetId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BalanceSheetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *BalanceSheetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *BalanceSheetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError("ImportState is not implemented", "ImportState is not implemented for balance sheet resource")
}

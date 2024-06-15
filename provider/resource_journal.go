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

var _ resource.Resource = &JournalResource{}
var _ resource.ResourceWithImportState = &JournalResource{}

func NewJournalResource() resource.Resource {
	return &JournalResource{}
}

type JournalResource struct {
	client *graphql.Client
}

type JournalResourceModel struct {
	JournalId types.String `tfsdk:"id"`
	Name 		types.String `tfsdk:"name"`
	Status 		types.String `tfsdk:"status"`
}

func (r *JournalResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_journal"
}

func (r *JournalResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Cala journal.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the journal.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the journal.",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "status",
				Default:             stringdefault.StaticString("ACTIVE"),
				Computed:            true,
			},
		},
	}
}

func (r *JournalResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *JournalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *JournalResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	status, err := toStatus(data.Status.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Status", fmt.Sprintf("Unable to convert status to Status: %s", err))
		return
	}

	input := JournalCreateInput{
		JournalId:         data.JournalId.ValueString(),
		Name:              data.Name.ValueString(),
		Status:            status,
	}

	response, err := journalCreate(ctx, *r.client, input)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create account, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "created an account")

	journal := response.JournalCreate.Journal

	data.JournalId = types.StringValue(journal.Id)
	data.Name = types.StringValue(journal.Name)
	data.Status = types.StringValue(string(journal.Status))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JournalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (r *JournalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *JournalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *JournalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}

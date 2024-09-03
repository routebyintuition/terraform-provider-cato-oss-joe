package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cato_models "github.com/routebyintuition/cato-go-sdk/models"
)

var (
	_ resource.Resource              = &internetFwPolicyResource{}
	_ resource.ResourceWithConfigure = &internetFwPolicyResource{}
)

func NewInternetFwPolicyResource() resource.Resource {
	return &internetFwPolicyResource{}
}

type internetFwPolicyResource struct {
	info *catoClientData
}

type InternetFirewall_Policy_Audit struct {
	PublishedTime types.String `tfsdk:"publishedtime"`
	PublishedBy   types.String `tfsdk:"publishedby"`
}

type InternetFirewall_Policy_Revision struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Changes     types.Int64  `tfsdk:"changes"`
	CreatedTime types.String `tfsdk:"createdtime"`
	UpdatedTime types.String `tfsdk:"updatedtime"`
}

type InternetFirewallAddRuleInput struct {
	ID   types.String `tfsdk:"id"`
	Rule types.Object `tfsdk:"rule"` //InternetFirewall_Policy_Rules
	At   types.Object `tfsdk:"at"`   //*PolicyRulePositionInput
}

type InternetFirewallCreateRuleInput struct {
	Rule types.Object `tfsdk:"rule"` //Policy_Policy_InternetFirewall_Policy_Rules_Rule
	At   types.Object `tfsdk:"at"`   //*PolicyRulePositionInput
	// Publish types.Bool   `tfsdk:"publish"` TO BE REMOVED
	// ID types.String `tfsdk:"id"`
}

type PolicyRulePositionInput struct {
	// this needs to be an emum PolicyRulePositionEnum
	Position types.String `tfsdk:"position"`
	Ref      types.String `tfsdk:"ref"`
}

type InternetFirewall_Policy_Rules struct {
	Audit types.Object `tfsdk:"audit"` //Policy_Policy_InternetFirewall_Policy_Rules_Audit
	Rule  types.Object `tfsdk:"rule"`  //Policy_Policy_InternetFirewall_Policy_Rules_Rule
	// need to switch properties to enum slice
	Properties types.List `tfsdk:"properties"` //[]types.String
}

type Policy_Policy_InternetFirewall_Policy_Rules_Audit struct {
	UpdatedTime types.String `tfsdk:"updatedtime"`
	UpdatedBy   types.String `tfsdk:"updatedby"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	Source           types.Object `tfsdk:"source"` //Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source
	ConnectionOrigin types.String `tfsdk:"connectionorigin"`
	Country          types.List   `tfsdk:"country"` //[]Policy_Policy_InternetFirewall_Policy_Rules_Rule_Country
	Device           types.List   `tfsdk:"device"`  //[]Policy_Policy_InternetFirewall_Policy_Rules_Rule_Device
	// needs to be enum OperatingSystem
	DeviceOs    types.List   `tfsdk:"deviceos"`
	Destination types.Object `tfsdk:"destination"` //Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination
	Service     types.Object `tfsdk:"service"`     //Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service
	// needs to be enum InternetFirewallActionEnum
	Action     types.String `tfsdk:"action"`
	Tracking   types.Object `tfsdk:"tracking"`   //Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking
	Schedule   types.Object `tfsdk:"schedule"`   //Policy_Policy_InternetFirewall_Policy_Rules_Rule_Schedule
	Exceptions types.List   `tfsdk:"exceptions"` //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Exceptions
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Country struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Device struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination struct {
	Application            types.List `tfsdk:"application"`            //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_Application
	CustomApp              types.List `tfsdk:"customapp"`              //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_CustomApp
	AppCategory            types.List `tfsdk:"appcategory"`            //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_AppCategory
	CustomCategory         types.List `tfsdk:"customcategory"`         //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_CustomCategory
	SanctionedAppsCategory types.List `tfsdk:"sanctionedappscategory"` //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_SanctionedAppsCategory
	Country                types.List `tfsdk:"country"`                //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_Country
	Domain                 types.List `tfsdk:"domain"`
	Fqdn                   types.List `tfsdk:"fqdn"`
	IP                     types.List `tfsdk:"ip"`
	Subnet                 types.List `tfsdk:"subnet"`
	IPRange                types.List `tfsdk:"iprange"`       //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_IPRange
	GlobalIPRange          types.List `tfsdk:"globaliprange"` //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_GlobalIPRange
	RemoteAsn              types.List `tfsdk:"remoteasn"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_Application struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_CustomApp struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_AppCategory struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_CustomCategory struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_SanctionedAppsCategory struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_Country struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_IPRange struct {
	From types.String `tfsdk:"from"`
	To   types.String `tfsdk:"to"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_GlobalIPRange struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source struct {
	IP                types.List `tfsdk:"ip"`
	Host              types.List `tfsdk:"host"` //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Host
	Site              types.List `tfsdk:"site"` //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Site
	Subnet            types.List `tfsdk:"subnet"`
	IPRange           types.List `tfsdk:"iprange"`           //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_IPRange
	GlobalIPRange     types.List `tfsdk:"globaliprange"`     //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_GlobalIPRange
	NetworkInterface  types.List `tfsdk:"networkinterface"`  //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_NetworkInterface
	SiteNetworkSubnet types.List `tfsdk:"sitenetworksubnet"` //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_SiteNetworkSubnet
	FloatingSubnet    types.List `tfsdk:"floatingsubnet"`    //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_FloatingSubnet
	User              types.List `tfsdk:"user"`              //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_User
	UsersGroup        types.List `tfsdk:"usersgroup"`        //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_UsersGroup
	Group             types.List `tfsdk:"group"`             //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Group
	SystemGroup       types.List `tfsdk:"systemgroup"`       //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_SystemGroup
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Host struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Site struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_IPRange struct {
	From types.String `tfsdk:"from"`
	To   types.String `tfsdk:"to"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_GlobalIPRange struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_NetworkInterface struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_SiteNetworkSubnet struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_User struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_UsersGroup struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Group struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_SystemGroup struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_FloatingSubnet struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service struct {
	Standard types.List `tfsdk:"standard"` //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service_Standard
	// To be enabled in later version
	// Custom types.List `tfsdk:"custom"` //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service_Custom
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service_Standard struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service_Custom struct {
	Port      types.List   `tfsdk:"port"`
	PortRange types.Object `tfsdk:"portrange"` //*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service_Custom_PortRange
	Protocol  types.String `tfsdk:"protocol"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service_Custom_PortRange struct {
	From types.String `tfsdk:"from"`
	To   types.String `tfsdk:"to"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking struct {
	Event types.Object `tfsdk:"event"` //Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Event
	Alert types.Object `tfsdk:"alert"` //Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Alert
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Event struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Alert struct {
	Enabled           types.Bool   `tfsdk:"enabled"`
	Frequency         types.String `tfsdk:"frequency"`
	SubscriptionGroup types.List   `tfsdk:"subscriptiongroup"` //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Alert_SubscriptionGroup
	Webhook           types.List   `tfsdk:"webhook"`           //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Alert_Webhook
	MailingList       types.List   `tfsdk:"mailinglist"`       //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Alert_MailingList
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Alert_SubscriptionGroup struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Alert_Webhook struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Alert_MailingList struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Schedule struct {
	ActiveOn        types.String `tfsdk:"activeon"`
	CustomTimeframe types.Object `tfsdk:"customtimeframe"` //*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Schedule_CustomTimeframe
	CustomRecurring types.Object `tfsdk:"customrecurring"` //*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Schedule_CustomRecurring
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Schedule_CustomTimeframe struct {
	From types.String `tfsdk:"from"`
	To   types.String `tfsdk:"to"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Schedule_CustomRecurring struct {
	From types.String `tfsdk:"from"`
	To   types.String `tfsdk:"to"`
	// needs to be an enum cato_query_models.DayOfWeek
	Days types.List `tfsdk:"days"` //[]DayOfWeek
}

type DayOfWeek types.String

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Exceptions struct {
	Name   types.String `tfsdk:"name"`
	Source types.Object `tfsdk:"source"` //Policy_Policy_InternetFirewall_Policy_Rules_Rule_Exceptions_Source
	// neeeds to be enum
	DeviceOs    types.List `tfsdk:"deviceos"`    //[]OperatingSystem
	Country     types.List `tfsdk:"country"`     //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Exceptions_Country
	Device      types.List `tfsdk:"device"`      //[]*Policy_Policy_InternetFirewall_Policy_Rules_Rule_Exceptions_Device
	Destination types.List `tfsdk:"destination"` //Policy_Policy_InternetFirewall_Policy_Rules_Rule_Exceptions_Destination
	// needs to be enum ConnectionOriginEnum
	ConnectionOrigin types.String `tfsdk:"connectionorigin"`
}

type OperatingSystem types.String

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Exceptions_Country struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Exceptions_Source struct {
	IP     types.List `tfsdk:"ip"`
	Subnet types.List `tfsdk:"subnet"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Exceptions_Device struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Exceptions_Destination struct {
	Domain    types.List `tfsdk:"domain"`
	Fqdn      types.List `tfsdk:"fqdn"`
	IP        types.List `tfsdk:"ip"`
	Subnet    types.List `tfsdk:"subnet"`
	RemoteAsn types.List `tfsdk:"remoteasn"`
}

type Policy_Policy_InternetFirewall_Policy_Rules_Rule_Section struct {
	By    types.String `tfsdk:"by"`
	Input types.String `tfsdk:"input"`
}

func (r *internetFwPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_if_policy"
}

func (r *internetFwPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"at": schema.SingleNestedAttribute{
				Description: "at",
				Required:    false,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"position": schema.StringAttribute{
						Description: "position",
						Required:    false,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"ref": schema.StringAttribute{
						Description: "ref",
						Required:    false,
						Optional:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
				},
			},
			"rule": schema.SingleNestedAttribute{
				Description: "rule item",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed: true,
						Optional: false,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Description: "Rule name",
						Required:    true,
					},
					"description": schema.StringAttribute{
						Description: "Rule description",
						Required:    false,
						Optional:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "enabled",
						Required:    true,
					},
					"source": schema.SingleNestedAttribute{
						Required: false,
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"ip": schema.ListAttribute{
								Description: "ip",
								ElementType: types.StringType,
								Required:    false,
								Optional:    true,
							},
							"host": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"site": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"subnet": schema.ListAttribute{
								ElementType: types.StringType,
								Description: "subnet",
								Required:    false,
								Optional:    true,
							},
							"iprange": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"from": schema.StringAttribute{
											Description: "from",
											Required:    false,
											Optional:    true,
										},
										"to": schema.StringAttribute{
											Description: "to",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"globaliprange": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"networkinterface": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"sitenetworksubnet": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"floatingsubnet": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"user": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"usersgroup": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"group": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"systemgroup": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
						},
					},
					"connectionorigin": schema.StringAttribute{
						Description: "connectionOrigin",
						Optional:    true,
						Required:    false,
					},
					"country": schema.ListNestedAttribute{
						Required: false,
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"by": schema.StringAttribute{
									Description: "by",
									Required:    false,
									Optional:    true,
								},
								"input": schema.StringAttribute{
									Description: "input",
									Required:    false,
									Optional:    true,
								},
							},
						},
					},
					"device": schema.ListAttribute{
						ElementType: types.StringType,
						Description: "device",
						Optional:    true,
						// Required:    true,
					},
					"deviceos": schema.ListAttribute{
						ElementType: types.StringType,
						Description: "deviceOS",
						Optional:    true,
						// Required:    true,
					},
					"destination": schema.SingleNestedAttribute{
						Optional: true,
						Required: false,
						Attributes: map[string]schema.Attribute{
							"application": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"customapp": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"appcategory": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"customcategory": schema.ListNestedAttribute{
								Description: "customCategory",
								Required:    false,
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"sanctionedappscategory": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"country": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"domain": schema.ListAttribute{
								ElementType: types.StringType,
								Description: "domain",
								Required:    false,
								Optional:    true,
							},
							"fqdn": schema.ListAttribute{
								ElementType: types.StringType,
								Description: "fqdn",
								Required:    false,
								Optional:    true,
							},
							"ip": schema.ListAttribute{
								ElementType: types.StringType,
								Description: "ip",
								Required:    false,
								Optional:    true,
							},
							"subnet": schema.ListAttribute{
								ElementType: types.StringType,
								Description: "subnet",
								Required:    false,
								Optional:    true,
							},
							"iprange": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"from": schema.StringAttribute{
											Description: "from",
											Required:    false,
											Optional:    true,
										},
										"to": schema.StringAttribute{
											Description: "to",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"globaliprange": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
							"remoteasn": schema.ListAttribute{
								ElementType: types.StringType,
								Description: "remoteAsn",
								Required:    false,
								Optional:    true,
							},
						},
					},
					"service": schema.SingleNestedAttribute{
						Required: false,
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"standard": schema.ListNestedAttribute{
								Required: false,
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"by": schema.StringAttribute{
											Description: "by",
											Required:    false,
											Optional:    true,
										},
										"input": schema.StringAttribute{
											Description: "input",
											Required:    false,
											Optional:    true,
										},
									},
								},
							},
						},
					},
					"action": schema.StringAttribute{
						Description: "action",
						Required:    true,
					},
					"tracking": schema.SingleNestedAttribute{
						Required: false,
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"event": schema.SingleNestedAttribute{
								Description: "event",
								Required:    true,
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "enabled",
										Required:    true,
									},
								},
							},
							"alert": schema.SingleNestedAttribute{
								Description: "alert",
								Required:    false,
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										Description: "enabled",
										Required:    true,
									},
									"frequency": schema.StringAttribute{
										Description: "frequency",
										Required:    true,
									},
									"subscriptiongroup": schema.ListAttribute{
										ElementType: types.StringType,
										Description: "subscriptionGroup",
										Required:    false,
										Optional:    true,
									},
									"webhook": schema.ListAttribute{
										ElementType: types.StringType,
										Description: "webhook",
										Required:    false,
										Optional:    true,
									},
									"mailinglist": schema.ListAttribute{
										ElementType: types.StringType,
										Description: "mailingList",
										Required:    false,
										Optional:    true,
									},
								},
							},
						},
					},
					"schedule": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"activeon": schema.StringAttribute{
								Description: "activeOn",
								Required:    false,
								Optional:    true,
							},
							"customtimeframe": schema.StringAttribute{
								Description: "customtimeframe",
								Required:    false,
								Optional:    true,
							},
							"customrecurring": schema.StringAttribute{
								Description: "customrecurring",
								Required:    false,
								Optional:    true,
							},
						},
					},
					"exceptions": schema.ListAttribute{
						ElementType: types.StringType,
						Description: "exceptions",
						Required:    false,
						Optional:    true,
					},
				},
			},
		},
	}
}

func (d *internetFwPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.info = req.ProviderData.(*catoClientData)

}

func (r *internetFwPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan InternetFirewallCreateRuleInput
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//initiate input
	input := cato_models.InternetFirewallAddRuleInput{
		At:   &cato_models.PolicyRulePositionInput{},
		Rule: &cato_models.InternetFirewallAddRuleDataInput{},
	}

	//retrieve & setting position
	positionInput := PolicyRulePositionInput{}
	diags = plan.At.As(ctx, &positionInput, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	input.At.Position = (*cato_models.PolicyRulePositionEnum)(positionInput.Position.ValueStringPointer())
	input.At.Ref = positionInput.Ref.ValueStringPointer()

	// retrieve & setting rule
	ruleInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule{}
	diags = plan.Rule.As(ctx, &ruleInput, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve & setting source
	if !ruleInput.Source.IsNull() {
		input.Rule.Source = &cato_models.InternetFirewallSourceInput{}

		sourceInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source{}
		diags = ruleInput.Source.As(ctx, &sourceInput, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting source IP
		diags = sourceInput.IP.ElementsAs(ctx, &input.Rule.Source.IP, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting source subnet
		diags = sourceInput.Subnet.ElementsAs(ctx, &input.Rule.Source.Subnet, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting source host
		elementsSourceHostInput := make([]types.Object, 0, len(sourceInput.Host.Elements()))
		diags = sourceInput.Host.ElementsAs(ctx, &elementsSourceHostInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceHostInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Host
		for _, item := range elementsSourceHostInput {
			diags = item.As(ctx, &itemSourceHostInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.Host = append(input.Rule.Source.Host, &cato_models.HostRefInput{
				By:    cato_models.ObjectRefBy(itemSourceHostInput.By.ValueString()),
				Input: itemSourceHostInput.Input.ValueString(),
			})
		}

		// retrieve & setting source site
		elementsSourceSiteInput := make([]types.Object, 0, len(sourceInput.Site.Elements()))
		diags = sourceInput.Site.ElementsAs(ctx, &elementsSourceSiteInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceSiteInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Site
		for _, item := range elementsSourceSiteInput {
			diags = item.As(ctx, &itemSourceSiteInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.Site = append(input.Rule.Source.Site, &cato_models.SiteRefInput{
				By:    cato_models.ObjectRefBy(itemSourceSiteInput.By.ValueString()),
				Input: itemSourceSiteInput.Input.ValueString(),
			})
		}

		// retrieve & setting source ip range
		elementsSourceIPRangeInput := make([]types.Object, 0, len(sourceInput.IPRange.Elements()))
		diags = sourceInput.IPRange.ElementsAs(ctx, &elementsSourceIPRangeInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceIPRangeInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_IPRange
		for _, item := range elementsSourceIPRangeInput {
			diags = item.As(ctx, &itemSourceIPRangeInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.IPRange = append(input.Rule.Source.IPRange, &cato_models.IPAddressRangeInput{
				From: itemSourceIPRangeInput.From.ValueString(),
				To:   itemSourceIPRangeInput.To.ValueString(),
			})
		}

		// retrieve & setting source global ip range
		elementsSourceGlobalIPRangeInput := make([]types.Object, 0, len(sourceInput.GlobalIPRange.Elements()))
		diags = sourceInput.GlobalIPRange.ElementsAs(ctx, &elementsSourceGlobalIPRangeInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceGlobalIPRangeInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_GlobalIPRange
		for _, item := range elementsSourceGlobalIPRangeInput {
			diags = item.As(ctx, &itemSourceGlobalIPRangeInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.GlobalIPRange = append(input.Rule.Source.GlobalIPRange, &cato_models.GlobalIPRangeRefInput{
				By:    cato_models.ObjectRefBy(itemSourceGlobalIPRangeInput.By.ValueString()),
				Input: itemSourceGlobalIPRangeInput.Input.ValueString(),
			})
		}

		// retrieve & setting source network interface
		elementsSourceNetworkInterfaceInput := make([]types.Object, 0, len(sourceInput.NetworkInterface.Elements()))
		diags = sourceInput.NetworkInterface.ElementsAs(ctx, &elementsSourceNetworkInterfaceInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceNetworkInterfaceInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_NetworkInterface
		for _, item := range elementsSourceNetworkInterfaceInput {
			diags = item.As(ctx, &itemSourceNetworkInterfaceInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.NetworkInterface = append(input.Rule.Source.NetworkInterface, &cato_models.NetworkInterfaceRefInput{
				By:    cato_models.ObjectRefBy(itemSourceNetworkInterfaceInput.By.ValueString()),
				Input: itemSourceNetworkInterfaceInput.Input.ValueString(),
			})
		}

		// retrieve & setting source site network subnet
		elementsSourceSiteNetworkSubnetInput := make([]types.Object, 0, len(sourceInput.SiteNetworkSubnet.Elements()))
		diags = sourceInput.SiteNetworkSubnet.ElementsAs(ctx, &elementsSourceSiteNetworkSubnetInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceSiteNetworkSubnetInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_SiteNetworkSubnet
		for _, item := range elementsSourceSiteNetworkSubnetInput {
			diags = item.As(ctx, &itemSourceSiteNetworkSubnetInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.SiteNetworkSubnet = append(input.Rule.Source.SiteNetworkSubnet, &cato_models.SiteNetworkSubnetRefInput{
				By:    cato_models.ObjectRefBy(itemSourceSiteNetworkSubnetInput.By.ValueString()),
				Input: itemSourceSiteNetworkSubnetInput.Input.ValueString(),
			})
		}

		// retrieve & setting source floating subnet
		elementsSourceFloatingSubnetInput := make([]types.Object, 0, len(sourceInput.FloatingSubnet.Elements()))
		diags = sourceInput.FloatingSubnet.ElementsAs(ctx, &elementsSourceFloatingSubnetInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceFloatingSubnetInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_FloatingSubnet
		for _, item := range elementsSourceFloatingSubnetInput {
			diags = item.As(ctx, &itemSourceFloatingSubnetInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.FloatingSubnet = append(input.Rule.Source.FloatingSubnet, &cato_models.FloatingSubnetRefInput{
				By:    cato_models.ObjectRefBy(itemSourceFloatingSubnetInput.By.ValueString()),
				Input: itemSourceFloatingSubnetInput.Input.ValueString(),
			})
		}

		// retrieve & setting source user
		elementsSourceUserInput := make([]types.Object, 0, len(sourceInput.User.Elements()))
		diags = sourceInput.User.ElementsAs(ctx, &elementsSourceUserInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceUserInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_User
		for _, item := range elementsSourceUserInput {
			diags = item.As(ctx, &itemSourceUserInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.User = append(input.Rule.Source.User, &cato_models.UserRefInput{
				By:    cato_models.ObjectRefBy(itemSourceUserInput.By.ValueString()),
				Input: itemSourceUserInput.Input.ValueString(),
			})
		}

		// retrieve & setting source users group
		elementsSourceUsersGroupInput := make([]types.Object, 0, len(sourceInput.UsersGroup.Elements()))
		diags = sourceInput.UsersGroup.ElementsAs(ctx, &elementsSourceUsersGroupInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceUsersGroupInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_UsersGroup
		for _, item := range elementsSourceUsersGroupInput {
			diags = item.As(ctx, &itemSourceUsersGroupInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.UsersGroup = append(input.Rule.Source.UsersGroup, &cato_models.UsersGroupRefInput{
				By:    cato_models.ObjectRefBy(itemSourceUsersGroupInput.By.ValueString()),
				Input: itemSourceUsersGroupInput.Input.ValueString(),
			})
		}

		// retrieve & setting source group
		elementsSourceGroupInput := make([]types.Object, 0, len(sourceInput.Group.Elements()))
		diags = sourceInput.Group.ElementsAs(ctx, &elementsSourceGroupInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceGroupInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Group
		for _, item := range elementsSourceGroupInput {
			diags = item.As(ctx, &itemSourceGroupInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.Group = append(input.Rule.Source.Group, &cato_models.GroupRefInput{
				By:    cato_models.ObjectRefBy(itemSourceGroupInput.By.ValueString()),
				Input: itemSourceGroupInput.Input.ValueString(),
			})
		}

		// retrieve & setting source system group
		elementsSourceSystemGroupInput := make([]types.Object, 0, len(sourceInput.SystemGroup.Elements()))
		diags = sourceInput.SystemGroup.ElementsAs(ctx, &elementsSourceSystemGroupInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceSystemGroupInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_SystemGroup
		for _, item := range elementsSourceSystemGroupInput {
			diags = item.As(ctx, &itemSourceSystemGroupInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.SystemGroup = append(input.Rule.Source.SystemGroup, &cato_models.SystemGroupRefInput{
				By:    cato_models.ObjectRefBy(itemSourceSystemGroupInput.By.ValueString()),
				Input: itemSourceSystemGroupInput.Input.ValueString(),
			})
		}
	}

	// retrieve & setting country
	elementsCountryInput := make([]types.Object, 0, len(ruleInput.Country.Elements()))
	diags = ruleInput.Country.ElementsAs(ctx, &elementsCountryInput, false)
	resp.Diagnostics.Append(diags...)

	var itemCountryInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Country
	for _, item := range elementsCountryInput {
		diags = item.As(ctx, &itemCountryInput, basetypes.ObjectAsOptions{})
		input.Rule.Country = append(input.Rule.Country, &cato_models.CountryRefInput{
			By:    cato_models.ObjectRefBy(itemCountryInput.By.ValueString()),
			Input: itemCountryInput.Input.ValueString(),
		})
	}

	// retrieve & setting device
	elementsDeviceInput := make([]types.Object, 0, len(ruleInput.Device.Elements()))
	diags = ruleInput.Device.ElementsAs(ctx, &elementsDeviceInput, false)
	resp.Diagnostics.Append(diags...)

	var itemDeviceInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Device
	for _, item := range elementsDeviceInput {
		diags = item.As(ctx, &itemDeviceInput, basetypes.ObjectAsOptions{})
		input.Rule.Device = append(input.Rule.Device, &cato_models.DeviceProfileRefInput{
			By:    cato_models.ObjectRefBy(itemDeviceInput.By.ValueString()),
			Input: itemDeviceInput.Input.ValueString(),
		})
	}

	// retrieve & setting device OS
	diags = ruleInput.DeviceOs.ElementsAs(ctx, &input.Rule.DeviceOs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve & setting destination
	if !ruleInput.Destination.IsNull() {
		input.Rule.Destination = &cato_models.InternetFirewallDestinationInput{}

		destinationInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination{}
		diags = ruleInput.Destination.As(ctx, &destinationInput, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination IP
		diags = destinationInput.IP.ElementsAs(ctx, &input.Rule.Destination.IP, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination subnet
		diags = destinationInput.Subnet.ElementsAs(ctx, &input.Rule.Destination.Subnet, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination domain
		diags = destinationInput.Domain.ElementsAs(ctx, &input.Rule.Destination.Domain, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination fqdn
		diags = destinationInput.Fqdn.ElementsAs(ctx, &input.Rule.Destination.Fqdn, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination remote asn
		diags = destinationInput.RemoteAsn.ElementsAs(ctx, &input.Rule.Destination.RemoteAsn, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination application
		elementsDestinationApplicationInput := make([]types.Object, 0, len(destinationInput.Application.Elements()))
		diags = destinationInput.Application.ElementsAs(ctx, &elementsDestinationApplicationInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationApplicationInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_Application
		for _, item := range elementsDestinationApplicationInput {
			diags = item.As(ctx, &itemDestinationApplicationInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.Application = append(input.Rule.Destination.Application, &cato_models.ApplicationRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationApplicationInput.By.ValueString()),
				Input: itemDestinationApplicationInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination custom app
		elementsDestinationCustomAppInput := make([]types.Object, 0, len(destinationInput.CustomApp.Elements()))
		diags = destinationInput.CustomApp.ElementsAs(ctx, &elementsDestinationCustomAppInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationCustomAppInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_CustomApp
		for _, item := range elementsDestinationCustomAppInput {
			diags = item.As(ctx, &itemDestinationCustomAppInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.CustomApp = append(input.Rule.Destination.CustomApp, &cato_models.CustomApplicationRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationCustomAppInput.By.ValueString()),
				Input: itemDestinationCustomAppInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination ip range
		elementsDestinationIPRangeInput := make([]types.Object, 0, len(destinationInput.IPRange.Elements()))
		diags = destinationInput.IPRange.ElementsAs(ctx, &elementsDestinationIPRangeInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationIPRangeInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_IPRange
		for _, item := range elementsDestinationIPRangeInput {
			diags = item.As(ctx, &itemDestinationIPRangeInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.IPRange = append(input.Rule.Destination.IPRange, &cato_models.IPAddressRangeInput{
				From: itemDestinationIPRangeInput.From.ValueString(),
				To:   itemDestinationIPRangeInput.To.ValueString(),
			})
		}

		// retrieve & setting destination global ip range
		elementsDestinationGlobalIPRangeInput := make([]types.Object, 0, len(destinationInput.GlobalIPRange.Elements()))
		diags = destinationInput.GlobalIPRange.ElementsAs(ctx, &elementsDestinationGlobalIPRangeInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationGlobalIPRangeInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_GlobalIPRange
		for _, item := range elementsDestinationGlobalIPRangeInput {
			diags = item.As(ctx, &itemDestinationGlobalIPRangeInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.GlobalIPRange = append(input.Rule.Destination.GlobalIPRange, &cato_models.GlobalIPRangeRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationGlobalIPRangeInput.By.ValueString()),
				Input: itemDestinationGlobalIPRangeInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination app category
		elementsDestinationAppCategoryInput := make([]types.Object, 0, len(destinationInput.AppCategory.Elements()))
		diags = destinationInput.AppCategory.ElementsAs(ctx, &elementsDestinationAppCategoryInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationAppCategoryInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_AppCategory
		for _, item := range elementsDestinationAppCategoryInput {
			diags = item.As(ctx, &itemDestinationAppCategoryInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.AppCategory = append(input.Rule.Destination.AppCategory, &cato_models.ApplicationCategoryRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationAppCategoryInput.By.ValueString()),
				Input: itemDestinationAppCategoryInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination custom app category
		elementsDestinationCustomCategoryInput := make([]types.Object, 0, len(destinationInput.CustomCategory.Elements()))
		diags = destinationInput.CustomCategory.ElementsAs(ctx, &elementsDestinationCustomCategoryInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationCustomCategoryInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_CustomCategory
		for _, item := range elementsDestinationCustomCategoryInput {
			diags = item.As(ctx, &itemDestinationCustomCategoryInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.CustomCategory = append(input.Rule.Destination.CustomCategory, &cato_models.CustomCategoryRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationCustomCategoryInput.By.ValueString()),
				Input: itemDestinationCustomCategoryInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination sanctionned apps category
		elementsDestinationSanctionedAppsCategoryInput := make([]types.Object, 0, len(destinationInput.SanctionedAppsCategory.Elements()))
		diags = destinationInput.SanctionedAppsCategory.ElementsAs(ctx, &elementsDestinationSanctionedAppsCategoryInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationSanctionedAppsCategoryInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_SanctionedAppsCategory
		for _, item := range elementsDestinationSanctionedAppsCategoryInput {
			diags = item.As(ctx, &itemDestinationSanctionedAppsCategoryInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.SanctionedAppsCategory = append(input.Rule.Destination.SanctionedAppsCategory, &cato_models.SanctionedAppsCategoryRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationSanctionedAppsCategoryInput.By.ValueString()),
				Input: itemDestinationSanctionedAppsCategoryInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination country
		elementsDestinationCountryInput := make([]types.Object, 0, len(destinationInput.Country.Elements()))
		diags = destinationInput.Country.ElementsAs(ctx, &elementsDestinationCountryInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationCountryInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_Country
		for _, item := range elementsDestinationCountryInput {
			diags = item.As(ctx, &itemDestinationCountryInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.Country = append(input.Rule.Destination.Country, &cato_models.CountryRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationCountryInput.By.ValueString()),
				Input: itemDestinationCountryInput.Input.ValueString(),
			})
		}
	}

	// retrieve & setting service
	if !ruleInput.Service.IsNull() {
		input.Rule.Service = &cato_models.InternetFirewallServiceTypeInput{}

		serviceInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service{}
		diags = ruleInput.Service.As(ctx, &serviceInput, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting service standard
		elementsServiceStandardInput := make([]types.Object, 0, len(serviceInput.Standard.Elements()))
		diags = serviceInput.Standard.ElementsAs(ctx, &elementsServiceStandardInput, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		var itemServiceStandardInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service_Standard
		for _, item := range elementsServiceStandardInput {
			diags = item.As(ctx, &itemServiceStandardInput, basetypes.ObjectAsOptions{})
			input.Rule.Service.Standard = append(input.Rule.Service.Standard, &cato_models.ServiceRefInput{
				By:    cato_models.ObjectRefBy(itemServiceStandardInput.By.ValueString()),
				Input: itemServiceStandardInput.Input.ValueString(),
			})
		}
	}

	// retrieve & setting tracking
	if !ruleInput.Tracking.IsNull() {
		input.Rule.Tracking = &cato_models.PolicyTrackingInput{
			Event: &cato_models.PolicyRuleTrackingEventInput{},
		}

		trackingInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking{}
		diags = ruleInput.Tracking.As(ctx, &trackingInput, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting tracking event
		trackingEventInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Event{}
		diags = trackingInput.Event.As(ctx, &trackingEventInput, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.Rule.Tracking.Event.Enabled = trackingEventInput.Enabled.ValueBool()

		// retrieve & setting tracking Alert
		if !trackingInput.Alert.IsNull() {
			input.Rule.Tracking.Alert = &cato_models.PolicyRuleTrackingAlertInput{}

			trackingAlertInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Alert{}
			diags = trackingInput.Alert.As(ctx, &trackingAlertInput, basetypes.ObjectAsOptions{})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			input.Rule.Tracking.Alert.Enabled = trackingAlertInput.Enabled.ValueBool()
			input.Rule.Tracking.Alert.Frequency = cato_models.PolicyRuleTrackingFrequencyEnum(trackingAlertInput.Frequency.ValueString())

			// retrieve & setting tracking alert subscription group
			diags = trackingAlertInput.SubscriptionGroup.ElementsAs(ctx, &input.Rule.Tracking.Alert.SubscriptionGroup, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			// retrieve & setting tracking alert webhook
			diags = trackingAlertInput.Webhook.ElementsAs(ctx, &input.Rule.Tracking.Alert.Webhook, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			// retrieve & setting tracking alert mailing list
			diags = trackingAlertInput.MailingList.ElementsAs(ctx, &input.Rule.Tracking.Alert.MailingList, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	// settings other rule attributes
	input.Rule.Name = ruleInput.Name.ValueString()
	input.Rule.Description = ruleInput.Description.ValueString()
	input.Rule.Enabled = ruleInput.Enabled.ValueBool()
	input.Rule.ConnectionOrigin = cato_models.ConnectionOriginEnum(ruleInput.ConnectionOrigin.ValueString())
	input.Rule.Action = cato_models.InternetFirewallActionEnum(ruleInput.Action.ValueString())

	//creating new rule
	policyChange, err := r.info.catov2.PolicyInternetFirewallAddRule(ctx, input, r.info.AccountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Catov2 API PolicyInternetFirewallAddRule error",
			err.Error(),
		)
		return
	}

	//publishing new rule
	tflog.Info(ctx, "publishing new rule")
	publishDataIfEnabled := &cato_models.PolicyPublishRevisionInput{}
	_, err = r.info.catov2.PolicyInternetFirewallPublishPolicyRevision(ctx, &cato_models.InternetFirewallPolicyMutationInput{}, publishDataIfEnabled, r.info.AccountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Catov2 API PolicyInternetFirewallPublishPolicyRevision error",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// overiding state with rule id
	resp.State.SetAttribute(
		ctx,
		path.Root("rule").AtName("id"),
		policyChange.GetPolicy().GetInternetFirewall().GetAddRule().Rule.GetRule().ID)

}

func (r *internetFwPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *internetFwPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan InternetFirewallCreateRuleInput
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//initiate input (whole schema is defined to have empty value)
	input := cato_models.InternetFirewallUpdateRuleInput{
		Rule: &cato_models.InternetFirewallUpdateRuleDataInput{
			Source: &cato_models.InternetFirewallSourceUpdateInput{
				// IP:                []string{},
				Host:              []*cato_models.HostRefInput{},
				Site:              []*cato_models.SiteRefInput{},
				Subnet:            []string{},
				IPRange:           []*cato_models.IPAddressRangeInput{},
				GlobalIPRange:     []*cato_models.GlobalIPRangeRefInput{},
				NetworkInterface:  []*cato_models.NetworkInterfaceRefInput{},
				SiteNetworkSubnet: []*cato_models.SiteNetworkSubnetRefInput{},
				FloatingSubnet:    []*cato_models.FloatingSubnetRefInput{},
				User:              []*cato_models.UserRefInput{},
				UsersGroup:        []*cato_models.UsersGroupRefInput{},
				Group:             []*cato_models.GroupRefInput{},
				SystemGroup:       []*cato_models.SystemGroupRefInput{},
			},
			Country:  []*cato_models.CountryRefInput{},
			Device:   []*cato_models.DeviceProfileRefInput{},
			DeviceOs: []cato_models.OperatingSystem{},
			Destination: &cato_models.InternetFirewallDestinationUpdateInput{
				Application:            []*cato_models.ApplicationRefInput{},
				CustomApp:              []*cato_models.CustomApplicationRefInput{},
				AppCategory:            []*cato_models.ApplicationCategoryRefInput{},
				CustomCategory:         []*cato_models.CustomCategoryRefInput{},
				SanctionedAppsCategory: []*cato_models.SanctionedAppsCategoryRefInput{},
				Country:                []*cato_models.CountryRefInput{},
				Domain:                 []string{},
				Fqdn:                   []string{},
				IP:                     []string{},
				Subnet:                 []string{},
				IPRange:                []*cato_models.IPAddressRangeInput{},
				GlobalIPRange:          []*cato_models.GlobalIPRangeRefInput{},
				RemoteAsn:              []string{},
			},
			Service: &cato_models.InternetFirewallServiceTypeUpdateInput{
				Standard: []*cato_models.ServiceRefInput{},
			},
			Tracking: &cato_models.PolicyTrackingUpdateInput{
				Event: &cato_models.PolicyRuleTrackingEventUpdateInput{},
				Alert: &cato_models.PolicyRuleTrackingAlertUpdateInput{},
			},
		},
	}

	// retrieve & setting rule
	ruleInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule{}
	diags = plan.Rule.As(ctx, &ruleInput, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve & setting source
	if !ruleInput.Source.IsNull() {
		sourceInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source{}
		diags = ruleInput.Source.As(ctx, &sourceInput, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting source IP
		diags = sourceInput.IP.ElementsAs(ctx, &input.Rule.Source.IP, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting source subnet
		diags = sourceInput.Subnet.ElementsAs(ctx, &input.Rule.Source.Subnet, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting source host
		elementsSourceHostInput := make([]types.Object, 0, len(sourceInput.Host.Elements()))
		diags = sourceInput.Host.ElementsAs(ctx, &elementsSourceHostInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceHostInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Host
		for _, item := range elementsSourceHostInput {
			diags = item.As(ctx, &itemSourceHostInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.Host = append(input.Rule.Source.Host, &cato_models.HostRefInput{
				By:    cato_models.ObjectRefBy(itemSourceHostInput.By.ValueString()),
				Input: itemSourceHostInput.Input.ValueString(),
			})
		}

		// retrieve & setting source site
		elementsSourceSiteInput := make([]types.Object, 0, len(sourceInput.Site.Elements()))
		diags = sourceInput.Site.ElementsAs(ctx, &elementsSourceSiteInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceSiteInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Site
		for _, item := range elementsSourceSiteInput {
			diags = item.As(ctx, &itemSourceSiteInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.Site = append(input.Rule.Source.Site, &cato_models.SiteRefInput{
				By:    cato_models.ObjectRefBy(itemSourceSiteInput.By.ValueString()),
				Input: itemSourceSiteInput.Input.ValueString(),
			})
		}

		// retrieve & setting source ip range
		elementsSourceIPRangeInput := make([]types.Object, 0, len(sourceInput.IPRange.Elements()))
		diags = sourceInput.IPRange.ElementsAs(ctx, &elementsSourceIPRangeInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceIPRangeInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_IPRange
		for _, item := range elementsSourceIPRangeInput {
			diags = item.As(ctx, &itemSourceIPRangeInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.IPRange = append(input.Rule.Source.IPRange, &cato_models.IPAddressRangeInput{
				From: itemSourceIPRangeInput.From.ValueString(),
				To:   itemSourceIPRangeInput.To.ValueString(),
			})
		}

		// retrieve & setting source global ip range
		elementsSourceGlobalIPRangeInput := make([]types.Object, 0, len(sourceInput.GlobalIPRange.Elements()))
		diags = sourceInput.GlobalIPRange.ElementsAs(ctx, &elementsSourceGlobalIPRangeInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceGlobalIPRangeInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_GlobalIPRange
		for _, item := range elementsSourceGlobalIPRangeInput {
			diags = item.As(ctx, &itemSourceGlobalIPRangeInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.GlobalIPRange = append(input.Rule.Source.GlobalIPRange, &cato_models.GlobalIPRangeRefInput{
				By:    cato_models.ObjectRefBy(itemSourceGlobalIPRangeInput.By.ValueString()),
				Input: itemSourceGlobalIPRangeInput.Input.ValueString(),
			})
		}

		// retrieve & setting source network interface
		elementsSourceNetworkInterfaceInput := make([]types.Object, 0, len(sourceInput.NetworkInterface.Elements()))
		diags = sourceInput.NetworkInterface.ElementsAs(ctx, &elementsSourceNetworkInterfaceInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceNetworkInterfaceInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_NetworkInterface
		for _, item := range elementsSourceNetworkInterfaceInput {
			diags = item.As(ctx, &itemSourceNetworkInterfaceInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.NetworkInterface = append(input.Rule.Source.NetworkInterface, &cato_models.NetworkInterfaceRefInput{
				By:    cato_models.ObjectRefBy(itemSourceNetworkInterfaceInput.By.ValueString()),
				Input: itemSourceNetworkInterfaceInput.Input.ValueString(),
			})
		}

		// retrieve & setting source site network subnet
		elementsSourceSiteNetworkSubnetInput := make([]types.Object, 0, len(sourceInput.SiteNetworkSubnet.Elements()))
		diags = sourceInput.SiteNetworkSubnet.ElementsAs(ctx, &elementsSourceSiteNetworkSubnetInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceSiteNetworkSubnetInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_SiteNetworkSubnet
		for _, item := range elementsSourceSiteNetworkSubnetInput {
			diags = item.As(ctx, &itemSourceSiteNetworkSubnetInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.SiteNetworkSubnet = append(input.Rule.Source.SiteNetworkSubnet, &cato_models.SiteNetworkSubnetRefInput{
				By:    cato_models.ObjectRefBy(itemSourceSiteNetworkSubnetInput.By.ValueString()),
				Input: itemSourceSiteNetworkSubnetInput.Input.ValueString(),
			})
		}

		// retrieve & setting source floating subnet
		elementsSourceFloatingSubnetInput := make([]types.Object, 0, len(sourceInput.FloatingSubnet.Elements()))
		diags = sourceInput.FloatingSubnet.ElementsAs(ctx, &elementsSourceFloatingSubnetInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceFloatingSubnetInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_FloatingSubnet
		for _, item := range elementsSourceFloatingSubnetInput {
			diags = item.As(ctx, &itemSourceFloatingSubnetInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.FloatingSubnet = append(input.Rule.Source.FloatingSubnet, &cato_models.FloatingSubnetRefInput{
				By:    cato_models.ObjectRefBy(itemSourceFloatingSubnetInput.By.ValueString()),
				Input: itemSourceFloatingSubnetInput.Input.ValueString(),
			})
		}

		// retrieve & setting source user
		elementsSourceUserInput := make([]types.Object, 0, len(sourceInput.User.Elements()))
		diags = sourceInput.User.ElementsAs(ctx, &elementsSourceUserInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceUserInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_User
		for _, item := range elementsSourceUserInput {
			diags = item.As(ctx, &itemSourceUserInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.User = append(input.Rule.Source.User, &cato_models.UserRefInput{
				By:    cato_models.ObjectRefBy(itemSourceUserInput.By.ValueString()),
				Input: itemSourceUserInput.Input.ValueString(),
			})
		}

		// retrieve & setting source users group
		elementsSourceUsersGroupInput := make([]types.Object, 0, len(sourceInput.UsersGroup.Elements()))
		diags = sourceInput.UsersGroup.ElementsAs(ctx, &elementsSourceUsersGroupInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceUsersGroupInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_UsersGroup
		for _, item := range elementsSourceUsersGroupInput {
			diags = item.As(ctx, &itemSourceUsersGroupInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.UsersGroup = append(input.Rule.Source.UsersGroup, &cato_models.UsersGroupRefInput{
				By:    cato_models.ObjectRefBy(itemSourceUsersGroupInput.By.ValueString()),
				Input: itemSourceUsersGroupInput.Input.ValueString(),
			})
		}

		// retrieve & setting source group
		elementsSourceGroupInput := make([]types.Object, 0, len(sourceInput.Group.Elements()))
		diags = sourceInput.Group.ElementsAs(ctx, &elementsSourceGroupInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceGroupInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_Group
		for _, item := range elementsSourceGroupInput {
			diags = item.As(ctx, &itemSourceGroupInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.Group = append(input.Rule.Source.Group, &cato_models.GroupRefInput{
				By:    cato_models.ObjectRefBy(itemSourceGroupInput.By.ValueString()),
				Input: itemSourceGroupInput.Input.ValueString(),
			})
		}

		// retrieve & setting source system group
		elementsSourceSystemGroupInput := make([]types.Object, 0, len(sourceInput.SystemGroup.Elements()))
		diags = sourceInput.SystemGroup.ElementsAs(ctx, &elementsSourceSystemGroupInput, false)
		resp.Diagnostics.Append(diags...)

		var itemSourceSystemGroupInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Source_SystemGroup
		for _, item := range elementsSourceSystemGroupInput {
			diags = item.As(ctx, &itemSourceSystemGroupInput, basetypes.ObjectAsOptions{})
			input.Rule.Source.SystemGroup = append(input.Rule.Source.SystemGroup, &cato_models.SystemGroupRefInput{
				By:    cato_models.ObjectRefBy(itemSourceSystemGroupInput.By.ValueString()),
				Input: itemSourceSystemGroupInput.Input.ValueString(),
			})
		}
	}

	// retrieve & setting country
	if !ruleInput.Country.IsNull() {
		elementsCountryInput := make([]types.Object, 0, len(ruleInput.Country.Elements()))
		diags = ruleInput.Country.ElementsAs(ctx, &elementsCountryInput, false)
		resp.Diagnostics.Append(diags...)

		var itemCountryInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Country
		for _, item := range elementsCountryInput {
			diags = item.As(ctx, &itemCountryInput, basetypes.ObjectAsOptions{})
			input.Rule.Country = append(input.Rule.Country, &cato_models.CountryRefInput{
				By:    cato_models.ObjectRefBy(itemCountryInput.By.ValueString()),
				Input: itemCountryInput.Input.ValueString(),
			})
		}
	}

	// retrieve & setting device
	if !ruleInput.Device.IsNull() {
		elementsDeviceInput := make([]types.Object, 0, len(ruleInput.Device.Elements()))
		diags = ruleInput.Device.ElementsAs(ctx, &elementsDeviceInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDeviceInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Device
		for _, item := range elementsDeviceInput {
			diags = item.As(ctx, &itemDeviceInput, basetypes.ObjectAsOptions{})
			input.Rule.Device = append(input.Rule.Device, &cato_models.DeviceProfileRefInput{
				By:    cato_models.ObjectRefBy(itemDeviceInput.By.ValueString()),
				Input: itemDeviceInput.Input.ValueString(),
			})
		}
	}

	// retrieve & setting device OS
	if !ruleInput.Source.IsNull() {
		diags = ruleInput.DeviceOs.ElementsAs(ctx, &input.Rule.DeviceOs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// retrieve & setting destination
	if !ruleInput.Destination.IsNull() {
		destinationInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination{}
		diags = ruleInput.Destination.As(ctx, &destinationInput, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination IP
		diags = destinationInput.IP.ElementsAs(ctx, &input.Rule.Destination.IP, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination subnet
		diags = destinationInput.Subnet.ElementsAs(ctx, &input.Rule.Destination.Subnet, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination domain
		diags = destinationInput.Domain.ElementsAs(ctx, &input.Rule.Destination.Domain, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination fqdn
		diags = destinationInput.Fqdn.ElementsAs(ctx, &input.Rule.Destination.Fqdn, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination remote asn
		diags = destinationInput.RemoteAsn.ElementsAs(ctx, &input.Rule.Destination.RemoteAsn, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting destination application
		elementsDestinationApplicationInput := make([]types.Object, 0, len(destinationInput.Application.Elements()))
		diags = destinationInput.Application.ElementsAs(ctx, &elementsDestinationApplicationInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationApplicationInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_Application
		for _, item := range elementsDestinationApplicationInput {
			diags = item.As(ctx, &itemDestinationApplicationInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.Application = append(input.Rule.Destination.Application, &cato_models.ApplicationRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationApplicationInput.By.ValueString()),
				Input: itemDestinationApplicationInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination custom app
		elementsDestinationCustomAppInput := make([]types.Object, 0, len(destinationInput.CustomApp.Elements()))
		diags = destinationInput.CustomApp.ElementsAs(ctx, &elementsDestinationCustomAppInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationCustomAppInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_CustomApp
		for _, item := range elementsDestinationCustomAppInput {
			diags = item.As(ctx, &itemDestinationCustomAppInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.CustomApp = append(input.Rule.Destination.CustomApp, &cato_models.CustomApplicationRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationCustomAppInput.By.ValueString()),
				Input: itemDestinationCustomAppInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination ip range
		elementsDestinationIPRangeInput := make([]types.Object, 0, len(destinationInput.IPRange.Elements()))
		diags = destinationInput.IPRange.ElementsAs(ctx, &elementsDestinationIPRangeInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationIPRangeInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_IPRange
		for _, item := range elementsDestinationIPRangeInput {
			diags = item.As(ctx, &itemDestinationIPRangeInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.IPRange = append(input.Rule.Destination.IPRange, &cato_models.IPAddressRangeInput{
				From: itemDestinationIPRangeInput.From.ValueString(),
				To:   itemDestinationIPRangeInput.To.ValueString(),
			})
		}

		// retrieve & setting destination global ip range
		elementsDestinationGlobalIPRangeInput := make([]types.Object, 0, len(destinationInput.GlobalIPRange.Elements()))
		diags = destinationInput.GlobalIPRange.ElementsAs(ctx, &elementsDestinationGlobalIPRangeInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationGlobalIPRangeInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_GlobalIPRange
		for _, item := range elementsDestinationGlobalIPRangeInput {
			diags = item.As(ctx, &itemDestinationGlobalIPRangeInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.GlobalIPRange = append(input.Rule.Destination.GlobalIPRange, &cato_models.GlobalIPRangeRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationGlobalIPRangeInput.By.ValueString()),
				Input: itemDestinationGlobalIPRangeInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination app category
		elementsDestinationAppCategoryInput := make([]types.Object, 0, len(destinationInput.AppCategory.Elements()))
		diags = destinationInput.AppCategory.ElementsAs(ctx, &elementsDestinationAppCategoryInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationAppCategoryInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_AppCategory
		for _, item := range elementsDestinationAppCategoryInput {
			diags = item.As(ctx, &itemDestinationAppCategoryInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.AppCategory = append(input.Rule.Destination.AppCategory, &cato_models.ApplicationCategoryRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationAppCategoryInput.By.ValueString()),
				Input: itemDestinationAppCategoryInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination custom app category
		elementsDestinationCustomCategoryInput := make([]types.Object, 0, len(destinationInput.CustomCategory.Elements()))
		diags = destinationInput.CustomCategory.ElementsAs(ctx, &elementsDestinationCustomCategoryInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationCustomCategoryInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_CustomCategory
		for _, item := range elementsDestinationCustomCategoryInput {
			diags = item.As(ctx, &itemDestinationCustomCategoryInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.CustomCategory = append(input.Rule.Destination.CustomCategory, &cato_models.CustomCategoryRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationCustomCategoryInput.By.ValueString()),
				Input: itemDestinationCustomCategoryInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination sanctionned apps category
		elementsDestinationSanctionedAppsCategoryInput := make([]types.Object, 0, len(destinationInput.SanctionedAppsCategory.Elements()))
		diags = destinationInput.SanctionedAppsCategory.ElementsAs(ctx, &elementsDestinationSanctionedAppsCategoryInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationSanctionedAppsCategoryInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_SanctionedAppsCategory
		for _, item := range elementsDestinationSanctionedAppsCategoryInput {
			diags = item.As(ctx, &itemDestinationSanctionedAppsCategoryInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.SanctionedAppsCategory = append(input.Rule.Destination.SanctionedAppsCategory, &cato_models.SanctionedAppsCategoryRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationSanctionedAppsCategoryInput.By.ValueString()),
				Input: itemDestinationSanctionedAppsCategoryInput.Input.ValueString(),
			})
		}

		// retrieve & setting destination country
		elementsDestinationCountryInput := make([]types.Object, 0, len(destinationInput.Country.Elements()))
		diags = destinationInput.Country.ElementsAs(ctx, &elementsDestinationCountryInput, false)
		resp.Diagnostics.Append(diags...)

		var itemDestinationCountryInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Destination_Country
		for _, item := range elementsDestinationCountryInput {
			diags = item.As(ctx, &itemDestinationCountryInput, basetypes.ObjectAsOptions{})
			input.Rule.Destination.Country = append(input.Rule.Destination.Country, &cato_models.CountryRefInput{
				By:    cato_models.ObjectRefBy(itemDestinationCountryInput.By.ValueString()),
				Input: itemDestinationCountryInput.Input.ValueString(),
			})
		}
	}

	// retrieve & setting service
	if !ruleInput.Service.IsNull() {
		serviceInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service{}
		diags = ruleInput.Service.As(ctx, &serviceInput, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting service standard
		if !serviceInput.Standard.IsNull() {
			elementsServiceStandardInput := make([]types.Object, 0, len(serviceInput.Standard.Elements()))
			diags = serviceInput.Standard.ElementsAs(ctx, &elementsServiceStandardInput, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			var itemServiceStandardInput Policy_Policy_InternetFirewall_Policy_Rules_Rule_Service_Standard
			for _, item := range elementsServiceStandardInput {
				diags = item.As(ctx, &itemServiceStandardInput, basetypes.ObjectAsOptions{})
				input.Rule.Service.Standard = append(input.Rule.Service.Standard, &cato_models.ServiceRefInput{
					By:    cato_models.ObjectRefBy(itemServiceStandardInput.By.ValueString()),
					Input: itemServiceStandardInput.Input.ValueString(),
				})
			}
		}
	}

	// retrieve & setting tracking

	if !ruleInput.Tracking.IsNull() {
		trackingInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking{}
		diags = ruleInput.Tracking.As(ctx, &trackingInput, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// retrieve & setting tracking event
		trackingEventInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Event{}
		diags = trackingInput.Event.As(ctx, &trackingEventInput, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		input.Rule.Tracking.Event.Enabled = trackingEventInput.Enabled.ValueBoolPointer()

		// retrieve & setting tracking Alert
		if !trackingInput.Alert.IsNull() {
			input.Rule.Tracking.Alert = &cato_models.PolicyRuleTrackingAlertUpdateInput{}

			trackingAlertInput := Policy_Policy_InternetFirewall_Policy_Rules_Rule_Tracking_Alert{}
			diags = trackingInput.Alert.As(ctx, &trackingAlertInput, basetypes.ObjectAsOptions{})
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			input.Rule.Tracking.Alert.Enabled = trackingAlertInput.Enabled.ValueBoolPointer()
			input.Rule.Tracking.Alert.Frequency = (*cato_models.PolicyRuleTrackingFrequencyEnum)(trackingAlertInput.Frequency.ValueStringPointer())

			// retrieve & setting tracking alert subscription group
			diags = trackingAlertInput.SubscriptionGroup.ElementsAs(ctx, &input.Rule.Tracking.Alert.SubscriptionGroup, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			// retrieve & setting tracking alert webhook
			diags = trackingAlertInput.Webhook.ElementsAs(ctx, &input.Rule.Tracking.Alert.Webhook, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}

			// retrieve & setting tracking alert mailing list
			diags = trackingAlertInput.MailingList.ElementsAs(ctx, &input.Rule.Tracking.Alert.MailingList, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	// settings other rule attributes
	input.ID = *ruleInput.ID.ValueStringPointer()
	input.Rule.Name = ruleInput.Name.ValueStringPointer()
	input.Rule.Description = ruleInput.Description.ValueStringPointer()
	input.Rule.Enabled = ruleInput.Enabled.ValueBoolPointer()
	input.Rule.ConnectionOrigin = (*cato_models.ConnectionOriginEnum)(ruleInput.ConnectionOrigin.ValueStringPointer())
	input.Rule.Action = (*cato_models.InternetFirewallActionEnum)(ruleInput.Action.ValueStringPointer())

	mutationInput := &cato_models.InternetFirewallPolicyMutationInput{}

	b, _ := json.Marshal(input)

	tflog.Info(ctx, "update input")
	tflog.Info(ctx, string(b))

	//creating new rule
	_, err := r.info.catov2.PolicyInternetFirewallUpdateRule(ctx, mutationInput, input, r.info.AccountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Catov2 API PolicyInternetFirewallUpdateRule error",
			err.Error(),
		)
		return
	}

	//publishing new rule
	tflog.Info(ctx, "publishing new rule")
	publishDataIfEnabled := &cato_models.PolicyPublishRevisionInput{}
	_, err = r.info.catov2.PolicyInternetFirewallPublishPolicyRevision(ctx, &cato_models.InternetFirewallPolicyMutationInput{}, publishDataIfEnabled, r.info.AccountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Catov2 API PolicyInternetFirewallPublishPolicyRevision error",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *internetFwPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state InternetFirewallCreateRuleInput
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//retrieve rule ID
	rule := Policy_Policy_InternetFirewall_Policy_Rules_Rule{}
	diags = state.Rule.As(ctx, &rule, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	removeMutations := &cato_models.InternetFirewallPolicyMutationInput{}
	removeRule := cato_models.InternetFirewallRemoveRuleInput{
		ID: rule.ID.ValueString(),
	}

	_, err := r.info.catov2.PolicyInternetFirewallRemoveRule(ctx, removeMutations, removeRule, r.info.AccountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to connect or request the Catov2 API",
			err.Error(),
		)
		return
	}

	publishDataIfEnabled := &cato_models.PolicyPublishRevisionInput{}
	_, err = r.info.catov2.PolicyInternetFirewallPublishPolicyRevision(ctx, &cato_models.InternetFirewallPolicyMutationInput{}, publishDataIfEnabled, r.info.AccountId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Catov2 API Delete/PolicyInternetFirewallPublishPolicyRevision error",
			err.Error(),
		)
		return
	}

}

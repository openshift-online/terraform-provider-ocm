package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift-online/ocm-sdk-go/logging"
)

type ClusterRosaClassicDataSourceType struct {
	logger logging.Logger
}

type ClusterRosaClassicDataSource struct {
	logger     logging.Logger
	collection *cmv1.ClustersClient
}

func (t *ClusterRosaClassicDataSourceType) GetSchema(ctx context.Context) (result tfsdk.Schema,
	diags diag.Diagnostics) {
	result = tfsdk.Schema{
		Description: "OpenShift managed cluster using rosa sts.",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description: "Name of the cluster.",
				Type:        types.StringType,
				Computed:    true,
			},
			//future : search with id or name instead of just name
			"id": {
				Description: "Unique identifier of the cluster.",
				Type:        types.StringType,
				Computed:    true,
			},
			"version": {
				Description: "Identifier of the version of OpenShift, for example 'openshift-v4.1.0'.",
				Type:        types.StringType,
				Computed:    true,
			},
			"state": {
				Description: "State of the cluster.",
				Type:        types.StringType,
				Computed:    true,
			},
			// to include -> STS Role ARN, Support Role ARN, Instance IAM Roles - Control plane, Instance IAM Roles - Worker, Operator IAM Roles (list), Created,Details Page,OIDC Endpoint URL
		},
	}
	return
}

func (t *ClusterRosaClassicDataSourceType) NewDataSource(ctx context.Context,
	p tfsdk.Provider) (result tfsdk.DataSource, diags diag.Diagnostics) {
	// Cast the provider interface to the specific implementation:
	parent := p.(*Provider)

	// Get the collection of clusters:
	collection := parent.connection.ClustersMgmt().V1().Clusters()

	// Create the resource:
	result = &ClusterRosaClassicDataSource{
		logger:     parent.logger,
		collection: collection,
	}
	return
}

func (s *ClusterRosaClassicDataSource) Read(ctx context.Context, request tfsdk.ReadDataSourceRequest,
	response *tfsdk.ReadDataSourceResponse) {
	// Get the state:
	state := &DataClusterRosaClassicState{}
	diags := request.Config.Get(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	// Fetch the complete details of the cluster:
	//var listItem *cmv1.Cluster
	listRequest := s.collection.List()
	if !state.Name.Unknown && !state.Name.Null {
		listRequest.Search("name is '" + state.Name.Value + "'")
	} else if !state.ID.Unknown && !state.ID.Null {
		listRequest.Search("id is '" + state.ID.Value + "'")
	} else {
		response.Diagnostics.AddError(
			"Data source requires either cluster Name or ID",
			"",
		)
		return
	}
	listResponse, err := listRequest.SendContext(ctx)
	if err != nil {
		response.Diagnostics.AddError(
			"Can't list Clusters",
			err.Error(),
		)
		return
	}
	if listResponse.Size() > 0 {
		// Populate the state:
		listResponse.Items().Each(func(cluster *cmv1.Cluster) bool {
			state.ID = types.String{
				Value: cluster.ID(),
			}
			state.Name = types.String{
				Value: cluster.Name(),
			}
			state.State = types.String{
				Value: fmt.Sprintf("%s", cluster.State()),
			}
			state.Version = types.String{
				Value: cluster.OpenshiftVersion(),
			}
			return true
		})
	}
	// Save the state:
	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
}

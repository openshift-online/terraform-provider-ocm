/*
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)
type ClusterResourceUtilsImpl struct {
	clusterClient ClusterClient
}

func (t *ClusterResourceType) GetSchema(ctx context.Context) (result tfsdk.Schema,
	diags diag.Diagnostics) {
	result = tfsdk.Schema{
		Description: "OpenShift managed cluster.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description: "Unique identifier of the cluster.",
				Type:        types.StringType,
				Computed:    true,
			},
			"product": {
				Description: "Product ID OSD or Rosa",
				Type:        types.StringType,
				Required:    true,
			},
			"name": {
				Description: "Name of the cluster.",
				Type:        types.StringType,
				Required:    true,
			},
			"cloud_provider": {
				Description: "Cloud provider identifier, for example 'aws'.",
				Type:        types.StringType,
				Required:    true,
			},
			"cloud_region": {
				Description: "Cloud region identifier, for example 'us-east-1'.",
				Type:        types.StringType,
				Required:    true,
			},
			"multi_az": {
				Description: "Indicates if the cluster should be deployed to " +
					"multiple availability zones. Default value is 'false'.",
				Type:     types.BoolType,
				Optional: true,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"properties": {
				Description: "User defined properties.",
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"api_url": {
				Description: "URL of the API server.",
				Type:        types.StringType,
				Computed:    true,
			},
			"console_url": {
				Description: "URL of the console.",
				Type:        types.StringType,
				Computed:    true,
			},
			"compute_nodes": {
				Description: "Number of compute nodes of the cluster.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"compute_machine_type": {
				Description: "Identifier of the machine type used by the compute nodes, " +
					"for example `r5.xlarge`. Use the `ocm_machine_types` data " +
					"source to find the possible values.",
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"ccs_enabled": {
				Description: "Enables customer cloud subscription.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"aws_account_id": {
				Description: "Identifier of the AWS account.",
				Type:        types.StringType,
				Optional:    true,
			},
			"aws_access_key_id": {
				Description: "Identifier of the AWS access key.",
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
			},
			"aws_secret_access_key": {
				Description: "AWS access key.",
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
			},
			"aws_subnet_ids": {
				Description: "aws subnet ids",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"aws_private_link": {
				Description: "aws subnet ids",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"availability_zones": {
				Description: "availability zones",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Optional: true,
			},
			"machine_cidr": {
				Description: "Block of IP addresses for nodes.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"proxy": {
				Description: "proxy",
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"http_proxy": {
						Description: "http proxy",
						Type:        types.StringType,
						Required:    true,
					},
					"https_proxy": {
						Description: "https proxy",
						Type:        types.StringType,
						Required:    true,
					},
					"no_proxy": {
						Description: "no proxy",
						Type:        types.StringType,
						Optional:    true,
					},
				}),
				Optional: true,
			},
			"service_cidr": {
				Description: "Block of IP addresses for services.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"pod_cidr": {
				Description: "Block of IP addresses for pods.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"host_prefix": {
				Description: "Length of the prefix of the subnet assigned to each node.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"version": {
				Description: "Identifier of the version of OpenShift, for example 'openshift-v4.1.0'.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"state": {
				Description: "State of the cluster.",
				Type:        types.StringType,
				Computed:    true,
			},
			"wait": {
				Description: "Wait till the cluster is ready.",
				Type:        types.BoolType,
				Optional:    true,
			},
		},
	}
	return
}

func (t *ClusterResourceType) NewResource(ctx context.Context,
	p tfsdk.Provider) (result tfsdk.Resource, diags diag.Diagnostics) {
	// Cast the provider interface to the specific implementation:
	parent := p.(*Provider)

	// Get the collection of clusters:
	collection := parent.connection.ClustersMgmt().V1().Clusters()

	// Create the resource:
	result = &ClusterResource{
		logger:     parent.logger,
		collection: collection,
	}

	return
}

func (clusterUtils ClusterResourceUtilsImpl) createClusterObject(ctx context.Context,
	state *ClusterState, diags Diagnostics) (*cmv1.Cluster, error) {
	// Create the cluster:
	builder := cmv1.NewCluster()
	builder.Name(state.Name.Value)
	builder.CloudProvider(cmv1.NewCloudProvider().ID(state.CloudProvider.Value))
	builder.Product(cmv1.NewProduct().ID(state.Product.Value))
	builder.Region(cmv1.NewCloudRegion().ID(state.CloudRegion.Value))
	if !state.MultiAZ.Unknown && !state.MultiAZ.Null {
		builder.MultiAZ(state.MultiAZ.Value)
	}
	if !state.Properties.Unknown && !state.Properties.Null {
		properties := map[string]string{}
		for k, v := range state.Properties.Elems {
			properties[k] = v.(types.String).Value
		}
		builder.Properties(properties)
	}
	nodes := cmv1.NewClusterNodes()
	if !state.ComputeNodes.Unknown && !state.ComputeNodes.Null {
		nodes.Compute(int(state.ComputeNodes.Value))
	}
	if !state.ComputeMachineType.Unknown && !state.ComputeMachineType.Null {
		nodes.ComputeMachineType(
			cmv1.NewMachineType().ID(state.ComputeMachineType.Value),
		)
	}

	if !state.AvailabilityZones.Unknown && !state.AvailabilityZones.Null {
		azs := make([]string, 0)
		for _, e := range state.AvailabilityZones.Elems {
			azs = append(azs, e.(types.String).Value)
		}
		nodes.AvailabilityZones(azs...)
	}

	if !nodes.Empty() {
		builder.Nodes(nodes)
	}
	ccs := cmv1.NewCCS()
	if !state.CCSEnabled.Unknown && !state.CCSEnabled.Null {
		ccs.Enabled(state.CCSEnabled.Value)
	}
	if !ccs.Empty() {
		builder.CCS(ccs)
	}
	aws := cmv1.NewAWS()
	if !state.AWSAccountID.Unknown && !state.AWSAccountID.Null {
		aws.AccountID(state.AWSAccountID.Value)
	}
	if !state.AWSAccessKeyID.Unknown && !state.AWSAccessKeyID.Null {
		aws.AccessKeyID(state.AWSAccessKeyID.Value)
	}
	if !state.AWSSecretAccessKey.Unknown && !state.AWSSecretAccessKey.Null {
		aws.SecretAccessKey(state.AWSSecretAccessKey.Value)
	}
	if !state.AWSPrivateLink.Unknown && !state.AWSPrivateLink.Null {
		aws.PrivateLink((state.AWSPrivateLink.Value))
		api := cmv1.NewClusterAPI()
		if state.AWSPrivateLink.Value {
			api.Listening(cmv1.ListeningMethodInternal)
		}
		builder.API(api)
	}

	if !state.AWSSubnetIDs.Unknown && !state.AWSSubnetIDs.Null {
		subnetIds := make([]string, 0)
		for _, e := range state.AWSSubnetIDs.Elems {
			subnetIds = append(subnetIds, e.(types.String).Value)
		}
		aws.SubnetIDs(subnetIds...)
	}

	if !aws.Empty() {
		builder.AWS(aws)
	}
	network := cmv1.NewNetwork()
	if !state.MachineCIDR.Unknown && !state.MachineCIDR.Null {
		network.MachineCIDR(state.MachineCIDR.Value)
	}
	if !state.ServiceCIDR.Unknown && !state.ServiceCIDR.Null {
		network.ServiceCIDR(state.ServiceCIDR.Value)
	}
	if !state.PodCIDR.Unknown && !state.PodCIDR.Null {
		network.PodCIDR(state.PodCIDR.Value)
	}
	if !state.HostPrefix.Unknown && !state.HostPrefix.Null {
		network.HostPrefix(int(state.HostPrefix.Value))
	}
	if !network.Empty() {
		builder.Network(network)
	}
	if !state.Version.Unknown && !state.Version.Null {
		builder.Version(cmv1.NewVersion().ID(state.Version.Value))
	}

	proxy := cmv1.NewProxy()
	if state.Proxy != nil {
		proxy.HTTPProxy(state.Proxy.HttpProxy.Value)
		proxy.HTTPSProxy(state.Proxy.HttpsProxy.Value)
		builder.Proxy(proxy)
	}

	object, err := builder.Build()

	return object, err
}

type Diagnostics interface {
	AddError(title, description string)
}

func (clusterUtils ClusterResourceUtilsImpl) create(ctx context.Context,
	state *ClusterState, diags Diagnostics) error {

	object, err := clusterUtils.createClusterObject(ctx, state, diags)
	if err != nil {
		diags.AddError(
			"Can't build cluster",
			fmt.Sprintf(
				"Can't build cluster with name '%s': %v",
				state.Name.Value, err,
			),
		)
		return err
	}

	object, err = clusterUtils.clusterClient.Create(ctx, object)
	if err != nil {
		diags.AddError(
			clusterCreationFailureMessage,
			fmt.Sprintf(
				"Can't create cluster with name '%s': %v",
				state.Name.Value, err,
			),
		)
		return err
	}

	// Wait till the cluster is ready unless explicitly disabled:
	wait := state.Wait.Unknown || state.Wait.Null || state.Wait.Value
	ready := object.State() == cmv1.ClusterStateReady
	if wait && !ready {
		err := clusterUtils.clusterClient.PollReady(ctx, object.ID())
		if err != nil {
			diags.AddError(
				clusterPollFailure,
				fmt.Sprintf(
					"Can't poll state of cluster with identifier '%s': %v",
					object.ID(), err,
				),
			)
			return err
		}
	}

	// Save the state:
	clusterUtils.populateClusterState(object, state)

	return nil
}

func (r *ClusterResource) Create(ctx context.Context,
	request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	// Get the plan:
	state := &ClusterState{}
	diags := request.Plan.Get(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	clusterClient := &ClusterClientImpl{InternalClient: r.collection}

	clusterUtils := ClusterResourceUtilsImpl{clusterClient: clusterClient}

	err := clusterUtils.create(ctx, state, &diags)
	if err != nil {
		return
	}

	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
}

func (r *ClusterResource) Read(ctx context.Context, request tfsdk.ReadResourceRequest,
	response *tfsdk.ReadResourceResponse) {
	// Get the current state:
	state := &ClusterState{}
	diags := request.State.Get(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	clusterClient := &ClusterClientImpl{InternalClient: r.collection}
	clusterUtils := ClusterResourceUtilsImpl{clusterClient: clusterClient}

	// Find the cluster:
	get, err := r.collection.Cluster(state.ID.Value).Get().SendContext(ctx)
	if err != nil {
		response.Diagnostics.AddError(
			"Can't find cluster",
			fmt.Sprintf(
				"Can't find cluster with identifier '%s': %v",
				state.ID.Value, err,
			),
		)
		return
	}
	object := get.Body()

	// Save the state:
	clusterUtils.populateClusterState(object, state)
	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
}

func (clusterUtils ClusterResourceUtilsImpl) update(ctx context.Context,
	currentState *ClusterState, desiredState *ClusterState, diags Diagnostics) error {

	builder := cmv1.NewCluster()
	nodes := cmv1.NewClusterNodes()
	compute, ok := shouldPatchInt(currentState.ComputeNodes, desiredState.ComputeNodes)
	if ok {
		nodes.Compute(int(compute))
	}
	if !nodes.Empty() {
		builder.Nodes(nodes)
	}
	patch, err := builder.Build()
	if err != nil {
		diags.AddError(
			"Can't build cluster patch",
			fmt.Sprintf(
				"Can't build patch for cluster with identifier '%s': %v",
				currentState.ID.Value, err,
			),
		)
		return err
	}
	object, err := clusterUtils.clusterClient.Update(ctx, currentState.ID.Value, patch)
	if err != nil {
		diags.AddError(
			clusterUpdateFailureMessage,
			fmt.Sprintf(
				"Can't update cluster with identifier '%s': %v",
				currentState.ID.Value, err,
			),
		)
		return err
	}

	// Update the state:
	clusterUtils.populateClusterState(object, currentState)

	return nil
}

func (r *ClusterResource) Update(ctx context.Context, request tfsdk.UpdateResourceRequest,
	response *tfsdk.UpdateResourceResponse) {
	var diags diag.Diagnostics

	// Get the state:
	state := &ClusterState{}
	diags = request.State.Get(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// Get the plan:
	plan := &ClusterState{}
	diags = request.Plan.Get(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	clusterClient := &ClusterClientImpl{InternalClient: r.collection}

	clusterUtils := ClusterResourceUtilsImpl{clusterClient: clusterClient}
	err := clusterUtils.update(ctx, state, plan, &diags)
	if err != nil {
		return
	}

	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
}

func (clusterUtils ClusterResourceUtilsImpl) delete(ctx context.Context,
	id string, shouldPoll bool, diags Diagnostics) error {
	err := clusterUtils.clusterClient.Delete(ctx, id)
	if err != nil {
		diags.AddError(
			clusterDeleteFailureMessage,
			fmt.Sprintf(
				"Can't delete cluster with identifier '%s': %v",
				id, err,
			),
		)
		return err
	}

	// Wait till the cluster has been effectively deleted:
	if shouldPoll {
		err := clusterUtils.clusterClient.PollRemoved(ctx, id)
		if err != nil {
			diags.AddError(
				clusterPollDeletionFailure,
				fmt.Sprintf(
					"Can't poll deletion of cluster with identifier '%s': %v",
					id, err,
				),
			)
			return err
		}
	}

	return nil
}

func (r *ClusterResource) Delete(ctx context.Context, request tfsdk.DeleteResourceRequest,
	response *tfsdk.DeleteResourceResponse) {
	// Get the state:
	state := &ClusterState{}
	diags := request.State.Get(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	clusterClient := &ClusterClientImpl{InternalClient: r.collection}

	clusterUtils := ClusterResourceUtilsImpl{clusterClient: clusterClient}
	shouldPoll := state.Wait.Unknown || state.Wait.Null || state.Wait.Value
	err := clusterUtils.delete(ctx, state.ID.Value, shouldPoll, &diags)
	if err != nil {
		return
	}

	// Remove the state:
	response.State.RemoveResource(ctx)
}

func (r *ClusterResource) ImportState(ctx context.Context, request tfsdk.ImportResourceStateRequest,
	response *tfsdk.ImportResourceStateResponse) {
	// Try to retrieve the object:
	get, err := r.collection.Cluster(request.ID).Get().SendContext(ctx)
	if err != nil {
		response.Diagnostics.AddError(
			"Can't find cluster",
			fmt.Sprintf(
				"Can't find cluster with identifier '%s': %v",
				request.ID, err,
			),
		)
		return
	}
	object := get.Body()

	clusterUtils := ClusterResourceUtilsImpl{}

	// Save the state:
	state := &ClusterState{}
	clusterUtils.populateClusterState(object, state)
	diags := response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
}

// populateClusterState copies the data from the API object to the Terraform state.
func (clusterUtils ClusterResourceUtilsImpl) populateClusterState(object *cmv1.Cluster, state *ClusterState) {
	state.ID = types.String{
		Value: object.ID(),
	}

	object.API()
	state.Product = types.String{
		Value: object.Product().ID(),
	}
	state.Name = types.String{
		Value: object.Name(),
	}
	state.CloudProvider = types.String{
		Value: object.CloudProvider().ID(),
	}
	state.CloudRegion = types.String{
		Value: object.Region().ID(),
	}
	state.MultiAZ = types.Bool{
		Value: object.MultiAZ(),
	}
	state.Properties = types.Map{
		ElemType: types.StringType,
		Elems:    map[string]attr.Value{},
	}
	for k, v := range object.Properties() {
		state.Properties.Elems[k] = types.String{
			Value: v,
		}
	}
	state.APIURL = types.String{
		Value: object.API().URL(),
	}
	state.ConsoleURL = types.String{
		Value: object.Console().URL(),
	}
	state.ComputeNodes = types.Int64{
		Value: int64(object.Nodes().Compute()),
	}
	state.ComputeMachineType = types.String{
		Value: object.Nodes().ComputeMachineType().ID(),
	}

	azs, ok := object.Nodes().GetAvailabilityZones()
	if ok {
		state.AvailabilityZones.Elems = make([]attr.Value, 0)
		for _, az := range azs {
			state.AvailabilityZones.Elems = append(state.AvailabilityZones.Elems, types.String{
				Value: az,
			})
		}
	}

	state.CCSEnabled = types.Bool{
		Value: object.CCS().Enabled(),
	}
	//The API does not return account id
	awsAccountID, ok := object.AWS().GetAccountID()
	if ok {
		state.AWSAccountID = types.String{
			Value: awsAccountID,
		}
	}
	awsAccessKeyID, ok := object.AWS().GetAccessKeyID()
	if ok {
		state.AWSAccessKeyID = types.String{
			Value: awsAccessKeyID,
		}
	} else {
		state.AWSAccessKeyID = types.String{
			Null: true,
		}
	}
	awsSecretAccessKey, ok := object.AWS().GetSecretAccessKey()
	if ok {
		state.AWSSecretAccessKey = types.String{
			Value: awsSecretAccessKey,
		}
	} else {
		state.AWSSecretAccessKey = types.String{
			Null: true,
		}
	}
	awsPrivateLink, ok := object.AWS().GetPrivateLink()
	if ok {
		state.AWSPrivateLink = types.Bool{
			Value: awsPrivateLink,
		}
	} else {
		state.AWSPrivateLink = types.Bool{
			Null: true,
		}
	}

	subnetIds, ok := object.AWS().GetSubnetIDs()
	if ok {
		state.AWSSubnetIDs.Elems = make([]attr.Value, 0)
		for _, subnetId := range subnetIds {
			state.AWSSubnetIDs.Elems = append(state.AWSSubnetIDs.Elems, types.String{
				Value: subnetId,
			})
		}
	}

	proxy, ok := object.GetProxy()
	if ok {
		state.Proxy.HttpProxy = types.String{
			Value: proxy.HTTPProxy(),
		}
		state.Proxy.HttpsProxy = types.String{
			Value: proxy.HTTPSProxy(),
		}
	}
	machineCIDR, ok := object.Network().GetMachineCIDR()
	if ok {
		state.MachineCIDR = types.String{
			Value: machineCIDR,
		}
	} else {
		state.MachineCIDR = types.String{
			Null: true,
		}
	}
	serviceCIDR, ok := object.Network().GetServiceCIDR()
	if ok {
		state.ServiceCIDR = types.String{
			Value: serviceCIDR,
		}
	} else {
		state.ServiceCIDR = types.String{
			Null: true,
		}
	}
	podCIDR, ok := object.Network().GetPodCIDR()
	if ok {
		state.PodCIDR = types.String{
			Value: podCIDR,
		}
	} else {
		state.PodCIDR = types.String{
			Null: true,
		}
	}
	hostPrefix, ok := object.Network().GetHostPrefix()
	if ok {
		state.HostPrefix = types.Int64{
			Value: int64(hostPrefix),
		}
	} else {
		state.HostPrefix = types.Int64{
			Null: true,
		}
	}
	version, ok := object.Version().GetID()
	if ok {
		state.Version = types.String{
			Value: version,
		}
	} else {
		state.Version = types.String{
			Null: true,
		}
	}
	state.State = types.String{
		Value: string(object.State()),
	}

}

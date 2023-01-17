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
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/logging"
)

const (
	clusterCreationFailureMessage = "Can't create cluster"
	clusterPollFailure = "Can't poll cluster state"
	clusterUpdateFailureMessage = "Can't update cluster"
	clusterDeleteFailureMessage = "Can't delete cluster"
	clusterPollDeletionFailure = "Can't poll cluster deletion"
)

type ClusterResourceType struct {
}

type ClusterResource struct {
	logger     logging.Logger
	collection *cmv1.ClustersClient
}

type ClusterClient interface {
	Get(ctx context.Context, id string) (*cmv1.Cluster, error)
	Create(ctx context.Context, cluster *cmv1.Cluster) (*cmv1.Cluster, error)
	PollReady(ctx context.Context, id string) error
	PollRemoved(ctx context.Context, id string) error
	Update(ctx context.Context, id string, cluster *cmv1.Cluster) (*cmv1.Cluster, error)
	Delete(ctx context.Context, id string) error
}

type ClusterClientImpl struct {
	InternalClient *cmv1.ClustersClient
}

func (clusterClient *ClusterClientImpl) Get(ctx context.Context, id string) (*cmv1.Cluster, error) {
	response, err := clusterClient.InternalClient.Cluster(id).Get().SendContext(ctx)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func (clusterClient *ClusterClientImpl) Create(ctx context.Context, cluster *cmv1.Cluster) (*cmv1.Cluster, error) {
	response, err := clusterClient.InternalClient.Add().Body(cluster).SendContext(ctx)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func (clusterClient *ClusterClientImpl) PollReady(ctx context.Context, id string) error {
	pollCtx, cancel := context.WithTimeout(ctx, 1*time.Hour)
	defer cancel()
	_, err := clusterClient.InternalClient.Cluster(id).Poll().
		Interval(30 * time.Second).
		Predicate(func(get *cmv1.ClusterGetResponse) bool {
			object := get.Body()
			return object.State() == cmv1.ClusterStateReady
		}).
		StartContext(pollCtx)
	if err != nil {
		return err
	}

	return nil
}

func (clusterClient *ClusterClientImpl) PollRemoved(ctx context.Context, id string) error {
	pollCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	_, err := clusterClient.InternalClient.Cluster(id).Poll().
		Interval(30 * time.Second).
		Status(http.StatusNotFound).
		StartContext(pollCtx)
	sdkErr, ok := err.(*errors.Error)
	if ok && sdkErr.Status() == http.StatusNotFound {
		err = nil
	}
	if err != nil {
		return err
	}

	return nil
}

func (clusterClient *ClusterClientImpl) Update(ctx context.Context, id string, patch *cmv1.Cluster) (*cmv1.Cluster, error) {
	update, err := clusterClient.InternalClient.Cluster(id).Update().
		Body(patch).
		SendContext(ctx)
	if err != nil {
		return nil, err
	}
	return update.Body(), nil
}

func (clusterClient *ClusterClientImpl) Delete(ctx context.Context, id string) error {
	_, err := clusterClient.InternalClient.Cluster(id).Delete().SendContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

type ClusterResourceUtils interface {
	populateClusterState(object *cmv1.Cluster, state *ClusterState)
	createClusterObject(ctx context.Context,
		state *ClusterState, diags diag.Diagnostics) (*cmv1.Cluster, error)
}
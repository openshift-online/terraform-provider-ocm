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
	"encoding/json"
	"fmt"

	gomock "github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/onsi/ginkgo/v2/dsl/core" // nolint
	. "github.com/onsi/gomega"             // nolint
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

type GinkgoTestReporter struct{}

func (g GinkgoTestReporter) Errorf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

func (g GinkgoTestReporter) Fatalf(format string, args ...interface{}) {
	Fail(fmt.Sprintf(format, args...))
}

type Matcher struct {
	comparator func(x interface{}) bool
	text       string
}

func (matcher Matcher) Matches(x interface{}) bool {
	return matcher.comparator(x)
}

func (matcher Matcher) String() string {
	return matcher.text
}

func generateBaseClusterMap() map[string]interface{} {
	return map[string]interface{}{
		"id": clusterId,
		"product": map[string]interface{}{
			"id": productId,
		},
		"cloud_provider": map[string]interface{}{
			"id": cloudProviderId,
		},
		"region": map[string]interface{}{
			"id": regionId,
		},
		"multi_az": multiAz,
		"properties": map[string]interface{}{
			"rosa_creator_arn": rosaCreatorArn,
		},
		"api": map[string]interface{}{
			"url": apiUrl,
		},
		"console": map[string]interface{}{
			"url": consoleUrl,
		},
		"nodes": map[string]interface{}{
			"compute_machine_type": map[string]interface{}{
				"id": machineType,
			},
			"availability_zones": []interface{}{
				availabilityZone1,
			},
			"compute": computeNodes,
		},
		"ccs": map[string]interface{}{
			"enabled": ccsEnabled,
		},
		"aws": map[string]interface{}{
			"account_id":        awsAccountID,
			"access_key_id":     awsAccessKeyID,
			"secret_access_key": awsSecretAccessKey,
			"private_link":      privateLink,
		},
		"version": map[string]interface{}{
			"id": clusterVersion,
		},
	}
}

var _ = Describe("Cluster creation", func() {
	

	parseClusterMapToObject := func(clusterMap map[string]interface{}) *cmv1.Cluster {
		clusterJsonString, err := json.Marshal(clusterMap)
		Expect(err).To(BeNil())

		clusterObject, err := cmv1.UnmarshalCluster(clusterJsonString)
		Expect(err).To(BeNil())

		return clusterObject
	}

	generateBaseClusterObject := func() *cmv1.Cluster {
		clusterJson := generateBaseClusterMap()
		clusterObject := parseClusterMapToObject(clusterJson)

		return clusterObject
	}

	generateBaseClusterState := func() *ClusterState {
		clusterState := &ClusterState{
			ID: types.String{
				Value: clusterId,
			},
			Name: types.String{
				Value: clusterName,
			},
			ComputeNodes: types.Int64{
				Value: int64(computeNodes),
			},
			Version: types.String{
				Value: clusterVersion,
			},
			CloudRegion: types.String{
				Value: regionId,
			},
			AWSAccountID: types.String{
				Value: awsAccountID,
			},
			AvailabilityZones: types.List{
				Elems: []attr.Value{
					types.String{
						Value: availabilityZone1,
					},
				},
			},
			Properties: types.Map{
				Elems: map[string]attr.Value{
					"rosa_creator_arn": types.String{
						Value: rosaCreatorArn,
					},
				},
			},
			Wait: types.Bool{
				Value: false,
			},
		}

		return clusterState
	}

	generateGinkgoController := func() *gomock.Controller {
		reporter := GinkgoTestReporter{}
		controller := gomock.NewController(reporter)
		return controller
	}

	It("Creates ClusterBuilder with correct field values", func() {
		ctrl := generateGinkgoController()
		clusterState := generateBaseClusterState()

		clusterUtils := ClusterResourceUtilsImpl{}
		clusterObject, err := clusterUtils.createClusterObject(context.Background(), clusterState, NewMockDiagnostics(ctrl))
		Expect(err).To(BeNil())

		Expect(err).To(BeNil())
		Expect(clusterObject).ToNot(BeNil())

		Expect(clusterObject.Name()).To(Equal(clusterName))
		Expect(clusterObject.Version().ID()).To(Equal(clusterVersion))

		id, ok := clusterObject.Region().GetID()
		Expect(ok).To(BeTrue())
		Expect(id).To(Equal(regionId))

		Expect(clusterObject.AWS().AccountID()).To(Equal(awsAccountID))

		availabilityZones := clusterObject.Nodes().AvailabilityZones()
		Expect(availabilityZones).To(HaveLen(1))
		Expect(availabilityZones[0]).To(Equal(availabilityZone1))

		arn, ok := clusterObject.Properties()["rosa_creator_arn"]
		Expect(ok).To(BeTrue())
		Expect(arn).To(Equal(arn))
	})

	It("populateClusterState converts correctly a Cluster object into a ClusterState", func() {
		// We builder a Cluster object by creating a json and using cmv1.UnmarshalCluster on it
		clusterObject := generateBaseClusterObject()

		//We convert the Cluster object into a ClusterState and check that the conversion is correct
		clusterState := &ClusterState{}
		clusterUtils := ClusterResourceUtilsImpl{}
		clusterUtils.populateClusterState(clusterObject, clusterState)

		Expect(clusterState.ID.Value).To(Equal(clusterId))
		Expect(clusterState.Version.Value).To(Equal(clusterVersion))
		Expect(clusterState.Product.Value).To(Equal(productId))
		Expect(clusterState.CloudProvider.Value).To(Equal(cloudProviderId))
		Expect(clusterState.CloudRegion.Value).To(Equal(regionId))
		Expect(clusterState.MultiAZ.Value).To(Equal(multiAz))
		Expect(clusterState.Properties.Elems["rosa_creator_arn"].Equal(types.String{Value: rosaCreatorArn})).To(Equal(true))
		Expect(clusterState.APIURL.Value).To(Equal(apiUrl))
		Expect(clusterState.ConsoleURL.Value).To(Equal(consoleUrl))
		Expect(clusterState.ComputeMachineType.Value).To(Equal(machineType))
		Expect(clusterState.AvailabilityZones.Elems).To(HaveLen(1))
		Expect(clusterState.AvailabilityZones.Elems[0].Equal(types.String{Value: availabilityZone1})).To(Equal(true))
		Expect(clusterState.CCSEnabled.Value).To(Equal(ccsEnabled))
		Expect(clusterState.AWSAccountID.Value).To(Equal(awsAccountID))
		Expect(clusterState.AWSAccessKeyID.Value).To(Equal(awsAccessKeyID))
		Expect(clusterState.AWSSecretAccessKey.Value).To(Equal(awsSecretAccessKey))
		Expect(clusterState.AWSPrivateLink.Value).To(Equal(privateLink))
	})

	Describe("Test the create function of the default clusterUtils", func() {
		var (
			ctrl *gomock.Controller
			clusterState *ClusterState
			clusterMap map[string]interface{}
			clusterObject *cmv1.Cluster
			mockDiags *MockDiagnostics
			mockClient *MockClusterClient

		)
		BeforeEach(func() {
			ctrl = generateGinkgoController()
			clusterState = generateBaseClusterState()
			clusterMap = generateBaseClusterMap()
			clusterObject = parseClusterMapToObject(clusterMap)
			mockDiags = NewMockDiagnostics(ctrl)
			mockClient = NewMockClusterClient(ctrl)
		})

		It("invokes cluster creation on the client", func() {
			comparator := func(x interface{}) bool {
				clusterObject := x.(*cmv1.Cluster)
				return clusterObject.Name() == clusterName
			}
			matcher := Matcher{
				comparator: comparator,
				text:       "cluster name is correct",
			}

			mockClient.EXPECT().Create(context.Background(), matcher).Return(clusterObject, nil)

			clusterUtils := ClusterResourceUtilsImpl{
				clusterClient: mockClient,
			}

			err := clusterUtils.create(context.Background(), clusterState, mockDiags)
			Expect(err).To(BeNil())
		})

		It("reports clusterCreationFailureMessage on client creation error", func() {
			mockDiags.EXPECT().AddError(clusterCreationFailureMessage, gomock.Any())

			comparator := func(x interface{}) bool {
				clusterObject := x.(*cmv1.Cluster)
				return clusterObject.Name() == clusterName
			}
			matcher := Matcher{
				comparator: comparator,
				text:       "cluster name is correct",
			}

			mockClient.EXPECT().Create(context.Background(), matcher).Return(clusterObject, fmt.Errorf("error"))

			clusterUtils := ClusterResourceUtilsImpl{
				clusterClient: mockClient,
			}

			err := clusterUtils.create(context.Background(), clusterState, mockDiags)
			Expect(err).ToNot(BeNil())
		})

		It("reports clusterPollFailure on client polling error", func() {
			clusterState.Wait.Value = true

			mockDiags.EXPECT().AddError(clusterPollFailure, gomock.Any())

			comparator := func(x interface{}) bool {
				clusterObject := x.(*cmv1.Cluster)
				return clusterObject.Name() == clusterName
			}
			matcher := Matcher{
				comparator: comparator,
				text:       "cluster name is correct",
			}

			mockClient.EXPECT().Create(context.Background(), matcher).Return(clusterObject, nil)
			mockClient.EXPECT().PollReady(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))

			clusterUtils := ClusterResourceUtilsImpl{
				clusterClient: mockClient,
			}

			err := clusterUtils.create(context.Background(), clusterState, mockDiags)
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Test the update function of the default clusterUtils", func() {
		var (
			ctrl *gomock.Controller
			clusterCurrentState *ClusterState
			clusterDesiredState *ClusterState
			clusterMap map[string]interface{}
			clusterObject *cmv1.Cluster
			mockDiags *MockDiagnostics
			mockClient *MockClusterClient

		)
		BeforeEach(func() {
			ctrl = generateGinkgoController()
			clusterCurrentState = generateBaseClusterState()
			clusterDesiredState = generateBaseClusterState()
			clusterMap = generateBaseClusterMap()
			clusterObject = parseClusterMapToObject(clusterMap)
			mockDiags = NewMockDiagnostics(ctrl)
			mockClient = NewMockClusterClient(ctrl)
		})

		It("invokes cluster update on the client", func() {
			clusterDesiredState.ComputeNodes.Value = clusterDesiredState.ComputeNodes.Value + 1
			clusterMap["nodes"].(map[string]interface{})["compute"] = clusterDesiredState.ComputeNodes.Value

			comparator := func(x interface{}) bool {
				clusterObject := x.(*cmv1.Cluster)
				return clusterObject.Nodes().Compute() == int(clusterDesiredState.ComputeNodes.Value)
			}
			matcher := Matcher{
				comparator: comparator,
				text:       "compute node conut is updated correctly",
			}

			
			mockClient.EXPECT().Update(context.Background(), clusterId, matcher).Return(clusterObject, nil)

			clusterUtils := ClusterResourceUtilsImpl{
				clusterClient: mockClient,
			}

			err := clusterUtils.update(context.Background(), clusterCurrentState, clusterDesiredState, mockDiags)
			Expect(err).To(BeNil())
		})

		It("reports clusterUpdateFailureMessage on client update error", func() {
			clusterDesiredState.ComputeNodes.Value = clusterDesiredState.ComputeNodes.Value + 1
			clusterMap["nodes"].(map[string]interface{})["compute"] = clusterDesiredState.ComputeNodes.Value

			mockDiags.EXPECT().AddError(clusterUpdateFailureMessage, gomock.Any())

			comparator := func(x interface{}) bool {
				clusterObject := x.(*cmv1.Cluster)
				return clusterObject.Nodes().Compute() == int(clusterDesiredState.ComputeNodes.Value)
			}
			matcher := Matcher{
				comparator: comparator,
				text:       "compute node conut is updated correctly",
			}

			mockClient.EXPECT().Update(context.Background(), clusterId, matcher).Return(clusterObject, fmt.Errorf("error"))

			clusterUtils := ClusterResourceUtilsImpl{
				clusterClient: mockClient,
			}

			err := clusterUtils.update(context.Background(), clusterCurrentState, clusterDesiredState, mockDiags)
			Expect(err).ToNot(BeNil())
		})

		It("reports clusterUpdateFailureMessage on client update error", func() {
			clusterDesiredState.ComputeNodes.Value = clusterDesiredState.ComputeNodes.Value + 1
			clusterMap["nodes"].(map[string]interface{})["compute"] = clusterDesiredState.ComputeNodes.Value

			mockDiags.EXPECT().AddError(clusterUpdateFailureMessage, gomock.Any())

			comparator := func(x interface{}) bool {
				clusterObject := x.(*cmv1.Cluster)
				return clusterObject.Nodes().Compute() == int(clusterDesiredState.ComputeNodes.Value)
			}
			matcher := Matcher{
				comparator: comparator,
				text:       "compute node conut is updated correctly",
			}

			mockClient.EXPECT().Update(context.Background(), clusterId, matcher).Return(clusterObject, fmt.Errorf("error"))

			clusterUtils := ClusterResourceUtilsImpl{
				clusterClient: mockClient,
			}

			err := clusterUtils.update(context.Background(), clusterCurrentState, clusterDesiredState, mockDiags)
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Test the delete function of the default clusterUtils", func() {
		var (
			ctrl *gomock.Controller
			mockDiags *MockDiagnostics
			mockClient *MockClusterClient

		)
		BeforeEach(func() {
			ctrl = generateGinkgoController()
			mockDiags = NewMockDiagnostics(ctrl)
			mockClient = NewMockClusterClient(ctrl)
		})

		It("invokes cluster deletion on the client", func() {
			mockClient.EXPECT().Delete(gomock.Any(), gomock.Eq(clusterId)).Return(nil)
			mockClient.EXPECT().PollRemoved(gomock.Any(), gomock.Eq(clusterId)).Return(nil)

			clusterUtils := ClusterResourceUtilsImpl{
				clusterClient: mockClient,
			}

			err := clusterUtils.delete(context.Background(), clusterId, true, mockDiags)
			Expect(err).To(BeNil())
		})

		It("reports clusterDeleteFailureMessage on client delete error", func() {
			mockDiags.EXPECT().AddError(clusterDeleteFailureMessage, gomock.Any())

			mockClient.EXPECT().Delete(gomock.Any(), gomock.Eq(clusterId)).Return(fmt.Errorf("error"))
			mockClient.EXPECT().PollRemoved(gomock.Any(), gomock.Eq(clusterId)).Return(nil)

			clusterUtils := ClusterResourceUtilsImpl{
				clusterClient: mockClient,
			}

			err := clusterUtils.delete(context.Background(), clusterId, true, mockDiags)
			Expect(err).ToNot(BeNil())
		})

		It("reports clusterPollDeletionFailure on client deletion polling error", func() {
			mockDiags.EXPECT().AddError(clusterPollDeletionFailure, gomock.Any())

			mockClient.EXPECT().Delete(gomock.Any(), gomock.Eq(clusterId)).Return(nil)
			mockClient.EXPECT().PollRemoved(gomock.Any(), gomock.Eq(clusterId)).Return(fmt.Errorf("error"))

			clusterUtils := ClusterResourceUtilsImpl{
				clusterClient: mockClient,
			}

			err := clusterUtils.delete(context.Background(), clusterId, true, mockDiags)
			Expect(err).ToNot(BeNil())
		})
	})

})

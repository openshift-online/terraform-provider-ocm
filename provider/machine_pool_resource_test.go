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
	"net/http"

	. "github.com/onsi/ginkgo/v2/dsl/core"             // nolint
	. "github.com/onsi/gomega"                         // nolint
	. "github.com/onsi/gomega/ghttp"                   // nolint
	. "github.com/openshift-online/ocm-sdk-go/testing" // nolint
)

var _ = Describe("Machine pool creation", func() {
	BeforeEach(func() {
		// The first thing that the provider will do for any operation on machine pools
		// is check that the cluster is ready, so we always need to prepare the server to
		// respond to that:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodGet, "/api/clusters_mgmt/v1/clusters/123"),
				RespondWithJSON(http.StatusOK, `{
				  "id": "123",
				  "name": "my-cluster",
				  "state": "ready"
				}`),
			),
		)
	})

	It("Can create machine pool", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodPost,
					"/api/clusters_mgmt/v1/clusters/123/machine_pools",
				),
				VerifyJSON(`{
				  "kind": "MachinePool",
				  "id": "my-pool",
				  "instance_type": "r5.xlarge",
				  "replicas": 10
				}`),
				RespondWithJSON(http.StatusOK, `{
				  "id": "my-pool",
				  "instance_type": "r5.xlarge",
				  "replicas": 10
				}`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_machine_pool" "my_pool" {
		    cluster      = "123"
		    name         = "my-pool"
		    machine_type = "r5.xlarge"
		    replicas     = 10
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())

		// Check the state:
		resource := terraform.Resource("ocm_machine_pool", "my_pool")
		Expect(resource).To(MatchJQ(".attributes.cluster", "123"))
		Expect(resource).To(MatchJQ(".attributes.id", "my-pool"))
		Expect(resource).To(MatchJQ(".attributes.name", "my-pool"))
		Expect(resource).To(MatchJQ(".attributes.machine_type", "r5.xlarge"))
		Expect(resource).To(MatchJQ(".attributes.replicas", 10.0))
	})
})

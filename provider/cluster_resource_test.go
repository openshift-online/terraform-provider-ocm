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

	. "github.com/onsi/ginkgo"                         // nolint
	. "github.com/onsi/gomega"                         // nolint
	. "github.com/onsi/gomega/ghttp"                   // nolint
	. "github.com/openshift-online/ocm-sdk-go/testing" // nolint
)

var _ = Describe("Cluster creation", func() {
	// This is the cluster that will be returned by the server when asked to create or retrieve
	// a cluster.
	const template = `{
	  "id": "123",
	  "name": "my-cluster",
	  "cloud_provider": {
	    "id": "aws"
	  },
	  "region": {
	    "id": "us-west-1"
	  },
	  "multi_az": false,
	  "properties": {},
	  "api": {
	    "url": "https://my-api.example.com"
	  },
	  "console": {
	    "url": "https://my-console.example.com"
	  },
	  "nodes": {
	    "compute": 3,
	    "compute_machine_type": {
	      "id": "r5.xlarge"
	    }
	  },
	  "ccs": {
	    "enabled": false
	  },
	  "network": {
	    "machine_cidr": "10.0.0.0/16",
	    "service_cidr": "172.30.0.0/16",
	    "pod_cidr": "10.128.0.0/14",
	    "host_prefix": 23
	  },
	  "version": {
		  "id": "openshift-4.8.0"
	  },
	  "state": "ready"
	}`

	It("Creates basic cluster", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.name`, "my-cluster"),
				VerifyJQ(`.cloud_provider.id`, "aws"),
				VerifyJQ(`.region.id`, "us-west-1"),
				RespondWithJSON(http.StatusCreated, template),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster" "my_cluster" {
		    name           = "my-cluster"
		    cloud_provider = "aws"
		    cloud_region   = "us-west-1"
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})

	It("Saves API and console URLs to the state", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				RespondWithJSON(http.StatusCreated, template),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster" "my_cluster" {
		    name           = "my-cluster"
		    cloud_provider = "aws"
		    cloud_region   = "us-west-1"
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())

		// Check the state:
		resource := terraform.Resource("ocm_cluster", "my_cluster")
		Expect(resource).To(MatchJQ(".attributes.api_url", "https://my-api.example.com"))
		Expect(resource).To(MatchJQ(".attributes.console_url", "https://my-console.example.com"))
	})

	It("Sets compute nodes and machine type", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(`.nodes.compute`, 3.0),
				VerifyJQ(`.nodes.compute_machine_type.id`, "r5.xlarge"),
				RespondWithJSON(http.StatusCreated, template),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster" "my_cluster" {
		    name                 = "my-cluster"
		    cloud_provider       = "aws"
		    cloud_region         = "us-west-1"
		    compute_nodes        = 3
		    compute_machine_type = "r5.xlarge"
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())

		// Check the state:
		resource := terraform.Resource("ocm_cluster", "my_cluster")
		Expect(resource).To(MatchJQ(".attributes.compute_nodes", 3.0))
		Expect(resource).To(MatchJQ(".attributes.compute_machine_type", "r5.xlarge"))
	})

	It("Creates CCS cluster", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(".ccs.enabled", true),
				VerifyJQ(".aws.account_id", "123"),
				VerifyJQ(".aws.access_key_id", "456"),
				VerifyJQ(".aws.secret_access_key", "789"),
				RespondWithPatchedJSON(http.StatusOK, template, `[
				  {
				    "op": "replace",
				    "path": "/ccs",
				    "value": {
				      "enabled": true
				    }
				  },
				  {
				    "op": "add",
				    "path": "/aws",
				    "value": {
				      "account_id": "123",
				      "access_key_id": "456",
				      "secret_access_key": "789"
				    }
				  }
				]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster" "my_cluster" {
		    name                  = "my-cluster"
		    cloud_provider        = "aws"
		    cloud_region          = "us-west-1"
		    ccs_enabled           = true
		    aws_account_id        = "123"
		    aws_access_key_id     = "456"
		    aws_secret_access_key = "789"
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())

		// Check the state:
		resource := terraform.Resource("ocm_cluster", "my_cluster")
		Expect(resource).To(MatchJQ(".attributes.ccs_enabled", true))
		Expect(resource).To(MatchJQ(".attributes.aws_account_id", "123"))
		Expect(resource).To(MatchJQ(".attributes.aws_access_key_id", "456"))
		Expect(resource).To(MatchJQ(".attributes.aws_secret_access_key", "789"))
	})

	It("Sets network configuration", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(".network.machine_cidr", "10.0.0.0/15"),
				VerifyJQ(".network.service_cidr", "172.30.0.0/15"),
				VerifyJQ(".network.pod_cidr", "10.128.0.0/13"),
				VerifyJQ(".network.host_prefix", 22.0),
				RespondWithPatchedJSON(http.StatusOK, template, `[
				  {
				    "op": "replace",
				    "path": "/network",
				    "value": {
				      "machine_cidr": "10.0.0.0/15",
				      "service_cidr": "172.30.0.0/15",
				      "pod_cidr": "10.128.0.0/13",
				      "host_prefix": 22
				    }
				  }
				]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster" "my_cluster" {
		    name           = "my-cluster"
		    cloud_provider = "aws"
		    cloud_region   = "us-west-1"
		    machine_cidr   = "10.0.0.0/15"
		    service_cidr   = "172.30.0.0/15"
		    pod_cidr       = "10.128.0.0/13"
		    host_prefix    = 22
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())

		// Check the state:
		resource := terraform.Resource("ocm_cluster", "my_cluster")
		Expect(resource).To(MatchJQ(".attributes.machine_cidr", "10.0.0.0/15"))
		Expect(resource).To(MatchJQ(".attributes.service_cidr", "172.30.0.0/15"))
		Expect(resource).To(MatchJQ(".attributes.pod_cidr", "10.128.0.0/13"))
		Expect(resource).To(MatchJQ(".attributes.host_prefix", 22.0))
	})

	It("Sets version", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(".version.id", "openshift-v4.8.1"),
				RespondWithPatchedJSON(http.StatusOK, template, `[
				  {
				    "op": "replace",
				    "path": "/version",
				    "value": {
				      "id": "openshift-v4.8.1"
				    }
				  }
				]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster" "my_cluster" {
		    name           = "my-cluster"
		    cloud_provider = "aws"
		    cloud_region   = "us-west-1"
		    version        = "openshift-v4.8.1"
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())

		// Check the state:
		resource := terraform.Resource("ocm_cluster", "my_cluster")
		Expect(resource).To(MatchJQ(".attributes.version", "openshift-v4.8.1"))
	})

	It("Sets STS attributes", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(
					`.aws.sts.role_arn`,
					`arn:aws:iam::000000000000:role/my-role`,
				),
				VerifyJQ(
					`.aws.sts.support_role_arn`,
					`arn:aws:iam::000000000000:role/my-support-role`,
				),
				VerifyJQ(
					`.aws.sts.operator_iam_roles[] | select(.namespace == "openshift-machine-api") | .name`,
					`aws-cloud-credentials`,
				),
				VerifyJQ(
					`.aws.sts.operator_iam_roles[] | select(.namespace == "openshift-machine-api") | .role_arn`,
					`arn:aws:iam::000000000000:role/my-machine-role`,
				),
				VerifyJQ(
					`.aws.sts.operator_iam_roles[] | select(.namespace == "openshift-cloud-credential-operator") | .name`,
					`cloud-credential-operator-iam-ro-creds`,
				),
				VerifyJQ(
					`.aws.sts.operator_iam_roles[] | select(.namespace == "openshift-cloud-credential-operator") | .role_arn`,
					`arn:aws:iam::000000000000:role/my-cloud-role`,
				),
				VerifyJQ(
					`.aws.sts.operator_iam_roles[] | select(.namespace == "openshift-image-registry") | .name`,
					`installer-cloud-credentials`,
				),
				VerifyJQ(
					`.aws.sts.operator_iam_roles[] | select(.namespace == "openshift-image-registry") | .role_arn`,
					`arn:aws:iam::000000000000:role/my-registry-role`,
				),
				VerifyJQ(
					`.aws.sts.operator_iam_roles[] | select(.namespace == "openshift-ingress-operator") | .name`,
					`cloud-credentials`,
				),
				VerifyJQ(
					`.aws.sts.operator_iam_roles[] | select(.namespace == "openshift-ingress-operator") | .role_arn`,
					`arn:aws:iam::000000000000:role/my-ingress-role`,
				),
				VerifyJQ(
					`.aws.sts.operator_iam_roles[] | select(.namespace == "openshift-cluster-csi-drivers") | .name`,
					`ebs-cloud-credentials`,
				),
				VerifyJQ(
					`.aws.sts.operator_iam_roles[] | select(.namespace == "openshift-cluster-csi-drivers") | .role_arn`,
					`arn:aws:iam::000000000000:role/my-csi-role`,
				),
				VerifyJQ(
					`.aws.sts.instance_iam_roles.master_role_arn`,
					`arn:aws:iam::000000000000:role/my-control-role`,
				),
				VerifyJQ(
					`.aws.sts.instance_iam_roles.worker_role_arn`,
					`arn:aws:iam::000000000000:role/my-worker-role`,
				),
				RespondWithPatchedJSON(http.StatusOK, template, `[
				  {
				    "op": "replace",
				    "path": "/ccs",
				    "value": {
				      "enabled": true
				    }
				  },
				  {
				    "op": "add",
				    "path": "/aws",
				    "value": {
				      "sts": {
				        "role_arn": "arn:aws:iam::000000000000:role/my-role",
				        "support_role_arn": "arn:aws:iam::000000000000:role/my-support-role",
				        "operator_iam_roles": [
				          {
				            "namespace": "openshift-machine-api",
				            "name": "aws-cloud-credentials",
				            "role_arn": "arn:aws:iam::000000000000:role/my-machine-role"
				          },
				          {
				            "namespace": "openshift-cloud-credential-operator",
				            "name": "cloud-credential-operator-iam-ro-creds",
				            "role_arn": "arn:aws:iam::000000000000:role/my-cloud-role"
				          },
				          {
				            "namespace": "openshift-image-registry",
				            "name": "installer-cloud-credentials",
				            "role_arn": "arn:aws:iam::000000000000:role/my-registry-role"
				          },
				          {
				            "namespace": "openshift-ingress-operator",
				            "name": "cloud-credentials",
				            "role_arn": "arn:aws:iam::000000000000:role/my-ingress-role"
				          },
				          {
				            "namespace": "openshift-cluster-csi-drivers",
				            "name": "ebs-cloud-credentials", 
				            "role_arn": "arn:aws:iam::000000000000:role/my-csi-role"
				          }
				        ],
				        "instance_iam_roles": {
				          "master_role_arn": "arn:aws:iam::00000000000:role/my-control-role",
				          "worker_role_arn": "arn:aws:iam::00000000000:role/my-worker-role"
				        }
				      }
				    }
				  }
				]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster" "my_cluster" {
		    name              = "my-cluster"
		    cloud_provider    = "aws"
		    cloud_region      = "us-west-1"
		    ccs_enabled       = true
		    sts_role          = "arn:aws:iam::000000000000:role/my-role"
		    sts_support_role  = "arn:aws:iam::000000000000:role/my-support-role"
		    sts_machine_role  = "arn:aws:iam::000000000000:role/my-machine-role"
		    sts_cloud_role    = "arn:aws:iam::000000000000:role/my-cloud-role"
		    sts_registry_role = "arn:aws:iam::000000000000:role/my-registry-role"
		    sts_ingress_role  = "arn:aws:iam::000000000000:role/my-ingress-role"
		    sts_csi_role      = "arn:aws:iam::000000000000:role/my-csi-role"
		    sts_control_role  = "arn:aws:iam::000000000000:role/my-control-role"
		    sts_worker_role   = "arn:aws:iam::000000000000:role/my-worker-role"
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())

		// Check the state:
		resource := terraform.Resource("ocm_cluster", "my_cluster")
		Expect(resource).To(MatchJQ(
			".attributes.sts_role",
			"arn:aws:iam::000000000000:role/my-role",
		))
		Expect(resource).To(MatchJQ(
			".attributes.sts_support_role",
			"arn:aws:iam::000000000000:role/my-support-role",
		))
		Expect(resource).To(MatchJQ(
			".attributes.sts_machine_role",
			"arn:aws:iam::000000000000:role/my-machine-role",
		))
		Expect(resource).To(MatchJQ(
			".attributes.sts_cloud_role",
			"arn:aws:iam::000000000000:role/my-cloud-role",
		))
		Expect(resource).To(MatchJQ(
			".attributes.sts_registry_role",
			"arn:aws:iam::000000000000:role/my-registry-role",
		))
		Expect(resource).To(MatchJQ(
			".attributes.sts_ingress_role",
			"arn:aws:iam::000000000000:role/my-ingress-role",
		))
		Expect(resource).To(MatchJQ(
			".attributes.sts_csi_role",
			"arn:aws:iam::000000000000:role/my-csi-role",
		))
		Expect(resource).To(MatchJQ(
			".attributes.sts_control_role",
			"arn:aws:iam::000000000000:role/my-control-role",
		))
		Expect(resource).To(MatchJQ(
			".attributes.sts_worker_role",
			"arn:aws:iam::000000000000:role/my-worker-role",
		))
	})

	It("Fails if the cluster already exists", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				RespondWithJSON(http.StatusBadRequest, `{
				  "id": "400",
				  "code": "CLUSTERS-MGMT-400",
				  "reason": "Cluster 'my-cluster' already exists"
				}`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster" "my_cluster" {
		    name           = "my-cluster"
		    cloud_provider = "aws"
		    cloud_region   = "us-west-1"
		  }
		`)
		Expect(terraform.Apply()).ToNot(BeZero())
	})
})

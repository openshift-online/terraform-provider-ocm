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
	"encoding/json"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2/dsl/core"             // nolint
	. "github.com/onsi/gomega"                         // nolint
	. "github.com/onsi/gomega/ghttp"                   // nolint
	. "github.com/openshift-online/ocm-sdk-go/testing" // nolint
)

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

var _ = Describe("Cluster creation", func() {
	// This is the cluster that will be returned by the server when asked to create or retrieve
	// a cluster.
	const template = `{
	  "id": "123",
	  "product": {
		"id": "osd"
	  },
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
			product		   = "osd"
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
			product		   = "osd"
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
			product		   		 = "osd"
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
			product		   		  = "osd"
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
			product		   = "osd"
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
			product		   = "osd"
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

	It("Sets additional trust bundle", func() {
		test_value := `-----BEGIN PRIVATE KEY-----
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
nX0+dzIk7LHwnRVF66MSfKjXMxhp+mSL
-----END PRIVATE KEY-----
-----BEGIN CERTIFICATE-----
MIIFmDCCA4ACCQC+egbkLjxMpzANBgkqhkiG9w0BAQsFADCBjTELMAkGA1UEBhMC
VVMxEjAQBgNVBAgMCVdhc2hpbmdvbjEQMA4GA1UEBwwHU2VhdHRsZTEQMA4GA1UE
CgwHUmVkIEhhdDENMAsGA1UECwwETU9CQjETMBEGA1UEAwwKbW9iYi1wcm94eTEi
MCAGCSqGSIb3DQEJARYTdGhodWJiYXJAcmVkaGF0LmNvbTAeFw0yMjExMjMxODI4
MDRaFw0yMzExMjMxODI4MDRaMIGNMQswCQYDVQQGEwJVUzESMBAGA1UECAwJV2Fz
aGluZ29uMRAwDgYDVQQHDAdTZWF0dGxlMRAwDgYDVQQKDAdSZWQgSGF0MQ0wCwYD
VQQLDARNT0JCMRMwEQYDVQQDDAptb2JiLXByb3h5MSIwIAYJKoZIhvcNAQkBFhN0
aGh1YmJhckByZWRoYXQuY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKC
AgEA4VozLPHboVi+8lEeLF0Bts/MZbOivssMc7Q8S6IbbPekEK5O+O1N2By1toZh
CWgyuNXNmMnwpCEaVklyn3KAyNImCWwS1GVqKekCWAIxZNmJcBo7PVIeQ1k5BtLu
ZKdF9CG1Xe/TJ1la8/mCjWW6ZWRM2V/LP3+4tS3V6/YsZ6O5K2riue3JueOCmKWc
9JAqOC9wY3EnVodT2QM8qUR+/dCL7yOpv0vYK/NmyBcrlU2obs9dQYbawmVyPf9W
xHJm5XM6BEweTvpp1CZg9HmctQqeFvswlauHCWX7HGXAScTr/2cJ8sIzElaXTYeT
vt9vtMNU6On2tGZQua1iaL3GtGntcE4OSAQmyCNsxuTNiNps1GSMchFt1Xtvbe1O
sY6gLpXjgjr55PgFrRvHzWYtfdtgfmJe5FxZMN5SsaGpxY2jfBrcsKlY6tlhyNR4
UpsJ4DCcWgbbHhDRAyIzC2L0409MJdkGbOz67YfaCpaHV6jl0ulrPbA2XQpZABn1
ibQka3h4nQA/sVyVCZPR5rVQHg/NEde3jVl1XPShfOWvQtsahSLYSc9+nGT4Rhcf
mkQn6mscN1qQVWrIrEm9zcWrxK/uKH0iaZAhcTHqUh6j8KmfODKn2nccekqsSh5S
qk37IfoTWs0lu3KVariw7zJjJ0jhgox7UQCvBEZi16O6U2MCAwEAATANBgkqhkiG
9w0BAQsFAAOCAgEAAihyYSYH6RvVF3zvNY3Mux2/uUvPOufS7kDcZpzTyCElc6yQ
dBDc4/u0Q6WWtEwLx7E7PqJT/vFKm2qpnMS4qjhLUJz62E1aT6qVnLGhZz0rFYHO
oCv2H08BUmzu3Gn1wRXvpxfjy5/ca1AJO7pm+GTdk6TLgLzkfHXseFkEcUQd8mwL
fwUcUjYJJoqZIwlDfseEs6vGX2lDRC6X1jev5N1A1wRr/a3V2lUxRjeEEyyLo0aY
NqAuLLkJP23qmXWhykiyXz69jjacmz1qr7bSVvmSCUxEGKIaeKvEwg1lUknsisKo
M3gEdjDYXZl86jQzqPiQPZFbKJHxNxguKcRZS2a+zy1p8lm+5CGPD2lB9j2qB5Yw
5YikqbsOIxUY5M+dFrsq8fpMcQjI2JOSXp5a23bF0iRdjYl+DcVNkvPBV+YWSDum
dAapVfCLcg+Mrswnf0OKd5vo4IquLE+chNcAEDs+dFZSiXo5y1UsmXUEKwesqd50
R5Fbq9B+RiPubxoP+bNc9O67JO1nx0NGxIHf59vCnWdSs6cAR8o9VOm40tGBxJ2j
TV4h0KjmDAJi27m1sfNJbMXrp9ng3X+GLM+R62NcO0huymPLq/RrZgcrRN63+ZbX
WlOQEm17QK7K3stHjDFVPZ7XJf/ys1cucBJKd8MpYchuwLwEVepdKLQqFYU=
-----END CERTIFICATE-----`

		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(http.MethodPost, "/api/clusters_mgmt/v1/clusters"),
				VerifyJQ(".additional_trust_bundle", test_value+"\n"),
				RespondWithPatchedJSON(http.StatusOK, template, `[
				  {
				    "op": "add",
				    "path": "/additional_trust_bundle",
				    "value": "`+strings.ReplaceAll(test_value, "\n", `\n`)+`"
				  }
				]`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_cluster" "my_cluster" {
		    name           = "my-cluster"
			product		   = "osd"
		    cloud_provider = "aws"
		    cloud_region   = "us-west-1"
		    additional_trust_bundle = <<EOT
` + test_value + `
EOT
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())

		// Check the state:
		resource := terraform.Resource("ocm_cluster", "my_cluster")
		Expect(resource).To(MatchJQ(".attributes.additional_trust_bundle", test_value+"\n"))
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
			product		   = "osd"
		    name           = "my-cluster"
		    cloud_provider = "aws"
		    cloud_region   = "us-west-1"
		  }
		`)
		Expect(terraform.Apply()).ToNot(BeZero())
	})
})

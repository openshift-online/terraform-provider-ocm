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

var _ = Describe("Identity provider creation", func() {
	BeforeEach(func() {
		// The first thing that the provider will do for any operation on identity providers
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

	It("Can create a 'htpasswd' identity provider", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodPost,
					"/api/clusters_mgmt/v1/clusters/123/identity_providers",
				),
				VerifyJSON(`{
				  "kind": "IdentityProvider",
				  "type": "HTPasswdIdentityProvider",
				  "name": "my-ip",
				  "htpasswd": {
				    "password": "my-password",
				    "username": "my-user"
				  }
				}`),
				RespondWithJSON(http.StatusOK, `{
				  "id": "456",
				  "name": "my-ip",
				  "htpasswd": {
				    "user": "my-user"
				  }
				}`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_identity_provider" "my_ip" {
		    cluster = "123"
		    name    = "my-ip"
		    htpasswd = {
		      username = "my-user"
		      password = "my-password"
		    }
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})

	It("Can create an LDAP identity provider", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodPost,
					"/api/clusters_mgmt/v1/clusters/123/identity_providers",
				),
				VerifyJSON(`{
				  "kind": "IdentityProvider",
				  "type": "LDAPIdentityProvider",
				  "name": "my-ip",
				  "ldap": {
				    "bind_dn": "my-bind-dn",
				    "bind_password": "my-bind-password",
				    "ca": "my-ca",
				    "insecure": false,
				    "url": "ldap://my-server.com",
				    "attributes": {
				      "id": ["my-id"],
				      "email": ["my-email"],
				      "name": ["my-name"],
				      "preferred_username": ["my-preferred-username"]
				    }
				  }
				}`),
				RespondWithJSON(http.StatusOK, `{
				  "id": "456",
				  "name": "my-ip",
				  "ldap": {
				    "bind_dn": "my-bind-dn",
				    "bind_password": "my-bind-password",
				    "ca": "my-ca",
				    "insecure": false,
				    "url": "ldap://my-server.com",
				    "attributes": {
				      "id": ["my-id"],
				      "email": ["my-email"],
				      "name": ["my-name"],
				      "preferred_username": ["my-preferred-username"]
				    }
				  }
				}`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_identity_provider" "my_ip" {
		    cluster    = "123"
		    name       = "my-ip"
		    ldap = {
		      bind_dn       = "my-bind-dn"
		      bind_password = "my-bind-password"
		      insecure      = false
		      ca            = "my-ca"
		      url           = "ldap://my-server.com"
		      attributes    = {
		        id                 = ["my-id"]
		        email              = ["my-email"]
		        name               = ["my-name"]
		        preferred_username = ["my-preferred-username"]
		      }
		    }
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})

	It("Can create an OpenID identity provider", func() {
		// Prepare the server:
		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest(
					http.MethodPost,
					"/api/clusters_mgmt/v1/clusters/123/identity_providers",
				),
				VerifyJSON(`{
				  "kind": "IdentityProvider",
				  "type": "OpenIDIdentityProvider",
				  "name": "my-ip",
				  "mapping_method": "claim",
				  "open_id": {
					"ca": "test_ca",
					"claims": {
						"email": [
							"email"
						],
						"groups": [
							"admins"
						],
						"name": [
							"name",
							"email"
						],
						"preferred_username": [
							"preferred_username",
							"email"
						]
					},
					"client_id": "test_client",
					"client_secret": "test_secret",
					"extra_scopes": [
					  "email",
					  "profile"
					],
					"issuer": "https://test.okta.com"
					}
				}`),
				RespondWithJSON(http.StatusOK, `{
					"kind": "IdentityProvider",
					"type": "OpenIDIdentityProvider",
					"href": "/api/clusters_mgmt/v1/clusters/123/identity_providers/456",
					"id": "456",
					"name": "my-ip",
					"mapping_method": "claim",
					"open_id": {
						"claims": {
							"email": [
								"email"
							],
							"groups": [
								"admins"
							],
							"name": [
								"name",
								"email"
							],
							"preferred_username": [
								"preferred_username",
								"email"
							]
						},
						"client_id": "test_client",
						"extra_scopes": [
							"email",
							"profile"
						],
						"issuer": "https://test.okta.com"
					}
				}`),
			),
		)

		// Run the apply command:
		terraform.Source(`
		  resource "ocm_identity_provider" "my_ip" {
		    cluster    				= "123"
		    name       				= "my-ip"
			mapping_method			= "claim"
		    openid = {
				ca            			= "test_ca"
				issuer					= "https://test.okta.com"
				client_id 				= "test_client"
				client_secret			= "test_secret"
				extra_scopes 			= ["email","profile"]
				claims = {
					email              = ["email"]
					groups			   = ["admins"]
					name               = ["name","email"]
					preferred_username = ["preferred_username","email"]
		      	}
		    }
		  }
		`)
		Expect(terraform.Apply()).To(BeZero())
	})

})

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

const (
	clusterId = "1n2j3k4l5m6n7o8p9q0r"
	clusterName = "my-cluster"
	clusterVersion = "openshift-v4.11.12"
	productId = "rosa"
	cloudProviderId = "aws"
	regionId = "us-east-1"
	multiAz = true
	rosaCreatorArn = "arn:aws:iam::123456789012:dummy/dummy"
	apiUrl = "https://api.my-cluster.com:6443"
	consoleUrl = "https://console.my-cluster.com"
	machineType = "m5.xlarge"
	availabilityZone1 = "us-east-1a"
	availabilityZone2 = "us-east-1b"
	ccsEnabled = true
	awsAccountID = "123456789012"
	awsAccessKeyID = "AKIAIOSFODNN7EXAMPLE"
	awsSecretAccessKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	privateLink = false
	computeNodes = 3
	oidcEndpointUrl   = "example.com"
	roleArn           = "arn:aws:iam::123456789012:role/role-name"
	httpProxy         = "http://proxy.com"
	httpsProxy        = "https://proxy.com"
)

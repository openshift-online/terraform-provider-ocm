#
# Copyright (c) 2021 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 3.67"
    }
    ocm = {
      version = ">= 0.1"
      source  = "openshift-online/ocm"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

provider "ocm" {
}

data "aws_caller_identity" "current" {
}
resource "ocm_cluster" "my_cluster" {
  name              = "my-cluster"
  cloud_provider    = "aws"
  cloud_region      = "us-east-1"
  ccs_enabled       = true
  aws_account_id    = data.aws_caller_identity.current.account_id
  sts_role          = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/my-role"
  sts_support_role  = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/my-support-role"
  sts_machine_role  = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/my-machine-role"
  sts_cloud_role    = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/my-cloud-role"
  sts_registry_role = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/my-registry-role"
  sts_ingress_role  = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/my-ingress-role"
  sts_csi_role      = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/my-csi-role"
  sts_control_role  = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/my-control-role"
  sts_worker_role   = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/my-worker-role"
}

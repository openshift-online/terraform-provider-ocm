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
    ocm = {
      version = ">= 0.1"
      source  = "rh-mobb/ocm"
    }
  }
}

provider "ocm" {
}

data "ocm_cloud_providers" "a" {
  search = "display_name like 'A%'"
  order  = "display_name asc"
}

output "result" {
  description = "Cloud providers with names starting with 'A'"
  value       = data.ocm_cloud_providers.a.items
}

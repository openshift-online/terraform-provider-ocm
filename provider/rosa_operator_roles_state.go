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

import "github.com/hashicorp/terraform-plugin-framework/types"

type RosaOperatorRolesState struct {
	ClusterID          types.String       `tfsdk:"cluster_id"`
	OperatorRolePrefix types.String       `tfsdk:"operator_role_prefix"`
	AccountRolePrefix  types.String       `tfsdk:"account_role_prefix"`
	OperatorIAMRoles   []*OperatorIAMRole `tfsdk:"operator_iam_roles"`
}

type OperatorIAMRole struct {
	Name            types.String `tfsdk:"operator_name"`
	Namespace       types.String `tfsdk:"operator_namespace"`
	RoleName        types.String `tfsdk:"role_name"`
	PolicyName      types.String `tfsdk:"policy_name"`
	ServiceAccounts types.List   `tfsdk:"service_accounts"`
}
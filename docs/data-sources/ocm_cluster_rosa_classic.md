---
page_title: "ocm_cluster_rosa_classic Data Source"
subcategory: ""
description: |-
  Get ROSA cluster details.
---

# ocm_cluster_rosa_classic (Data Source)

Get ROSA cluster details.

## Schema

### Optional

- **name** (String) Name of the cluster.

- **id** (String) Unique identifier of the cluster.

Either name or id must be provided.

  
### Read-Only

- **version** (String) Version of OpenShift used to create the cluster

- **state** (String) State of the cluster.
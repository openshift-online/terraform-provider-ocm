terraform {
  required_providers {
    ocm = {
      version = ">= 0.1"
      source  = "openshift-online/ocm"
    }
  }
}

provider "ocm" {
}

data "ocm_cluster_rosa_classic" "mycluster-name" {
    name = "my-cluster"
}

output "cluster_id"{
    value = data.ocm_cluster_rosa_classic.mycluster-name.id //outputs: cluster_id = "1n2j3k4l5m6n7o8p9q0r"
}

data "ocm_cluster_rosa_classic" "mycluster-id" {
    id = "1n2j3k4l5m6n7o8p9q0r"
}

output "cluster_name"{
    value = data.ocm_cluster_rosa_classic.mycluster-id.name //outputs: cluster_name = "my-cluster"
}

output "cluster_status"{
    value = data.ocm_cluster_rosa_classic.mycluster-id.state //outputs: cluster_status = "<cluster state>"
}

output "cluster_version"{
    value = data.ocm_cluster_rosa_classic.mycluster-id.version // outputs: cluster_version = "<x.x.x>"

}

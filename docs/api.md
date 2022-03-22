This doc descibes REST API for ncn lifecycle management. Note that in this version, we only provide APIs for individual operation. A full end to end lifecycle management API is out of scope in Phase I

---

## Kubernetes Nodes

#### e2e upgrade flow

1. `/etcd/{hostname}/prepare`
   > NOTE: no-op for **worker** nodes
1. `/kubernetes/{hostname}/pre-rebuild`
1. `/kubernetes/{hostname}/drain`
1. `/ncn/{hostname}/backup`
1. `/ncn/{hostname}/wipe`
1. PUT `/ncn/{hostname}/boot-parameters`
1. `/ncn/{hostname}/reboot`

   > NOTE: how do we wait for boot? maybe wait for ncn ready on k8s?

1. `/ncn/{hostname}/restore`
1. `/ncn/{hostname}/post-rebuild`
1. `/kubernetes/{hostname}/post-rebuild`
1. `/ncn/{hostname}/validate`

##### After all Kubernetes nodes are upgraded

1. `/ncn/kubernetes/post-upgrade`

---

## Ceph Storage Node

---

[API Doc](swagger.md)

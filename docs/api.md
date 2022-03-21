This doc descibes REST API for ncn lifecycle management. Note that in this version, we only provide APIs for individual operation. A full end to end lifecycle management API is out of scope in Phase I

---

## Kubernetes Nodes

#### e2e upgrade flow

1. `/etcd/{hostname}/prepare`
1. `/kubernetes/{hostname}/pre-rebuild`
1. `/kubernetes/{hostname}/drain`
1. `/ncn/{hostname}/backup`
1. `/ncn/{hostname}/wipe`
1. `/ncn/{hostname}/reboot`

   > NOTE: how do we wait for boot? maybe wait for ncn ready on k8s?

1. `/ncn/{hostname}/restore`
1. `/ncn/{hostname}/post-rebuild`
1. `/kubernetes/{hostname}/post-rebuild`
1. `/ncn/{hostname}/validate`

---

## Kubernetes Worker node

---

## Ceph Storage Node

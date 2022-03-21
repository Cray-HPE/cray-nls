## Prepare baremetal ETCD for rejoining

Prepare a master ncn to rejoin baremetal etcd cluster

#### Pre-condition

1. **NCN** is a **master** node
1. Baremetal etcd cluster is in **healthy** state
1. quorum after removal

#### Action

1. Remove a ncn from baremetal etcd cluster
1. Stop etcd services on the ncn
1. Add the ncn back to etcd cluster so it can rejoin on boot

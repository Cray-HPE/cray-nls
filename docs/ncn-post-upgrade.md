## NCN Post Upgrade

After all ncn of a certain type has been rebuilt, some `CSM specific` steps are required.

---

### Master

##### Pre-condition

1. **NCN** is a **master/worker** 

#### Actions

1. `/srv/cray/scripts/common/apply-networking-manifests.sh`
    NOTE: this is taking quite long. we may want to use async here


1. `/usr/share/doc/csm/upgrade/1.2/scripts/k8s/apply-coredns-pod-affinity.sh`


1. `/usr/share/doc/csm/upgrade/1.2/scripts/k8s/upgrade_control_plane.sh`

---

### Storage

##### Pre-condition
1. **NCN** is a **storage** 

#### Actions
1. Deploy node-exporter and alertmanager


1. Update BSS to ensure the Ceph images are loaded if a node is rebuilt

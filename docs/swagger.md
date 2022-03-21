# NCN Lifecycle Management API
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
1. `/ncn/{hostname}/reboot`

   > NOTE: how do we wait for boot? maybe wait for ncn ready on k8s?

1. `/ncn/{hostname}/restore`
1. `/ncn/{hostname}/post-rebuild`
1. `/kubernetes/{hostname}/post-rebuild`
1. `/ncn/{hostname}/validate`

---

## Ceph Storage Node

---

[API Doc](swagger.md)

## Version: 1.0

**License:** [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)

### /etcd/{hostname}/prepare

#### PUT
##### Summary

Prepare baremetal etcd for a master node to rejoin

##### Description

# Prepare baremetal ETCD for rejoining

Prepare a master ncn to rejoin baremetal etcd cluster

## Pre-condition

1. **NCN** is a **master** node
1. Baremetal etcd cluster is in **healthy** state
1. quorum after removal

## Action

1. Remove a ncn from baremetal etcd cluster
1. Stop etcd services on the ncn
1. Add the ncn back to etcd cluster so it can rejoin on boot

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname of target ncn | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | ok | string |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /kubernetes/{hostname}/drain

#### POST
##### Summary

Drain a Kubernetes node

##### Description

# Drain Kubernetes Node

Before we can safely drain/remove a node from k8s cluster, we need to run some `CSM specific logic` to make sure a node can be drained from k8s cluster safely

---

#### Pre-condition

1. **NCN** is a **master** node
1. quorum after removal

#### Actions

1. drain node

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /kubernetes/{hostname}/post-rebuild

#### POST
##### Summary

Kubernetes node post rebuild action

##### Description

# K8s Post Rebuild

After a node rejoined k8s cluster after rebuild, certain `CSM specific steps` are required. We need to perform such action so we put a system back up health state.

---

## Master

#### Pre-condition

1. **NCN** is a **master** node

#### Actions

1. `scripts/k8s/update_kubeapi_istio_ca.sh`

---

## Worker

#### Pre-condition

1. **NCN** is a **worker** node

#### Actions

1. redeploy cps
1. `cfs/wait_for_configuration.sh`
1. ENSURE_KEY_PODS_HAVE_STARTED

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /kubernetes/{hostname}/pre-rebuild

#### POST
##### Summary

Kubernetes node pre rebuild action

##### Description

# K8s Pre Rebuild

Actions we need to perform before rebuild a k8s node

---

## Master

### Pre-condition

1. **NCN** is a **master** node
1. **NCN** is already the **first master**

### Action

1. Loop through other master nodes until `scripts/k8s/promote-initial-master.sh` returns 0
2. Update `meta-data.first-master-hostname`

---

## worker

### Pre-condition

1. **NCN** is a **worker** node

### Action

1. ENSURE_NEXUS_CAN_START_ON_ANY_NODE
1. ENSURE_ETCD_PODS_RUNNING
1. ENSURE_POSTGRES_HEALTHY
1. `cfs/wait_for_configuration.sh`
1. snapshot cps deployment

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/backup

#### POST
##### Summary

Create a NCN backup

##### Description

# NCN

Create backup of a ncn based on a predefined list so critical files can be restored after rebuild.

---

## Master

#### Pre-condition

1. **NCN** is a **master** node

#### Actions

1. backup local **sat** file
1. (m001 only) backup **ifcfg-lan0**
1. upload backup to s3

---

## Worker

#### Pre-condition

1. **NCN** is a **worker** node

#### Actions

1. bakcup ssh keys/authroized_keys
1. upload backup to s3

---

## Storage

1. **NCN** is a **ceph storage** node

#### Pre-condition

#### Actions

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/post-rebuild

#### POST
##### Summary

Perform post rebuild action on a NCN

##### Description

# NCN Post Rebuild

After a ncn has been rebuilt, some `CSM specific` steps are required.

---

## Master/Worker

#### Pre-condition

1. **NCN** is a **master** node

#### Actions

1. install latest docs-csm rpm
1. set `metal.no-wipe=1`

---

## Storage

#### Pre-condition

#### Actions

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/reboot

#### POST
##### Summary

Perform reboot on a NCN

##### Description

# NCN

Set to boot from pxe and power cycle the ncn

---

## Master/Worker/Storage

#### Pre-condition

#### Actions

1. Set boot to pxe
2. `ipmitool` power cycle the ncn

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/restore

#### POST
##### Summary

Restore a NCN backup

##### Description

# NCN

Restore previously backup files to a ncn.

---

## Master/Worker/Storage

#### Pre-condition

`N/A`

#### Actions

1. download backup from s3
1. untar/restore backup

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/validate

#### POST
##### Summary

Perform validation on a NCN

##### Description

# NCN

Run validation step of a ncn

---

## Master/Worker/Storage

#### Pre-condition

#### Actions

1. run goss test

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### /ncn/{hostname}/wipe

#### POST
##### Summary

Perform disk wipe on a NCN

##### Description

# NCN

Wipe a ncn's disk and set BSS `metal.no-wipe` to `0` so it actually gets wiped on boot

---

## Master

#### Pre-condition

1. **NCN** is a **master** node

#### Actions

1. Wipe disk

```
usb_device_path=$(lsblk -b -l -o TRAN,PATH \| awk /usb/'{print $2}')
usb_rc=$?
set -e
if [[ "$usb_rc" -eq 0 ]]; then
    if blkid -p $usb_device_path; then
    have_mnt=0
    for mnt_point in /mnt/rootfs /mnt/sqfs /mnt/livecd /mnt/pitdata; do
        if mountpoint $mnt_point; then
        have_mnt=1
        umount $mnt_point
        fi
    done
    if [ "$have_mnt" -eq 1 ]; then
        eject $usb_device_path
    fi
    fi
fi
umount /var/lib/etcd /var/lib/sdu \|\| true
for md in /dev/md/*; do mdadm -S $md \|\| echo nope ; done
vgremove -f --select 'vg_name=~metal*' \|\| true
pvremove /dev/md124 \|\| true
# Select the devices we care about; RAID, SATA, and NVME devices/handles (but *NOT* USB)
disk_list=$(lsblk -l -o SIZE,NAME,TYPE,TRAN \| grep -E '(raid\|sata\|nvme\|sas)' \| sort -u \| awk '{print "/dev/"$2}' \| tr '\\n' ' ')
for disk in $disk_list; do
    wipefs --all --force wipefs --all --force "$disk" \|\| true
    sgdisk --zap-all "$disk"
done
```

2. set `metal.no-wipe=0`

---

## Worker

1. **NCN** is a **worker** node

#### Actions

1. Wipe disk

```
lsblk \| grep -q /var/lib/sdu
sdu_rc=$?
vgs \| grep -q metal
vgs_rc=$?
set -e
systemctl disable kubelet.service \|\| true
systemctl stop kubelet.service \|\| true
systemctl disable containerd.service \|\| true
systemctl stop containerd.service \|\| true
umount /var/lib/containerd /var/lib/kubelet \|\| true
if [[ "$sdu_rc" -eq 0 ]]; then
    umount /var/lib/sdu \|\| true
fi
for md in /dev/md/*; do mdadm -S $md \|\| echo nope ; done
if [[ "$vgs_rc" -eq 0 ]]; then
    vgremove -f --select 'vg_name=~metal*' \|\| true
    pvremove /dev/md124 \|\| true
fi
wipefs --all --force /dev/sd* /dev/disk/by-label/* \|\| true
sgdisk --zap-all /dev/sd*
```

2. set `metal.no-wipe=0`

---

## Storage

#### Pre-condition

1. **NCN** is a **storage** node

#### Actions

1. Wipe disk

```
for d in $(lsblk \| grep -B2 -F md1 \| grep ^s \| awk '{print $1}'); do wipefs -af "/dev/$d"; done
```

2. set `metal.no-wipe=0`

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| hostname | path | Hostname | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 400 | Bad Request | [utils.ResponseError](#utilsresponseerror) |
| 404 | Not Found | [utils.ResponseError](#utilsresponseerror) |
| 500 | Internal Server Error | [utils.ResponseError](#utilsresponseerror) |

### Models

#### utils.ResponseError

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string |  | No |

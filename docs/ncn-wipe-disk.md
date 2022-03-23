## NCN wipe disk

Wipe a ncn's disk and set BSS `metal.no-wipe` to `0` so it actually gets wiped on boot

---

### Master

#### Pre-condition

1. **NCN** is a **master** node

#### Actions

1. Wipe disk

```
usb_device_path=$(lsblk -b -l -o TRAN,PATH | awk /usb/'{print $2}')
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
umount /var/lib/etcd /var/lib/sdu || true
for md in /dev/md/*; do mdadm -S $md || echo nope ; done
vgremove -f --select 'vg_name=~metal*' || true
pvremove /dev/md124 || true
# Select the devices we care about; RAID, SATA, and NVME devices/handles (but *NOT* USB)
disk_list=$(lsblk -l -o SIZE,NAME,TYPE,TRAN | grep -E '(raid|sata|nvme|sas)' | sort -u | awk '{print "/dev/"$2}' | tr '\n' ' ')
for disk in $disk_list; do
    wipefs --all --force wipefs --all --force "$disk" || true
    sgdisk --zap-all "$disk"
done
```

2. set `metal.no-wipe=0`

#### Microservices

| name | protocol/client | credentials | Note |
| ---- | --------------- | ----------- | ---- |
| wipe | ssh as root     | k8s secret  |      |
| bss  | bss go client   | jwt token   |      |

---

### Worker

1. **NCN** is a **worker** node

#### Actions

1. Wipe disk

```
lsblk | grep -q /var/lib/sdu
sdu_rc=$?
vgs | grep -q metal
vgs_rc=$?
set -e
systemctl disable kubelet.service || true
systemctl stop kubelet.service || true
systemctl disable containerd.service || true
systemctl stop containerd.service || true
umount /var/lib/containerd /var/lib/kubelet || true
if [[ "$sdu_rc" -eq 0 ]]; then
    umount /var/lib/sdu || true
fi
for md in /dev/md/*; do mdadm -S $md || echo nope ; done
if [[ "$vgs_rc" -eq 0 ]]; then
    vgremove -f --select 'vg_name=~metal*' || true
    pvremove /dev/md124 || true
fi
wipefs --all --force /dev/sd* /dev/disk/by-label/* || true
sgdisk --zap-all /dev/sd*
```

2. set `metal.no-wipe=0`

#### Microservices

| name | protocol/client | credentials | Note |
| ---- | --------------- | ----------- | ---- |
| wipe | ssh as root     | k8s secret  |      |
| bss  | bss go client   | jwt token   |      |

---

### Storage

#### Pre-condition

1. **NCN** is a **storage** node

#### Actions

1. Wipe disk

```
for d in $(lsblk | grep -B2 -F md1 | grep ^s | awk '{print $1}'); do wipefs -af "/dev/$d"; done
```

2. set `metal.no-wipe=0`

#### Microservices

| name | protocol/client | credentials | Note |
| ---- | --------------- | ----------- | ---- |
| wipe | ssh as root     | k8s secret  |      |
| bss  | bss go client   | jwt token   |      |

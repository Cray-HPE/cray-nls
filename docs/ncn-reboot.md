## NCN Reboot

Set to boot from pxe and power cycle the ncn

---

### Master/Worker/Storage

##### Pre-condition

#### Actions

1. Set boot to pxe

2. `ipmitool` power cycle the ncn

#### Microservices

| name         | protocol/client | credentials | Note |
| ------------ | --------------- | ----------- | ---- |
| set pxe boot | ipmi            | k8s secret  |      |
| power cycle  | ipmi            | k8s secret  |      |

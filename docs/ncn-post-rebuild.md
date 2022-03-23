## NCN Post Rebuild

After a ncn has been rebuilt, some `CSM specific` steps are required.

---

### Master/Worker

##### Pre-condition

1. **NCN** is a **master** node

#### Actions

1. install latest docs-csm rpm

1. set `metal.no-wipe=1`

#### Microservices

| name            | protocol/client | credentials | Note                                                                                |
| --------------- | --------------- | ----------- | ----------------------------------------------------------------------------------- |
| install doc rpm | ssh as root     | k8s secret  | we should look into bss/cloud-init so it always install what we specify during boot |
| set no wipe     | bss client      | jwt token   |                                                                                     |

---

### Storage

##### Pre-condition

#### Actions

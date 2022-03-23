## K8s Pre Rebuild

Actions we need to perform before rebuild a k8s node

---

### Master

#### Pre-condition

1. **NCN** is a **master** node

1. **NCN** is already the **first master**

#### Action

1. Loop through other master nodes until `scripts/k8s/promote-initial-master.sh` returns 0

2. Update `meta-data.first-master-hostname`

#### Microservices

| name              | protocol/client | credentials | Note                                                                                |
| ----------------- | --------------- | ----------- | ----------------------------------------------------------------------------------- |
| move first master | ssh as root     | k8s secret  | we need to look into the script and figure out exactly which microservices it calls |
| bss               | bss go client   | jwt token   |                                                                                     |

---

### worker

#### Pre-condition

1. **NCN** is a **worker** node

#### Action

1. ENSURE_NEXUS_CAN_START_ON_ANY_NODE

1. ENSURE_ETCD_PODS_RUNNING

1. ENSURE_POSTGRES_HEALTHY

1. `cfs/wait_for_configuration.sh`

1. snapshot cps deployment

#### Microservices

| name                         | protocol/client | credentials | Note                                                                  |
| ---------------------------- | --------------- | ----------- | --------------------------------------------------------------------- |
| ensure some pods are running | ssh as root     | k8s secret  | All can be done by using k8s client                                   |
| ensure pg healthy            | ssh as root     | k8s secret  | All can be done by using k8s client                                   |
| wait for cfs                 | ssh as root     | k8s secret  | All can be done by using k8s client (?)                               |
| snapshot cps deployment      | ssh as root     | k8s secret  | is `cray cps` an api call? if so we can make api calls instead of ssh |

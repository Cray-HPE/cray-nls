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

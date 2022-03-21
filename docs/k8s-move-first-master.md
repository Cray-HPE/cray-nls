# Move First Master

We need to make sure first master is not the node being rebuit. We need to move `first_master` to a different master node

### Pre-condition

1. **NCN** is a **master** node
1. **NCN** is already the **first master**

### Action

1. Loop through other master nodes until `scripts/k8s/promote-initial-master.sh` returns 0
2. Update `meta-data.first-master-hostname`
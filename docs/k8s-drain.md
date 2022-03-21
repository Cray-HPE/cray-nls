# Drain Kubernetes Node

Before we can safely drain/remove a node from k8s cluster, we need to run some `CSM specific logic` to make sure a node can be drained from k8s cluster safely

---

## Master

#### Pre-condition

1. **NCN** is a **master** node
1. quorum after removal

#### Actions

1. drain node

---

## Worker

#### Pre-condition

#### Actions

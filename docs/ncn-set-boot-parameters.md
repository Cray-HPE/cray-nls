## NCN set boot parameters

After a node rejoined k8s cluster after rebuild, certain `CSM specific steps` are required. We need to perform such action so we put a system back up health state.

---

#### Actions
1. update cloud-init global data
1. set which image to boot
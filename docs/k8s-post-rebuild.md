## K8s Post Rebuild

After a node rejoined k8s cluster after rebuild, certain `CSM specific steps` are required. We need to perform such action so we put a system back up health state.

---

### Master

#### Pre-condition

1. **NCN** is a **master** node

#### Actions

1. `scripts/k8s/update_kubeapi_istio_ca.sh`

#### Microservices

| name                    | protocol/client | credentials | Note |
| ----------------------- | --------------- | ----------- | ---- |
| update_kubeapi_istio_ca | ssh as root     | k8s secret  |      |

---

### Worker

#### Pre-condition

1. **NCN** is a **worker** node

#### Actions

1. redeploy cps

1. `cfs/wait_for_configuration.sh`

1. ENSURE_KEY_PODS_HAVE_STARTED

#### Microservices

| name                        | protocol/client | credentials | Note                                                                  |
| --------------------------- | --------------- | ----------- | --------------------------------------------------------------------- |
| cps redeploy                | ssh as root     | k8s secret  | is `cray cps` an api call? if so we can make api calls instead of ssh |
| wait for cfs                | ssh as root     | k8s secret  | All can be done by using k8s client (?)                               |
| ensure key pods are running | ssh as root     | k8s secret  | All can be done by using k8s client                                   |

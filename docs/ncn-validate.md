## NCN Validation

Run validation step of a ncn

---

### Master/Worker/Storage

#### Pre-condition

#### Actions

1. run goss test

#### Microservices

| name          | protocol/client | credentials | Note                                                  |
| ------------- | --------------- | ----------- | ----------------------------------------------------- |
| run goss test | ssh as root     | k8s secret  | goss has a server that accepts REST call to run tests |

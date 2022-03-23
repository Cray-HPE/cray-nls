## NCN restore backup

Restore previously backup files to a ncn.

---

### Master/Worker/Storage

##### Pre-condition

`N/A`

#### Actions

1. download backup from s3

1. untar/restore backup

#### Microservices

| name             | protocol/client | credentials | Note |
| ---------------- | --------------- | ----------- | ---- |
| download from s3 | s3 client       | jwt token   |      |
| restore backup   | ssh as root     | k8s secret  |      |

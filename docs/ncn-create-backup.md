# NCN

Create backup of a ncn based on a predefined list so critical files can be restored after rebuild.

---

## Master

#### Pre-condition

1. **NCN** is a **master** node

#### Actions

1. backup local **sat** file
1. (m001 only) backup **ifcfg-lan0**
1. upload backup to s3

---

## Worker

#### Pre-condition

1. **NCN** is a **worker** node

#### Actions

1. bakcup ssh keys/authroized_keys
1. upload backup to s3

---

## Storage

1. **NCN** is a **ceph storage** node

#### Pre-condition

#### Actions

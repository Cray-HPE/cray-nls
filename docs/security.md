## Security

### Authentication

Similar to other exposed services, authentication is done by keycloak. Keycloak issued jwt token is verified and passed to API gateway.

### Routes/AuthZ

Each route of these APIs are protected by configuring OPA policy.

- **Crawl Phase**

  we will only `admin` and `user` roles. Users have `admin` role are allowed to invoke any APIs. Users with `user` role will only be able to call **GET** APIs.

- **Walk Phase**

  we can introduce more granular permissions/roles based on future requirements.

- **Run Phase**

  we can even go to resources level. For example, `User A` could have all permissions of `ceph nodes` but this user won't be able to rebuild/reboot any k8s nodes. `Monitoring User` can rerun/retry any failed rebuild/reboots but can't initiate such operation.

### Microservices

The jwt token will be passed down to each microservices and individual microservice should enforce authZ in its own domain. Any credentials needed by each microservice should be obtained in a secure manner. SSH as root should be avoided if possible. However, there are certain operations requires root access via ssh. In those cases, we should use Vault to generate one time, short lived temporary SSH keys. Note that these goals will be achieved phase by phase.

- **Crawl Phase**

  In crawl phase, we execute steps almost identical to what we have today. Most steps need direct root access via SSH. SSH credentials are mounted onto each short lived _Job Pods_ as `hostPath`. JWT tokens needed for other microservice calls are obtained from `ncn-m001` over SSH:

  ```
  export TOKEN=$(curl -k -s -S -d grant_type=client_credentials \
   -d client_id=admin-client \
   -d client_secret=`kubectl get secrets admin-client-auth -o jsonpath='{.data.client-secret}' | base64 -d` \
   https://api-gw-service-nmn.local/keycloak/realms/shasta/protocol/openid-connect/token | jq -r '.access_token')
  ```

  > NOTE: This is exactly what our 1.0.x and 1.2.x does

- **Walk Phase**

  - SSH credentials need to be controlled by Vault and only one time credentials should be used
  - JWT token should be passed from API gateway instead of getting it from `ncn-m001` as root user
  - Any steps can be performed by make REST/gRPC request to a microservice should not use SSH any more

- **Run Phase**

  Each microservice should implement it's own granular/resources level authZ

### Logging/Audit

- **Request info**

  API Gateway should log user information from validated JWT token so we know "who did what at when". Each microservice should also log the same information. Additionally, a unique request id should be passed/logged as well such that we can track a request in every microservice. Note that this is slightly different than what istio tracking is because of async operations. It won't carry istio injected `x-b3-traceid` in some cases.

  Required fields:

  - User Info: `name/id/email`, `roles`
  - HTTP path: `REST API URI`
  - HTTP method: `GET|POST|PUT|DELETE`
  - Resources list: `ncn-w001,ncn-w002...`
  - Operation Result: `failed|succeed|terminated`

- **Operation logs**

  Each steps of automation should be logged in order to troubleshoot/audit what exactly happened on ncn(s). This is done by _Argo Workflow_ engine.
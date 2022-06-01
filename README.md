# Development Setup

### Prereq

1. [Go 1.17](https://go.dev/doc/install)
1. [Go Fresh](https://github.com/gravityblast/fresh)
   ```
   go get github.com/pilu/fresh
   ```
1. [K3D](https://k3d.io/) (required for running locally)
1. [NodeJS](https://nodejs.org/en/download/) (required for markdown version of swagger doc)

### Start

1. Run argo workflow in k3d

   ```
   scripts/argo.local.sh
   ```

   This will start a k3d cluster and deploy minimal Argo Workflow. It also port-forward `2746` to localhost


1. Modify `cluster-admin` file

    > Run the following command

    ```bash
    kubectl get nodes
    ```
    If you have 1 master and 3 worker nodes running you can proceed, otherwise re-try *step 1: Run argo workflow in k3d*


    > open the `cluster-admin` file in your default editor by running:

    ```
     kubectl edit ClusterRoleBindings/cluster-admin
    ```

    Replace the entire `subjects:` section with the following block

    ```yaml
    subjects:
      - apiGroup: rbac.authorization.k8s.io
        kind: Group
        name: system:masters
      - kind: ServiceAccount
        name: default
        namespace: argo
    ```


    Save and exit


1. Create and update `.env` file

     > Copy .env.example file to .env file

    ```bash
    cp .env.example .env
    ```

    Then replace the contents of the `.env` with the following block:

    ```
    SERVER_PORT=3000
    ENV=development
    ARGO_TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6Ijk4dWFiYTF4ZU56SFA2OGFETWtwVXNlUGJhekdfV3B4NE9zYnpOVWZFZlkifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJhcmdvIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImFyZ28tc2VydmVyLXRva2VuLThzMmZkIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImFyZ28tc2VydmVyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiNDRkMjZkYzUtM2JhNi00Nzc1LTk3NjgtYWMzYmI4ZDUwNjEwIiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50OmFyZ286YXJnby1zZXJ2ZXIifQ.nGIrVeWodcv3TdvP9A5FpAH1mjUUAzdVyiRB6zZ65Nd_qRsAgKMJaDoasIEaCTG3O1NqNH-L6EBmK4_wMkx5hCV4zgShkWyDwmQBGccL52-0g9r2EPzZUrK2djRoQFNBAVHhYdL8kcUpNxl-nKLBrLCNyinGnHRWblRmt021nZZsE62ljXj1TDmNRbg71oC94tnBeGi4j3Iza9KBf8cTjTeKYVyDLQrd6gGQuUpWZ9l4vfK7X4Ke8dDEuWtWkxDZgnNewQLZmPP7K-QUZEbNlL82r4KJSExaCJmBZfPx21QH0bSJbrl6Xr1m22W1OqTCEQEsNv2uIm7nV4A2vmhdZA
    ARGO_SERVER_URL=localhost:2746
    ```



1. Run server

   > One time setup: copy and rename `.env.example` to `.env`

   ```
   scripts/runDev.sh
   ```

   Automatically rebuild and launch API server when a change is made. Fresh configuration file: `{rootDir}/runner.conf`

1. Load swagger and Argo

    > These sites are on localhost and will only open with the previous scripts running

    - [Swagger URL](http://localhost:3000/swagger/index.html)
    - [Argo URL](http://localhost:2746/)


### Optional Steps

- Update swagger
   ```
   scripts/swagger.gen.sh
   ```
   > Note: This script will try to update `docs/swagger.md` if nodejs is installed. Otherwise, it will only update `docs/swagger.yaml`



# Reference

[Dependency Injection](https://medium.com/swlh/dependency-injection-in-go-using-fx-6a623c5c5e01)

- [uber fx](https://github.com/uber-go/fx)

[Clean gin template](https://github.com/dipeshdulal/clean-gin)

[Argo Workflow](https://argoproj.github.io/argo-workflows)

[K3D](https://k3d.io/)

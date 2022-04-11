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

1. Run server

   > One time setup: copy and rename `.env.example` to `.env`

   ```
   scripts/runDev.sh
   ```

   Automatically rebuild and launch API server when a change is made. Fresh configuration file: `{rootDir}/runner.conf`

1. Update swagger
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

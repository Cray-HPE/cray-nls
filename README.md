# Development

### Prereq

1. golang 1.17
2. npm: required for swagger doc gen

### Development

1. Run argo workflow in k3d

   ```
   scripts/argo.local.sh
   ```

1. Run server
   ```
   go run main.go
   ```
1. Update swagger
   ```
   scripts/swagger.gen.sh
   ```

# Reference

[clean gin template](https://github.com/dipeshdulal/clean-gin)

[Argo Workflow](https://argoproj.github.io/argo-workflows)

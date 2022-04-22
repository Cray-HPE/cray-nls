# Temporary script to deploy argo onto a system for testing
#     NOTE: this will be removed once we have deployment work done

# csm 1.2
NEXUS_USERNAME="$(kubectl -n nexus get secret nexus-admin-credential --template {{.data.username}} | base64 -d)"
NEXUS_PASSWORD="$(kubectl -n nexus get secret nexus-admin-credential --template {{.data.password}} | base64 -d)"
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-username "$NEXUS_USERNAME" --dest-password "$NEXUS_PASSWORD"  docker://quay.io/argoproj/argocli:latest docker://registry.local/argoproj/argocli:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-username "$NEXUS_USERNAME" --dest-password "$NEXUS_PASSWORD"  docker://quay.io/argoproj/argoexec:latest docker://registry.local/argoproj/argoexec:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-username "$NEXUS_USERNAME" --dest-password "$NEXUS_PASSWORD"  docker://quay.io/argoproj/workflow-controller:latest docker://registry.local/argoproj/workflow-controller:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-username "$NEXUS_USERNAME" --dest-password "$NEXUS_PASSWORD"  docker://postgres:12-alpine docker://registry.local/library/postgres:12-alpine
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-username "$NEXUS_USERNAME" --dest-password "$NEXUS_PASSWORD"  docker://minio/minio docker://registry.local/docker.io/minio/minio


# csm 1.0
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://quay.io/argoproj/argocli:latest docker://registry.local/argoproj/argocli:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://quay.io/argoproj/argoexec:latest docker://registry.local/argoproj/argoexec:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://quay.io/argoproj/workflow-controller:latest docker://registry.local/argoproj/workflow-controller:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://postgres:12-alpine docker://registry.local/library/postgres:12-alpine
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://minio/minio docker://registry.local/library/minio/minio
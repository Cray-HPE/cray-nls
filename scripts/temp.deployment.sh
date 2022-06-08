#
# MIT License
#
# (C) Copyright 2022 Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.
#
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
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-username "$NEXUS_USERNAME" --dest-password "$NEXUS_PASSWORD"  docker://kennethreitz/httpbin:latest docker://registry.local/docker.io/kennethreitz/httpbin:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-username "$NEXUS_USERNAME" --dest-password "$NEXUS_PASSWORD"  docker://portainer/kubectl-shell:latest-v1.21.1-amd64 docker://registry.local/docker.io/portainer/kubectl-shell:latest-v1.21.1-amd64


# csm 1.0
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://quay.io/argoproj/argocli:latest docker://registry.local/argoproj/argocli:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://quay.io/argoproj/argoexec:latest docker://registry.local/argoproj/argoexec:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://quay.io/argoproj/workflow-controller:latest docker://registry.local/argoproj/workflow-controller:latest
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://postgres:12-alpine docker://registry.local/library/postgres:12-alpine
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-creds "admin:admin" docker://minio/minio docker://registry.local/library/minio/minio
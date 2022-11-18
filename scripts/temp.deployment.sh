#!/bin/bash
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
set -e

function deployNLS() {
    BUILDDIR="/tmp/build"
    mkdir -p "$BUILDDIR"
    kubectl get secrets -n loftsman site-init -o jsonpath='{.data.customizations\.yaml}' | base64 -d > "${BUILDDIR}/customizations.yaml"
    kubectl get configmap -n loftsman loftsman-platform -o jsonpath='{.data.manifest\.yaml}' > "${BUILDDIR}/iuf.yaml"
    manifestgen -i "${BUILDDIR}/iuf.yaml" -c "${BUILDDIR}/customizations.yaml" -o "${BUILDDIR}/platform.yaml"
    yq w -i "${BUILDDIR}/platform.yaml" 'spec.charts[0].version' "$2"
    yq w -i "${BUILDDIR}/platform.yaml" 'metadata.name' "iuf"
    charts="$(yq r "${BUILDDIR}/platform.yaml" 'spec.charts[*].name')"
    for chart in $charts; do
        if [[ $chart != "cray-iuf" ]] && [[ $chart != "cray-nls" ]]; then
            yq d -i "${BUILDDIR}/platform.yaml" "spec.charts.(name==$chart)"
        fi
    done

    yq d -i "${BUILDDIR}/platform.yaml" "spec.sources"

    loftsman ship --charts-path "$1" --manifest-path /tmp/build/platform.yaml
}
CHART_PATH="/tmp/nls"
rm -rf ${CHART_PATH}
git clone https://github.com/Cray-HPE/cray-nls.git ${CHART_PATH}
if [[ -n $1 ]]; then
    cd ${CHART_PATH} && git checkout "$1"
fi
cd ${CHART_PATH} && helm dep update charts/v1.0/cray-nls/
cd ${CHART_PATH} && helm package charts/v1.0/cray-nls/
cd ${CHART_PATH} && helm dep update charts/v1.0/cray-iuf/
cd ${CHART_PATH} && helm package charts/v1.0/cray-iuf/


echo "Get NLS chart version"
tarFileName=$(ls | cat | grep "cray-nls")
tarFileName=${tarFileName#"cray-nls-"}
version=${tarFileName%".tgz"}
echo "Version: $version"

echo "Get image version"
imageVersion=$(yq r ${CHART_PATH}/charts/v1.0/cray-nls/values.yaml 'cray-service.containers.cray-nls.image.tag')
echo "image version: $imageVersion"
echo "sync docker images"
NEXUS_USERNAME="$(kubectl -n nexus get secret nexus-admin-credential --template {{.data.username}} | base64 -d)"
NEXUS_PASSWORD="$(kubectl -n nexus get secret nexus-admin-credential --template {{.data.password}} | base64 -d)"
podman run --rm --network host quay.io/skopeo/stable copy --src-tls-verify=false --dest-tls-verify=false --dest-username "$NEXUS_USERNAME" --dest-password "$NEXUS_PASSWORD"  docker://artifactory.algol60.net/csm-docker/stable/cray-nls:$imageVersion docker://registry.local/artifactory.algol60.net/csm-docker/stable/cray-nls:$imageVersion

deployNLS "$CHART_PATH" "$version"
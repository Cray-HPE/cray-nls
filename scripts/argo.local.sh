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

if ! command -v k3d &> /dev/null
then
    wget -q -O - https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
fi

k3d cluster list  | grep "mycluster"

if [[ $? -ne 0 ]]; then
    k3d cluster create mycluster \
        -a 2 \
        --agents-memory 1g \
        --servers-memory 1g \
        --no-lb \
        --k3s-arg "--node-name=ncn-w001"@agent:0 \
        --k3s-arg "--node-name=ncn-w002"@agent:1 \
        --k3s-arg "--node-name=ncn-m001"@server:0
else 
    k3d cluster start --wait
fi
k3d kubeconfig merge mycluster --kubeconfig-switch-context
kubectl wait --for=condition=ready nodes -l node.kubernetes.io/instance-type=k3s
kubectl get nodes

kubectl create ns argo
kubectl apply -n argo -f https://raw.githubusercontent.com/argoproj/argo-workflows/master/manifests/quick-start-minimal.yaml
kubectl wait --for=condition=ready pod -l app=argo-server -n argo
kubectl -n argo port-forward svc/argo-server 2746:2746
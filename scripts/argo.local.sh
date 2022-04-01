
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
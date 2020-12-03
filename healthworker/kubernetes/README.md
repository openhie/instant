# Kubernetes

To start use the following:
```bash
kubectl apply -k dev/
# to remove
kubectl delete -k dev/
```

## Troubleshooting

[`kompose`](https://kompose.io) was used to create the resource manifests. There are two issues that this process generates. 

One is that a port is not open for the service. Make the following change.
```yaml
spec:
  type: LoadBalancer
  ports:
  - name: "3000"
    port: 3000
    targetPort: 3000
  selector:
    io.kompose.service: facility-recon
# status:
#   loadBalancer: {}
```

Another is that the redis-service is not created. This is easily copied from another place.
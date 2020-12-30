# Kubernetes

To start use the following:
```bash
kubectl apply -k main/
# to remove
kubectl delete -k main/
```

## Troubleshooting

[`kompose`](https://kompose.io) was used to create the resource manifests. There are two issues that this process generates.  One is that a port is not open for the service. 

Another is that there is potential for port contention, so the exposed port is remapped to 3001.

If rebuilding the Kompose-generated manifests, make the following changes:
```yaml
spec:
  type: LoadBalancer
  # change default port to avoid contention
  ports:
    - name: "3001"
      port: 3001
      targetPort: 3000
  selector:
    io.kompose.service: ihris
# status:
#   loadBalancer: {}
```

Another issue is that the redis-service is not created. This is easily copied from another set of manifests.
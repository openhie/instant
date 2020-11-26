# Kubernetes


To start use the following:

```bash
kubectl apply -k dev/
# to remove
kubectl delete -k dev/
```

The OpenCR service may not become available for one or two minutes waiting for HAPI FHIR to fully start.

## Troubleshooting

[`kompose`](https://kompose.io) was used to create the resource manifests. There is one issue that this process generates which means that a port is not open for the service. Make the following change.
```yaml
spec:
  type: LoadBalancer
  ports:
    - name: "3000"
      port: 3000
      targetPort: 3000
  selector:
    io.kompose.service: opencr
# status:
#   loadBalancer: {}
```
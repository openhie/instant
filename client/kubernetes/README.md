# Kubernetes

To start use the following:

```bash
kubectl apply -k main/
# to remove
kubectl delete -k main/
```

The OpenCR service may not become available for one or two minutes waiting for HAPI FHIR to fully start.

## Troubleshooting

[`kompose`](https://kompose.io) was used to create the resource manifests. There is an issue that this process generates which means that a port is not open for the service. Make the following change.

Another is that there is potential for port contention, so the exposed port is remapped to 3003.

```yaml
spec:
  type: LoadBalancer
  # change to 3004 from default to avoid port contention
  ports:
    - name: "3004"
      port: 3004
      targetPort: 3000
  selector:
    io.kompose.service: opencr
```

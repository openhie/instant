# Setup

## Kubernetes

Navigate to the folder `/kubernetes/main` to execute the commands.

To set up the healthworkforce services run the following command:

```sh
./k8s.sh up
```

To tear down the deployments use the opposing command:

```bash
./k8s.sh down
```

To completely remove all project components use the following option:

```bash
./k8s.sh destroy
```

To configure the services set up, navigate to the folder `/kubernetes/importer` and run the following:

```bash
./k8s.sh up
```

To clean up the services for configuration run the following:

```bash
./k8s.sh clean
```

## Docker

Navigate to the folder `/Docker` to execute the commands.

To set up the healthworkforce services run the following command:

```sh
./compose.sh init
```

To tear down the deployments use the opposing command:

```bash
./compose.sh down
```

To start up the services after a tear down, use the following command:

```bash
./compose.sh up
```

To completely remove all project components use the following option:

```bash
./compose.sh destroy
```

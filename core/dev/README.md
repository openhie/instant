# Instant OpenHIE development instance setup

A development instance of the instant openHIE can be setup by running the command below in this directory

```bash
  ./setup.sh
```

To configure the openhim console to communicate with the openhim core, you need the `host` (external ip address of the container) and api `port` (port that maps to 8080) of the core. These can be retrieved by running the following command

```bash
  kubectl get services core
```

Update the openhim console configuration (ie substitute the port and host) by running the following command

```bash
  ./configure-console.sh
```

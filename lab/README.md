

### Run with Docker-compose

Navigate to the main folder to execute the commands.

To set the Instant OpenHIE services run the following command:

```sh
yarn
yarn docker:build
```
To Start the Core and Lab Services 

```sh
yarn docker:instant init -t docker core lab
```

To tear down the deployments use the opposing command:

```bash
yarn docker:instant down -t docker
```

To start up the services after a tear down, use the following command:

```bash
yarn docker:instant up -t docker
```

To completely remove all project components use the following option:

```bash
yarn docker:instant destroy -t docker
```

## Accessing the Laboratory service

### OpenELIS

* URL: https://localhost:8443/OpenELIS-Global
* Username: **admin**
* Password: **adminADMIN!**

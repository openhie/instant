# Instant OpenHIE

> Note: This repo is not for production. Rather, it contains strawpersons to facilitate discussion and demonstrations of progress on the Instant OpenHIE project. As such, it implies no endorsement or support from any institution especially and including the OpenHIE Community of Practice. This is for open discussion in the community. Please join the OpenHIE Dev-Ops Sub-community and give us your thoughts!

```sh
docker run -it --rm -v /var/run/docker.sock:/var/run/docker.sock -v /usr/bin/docker:/usr/bin/docker -v /usr/local/bin/docker-compose:/usr/local/bin/docker-compose -v /usr/bin/docker:/usr/bin/docker -v /usr/lib/libdevmapper.so.1.02:/usr/lib/libdevmapper.so.1.02 instant <cmd>
```

```sh
docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v /usr/bin/docker:/usr/bin/docker -v /usr/local/bin/docker-compose:/usr/local/bin/docker-compose -v /usr/bin/docker:/usr/bin/docker -v /usr/lib/libdevmapper.so.1.02:/usr/lib/libdevmapper.so.1.02 -v ~/.kube/config:/root/.kube/config -v ~/.minikube:/home/ryan/.minikube -v ~/.minikube:/root/.minikube -v /usr/local/bin/minikube:/usr/local/bin/minikube --network host instant up k8s
```
# Health Worker Infrastructure Package

The Health Worker Infrastructure Package deploys iHRIS v5 as a prototypical reference health worker registry and interface to the same.

This package consists of several services:

- [iHRIS](https://ihris.org)
- [HAPI FHIR Server](https://hapifhir.io/)
- [ElasticSearch + Kibana](https://elastic.co)
- [Redis] in-memory datastore

The package has strict requirements. iHRIS data, forms, and internals are based on FHIR structure definitions. Whether in Kubernetes or Docker, the containers must be run before iHRIS is launched as they load structure definitions in HAPI. 

echo -e "\nHAPI FHIR test..."
curl https://openhim-core.ssl.instant/hapi-fhir-jpaserver/fhir/Patient -k -H "Content-Type: application/fhir+json" -f -d '{ resourceType: "Patient" }' && echo -e '\n\nSUCCESS' || (echo -e '\nFAILED'; exit 1)

echo -e "\nHAPI FHIR test..."

HOST=$1

curl https://$HOST/hapi-fhir-jpaserver/fhir/Patient -k -H "Content-Type: application/fhir+json" -f -d '{ resourceType: "Patient" }' && echo -e '\n\nSUCCESS' || (echo -e '\nFAILED'; exit 1)

echo -e "\nHAPI FHIR test..."

URLBASE=$1

curl $URLBASE/fhir/Patient -k -H "Content-Type: application/fhir+json" -f -d '{ "resourceType": "Patient" }' && echo -e '\n\nSUCCESS' || (echo -e '\nFAILED'; exit 1)

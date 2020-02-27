echo -e "\nMapping medaitor test..."
curl https://openhim-core.ssl.instant/dhis2part1 -X POST -k -H "Content-Type: application/fhir+json" -f && echo -e '\n\nSUCCESS' || (echo -e '\nFAILED'; exit 1)

echo -e "\nMapping mediator test..."

HOST=$1

curl https://$HOST/dhis2part1 -X POST -k -H "Content-Type: application/json" -f && echo -e '\n\nSUCCESS' || (echo -e '\nFAILED'; exit 1)

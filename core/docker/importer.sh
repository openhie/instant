
curl -X POST --insecure -u "$OPENHIM_API_USERNAME:$OPENHIM_API_PASSWORD" -H "Content-Type: application/json" -d @openhim-import.json "$OPENHIM_API_SERVER/metadata" -v
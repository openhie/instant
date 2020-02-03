#!/bin/sh

OPENHIM_RESPONSE="";
TIMEOUT=0

while [ "$OPENHIM_RESPONSE" != "200" ]
do
  TIMEOUT=$((TIMEOUT+1))
  # Fail script if incomplete at 30 seconds
  if [ "$TIMEOUT" == "15" ]
  then
    echo -e "Sleep timeout reached: $TIMEOUT.\nFailed to import config. Could not reach the OpenHIM."
    exit 1
  fi

  echo "OpenHIM not ready ( $OPENHIM_RESPONSE ) - sleeping"
  sleep 2
  OPENHIM_RESPONSE=$(curl -X GET --insecure --write-out %{http_code} --silent --output /dev/null "$OPENHIM_API_SERVER/heartbeat");
done

echo -e "\nSTART Posting OpenHIM Config\n----------------------------\n"
curl --insecure -u "$OPENHIM_API_USERNAME:$OPENHIM_API_PASSWORD" -H "Content-Type: application/json" -d @openhim-import.json "$OPENHIM_API_SERVER/metadata" -v
echo -e "\n\n--------------------------\nEND Posting OpenHIM Config\n"

#!/bin/bash

: "${FIOS_SESSION:?Variable FIOS_SESSION not set or empty}"
: "${FIOS_XSRF:?Variable FIOS_XSRF not set or empty}"


curl -k \
    --cookie "test; Session=$FIOS_SESSION; XSRF-TOKEN=$FIOS_XSRF; bhr4HasEnteredAdvanced=false; bhr4UI2HasToRefresh=false" \
    -H "X-XSRF-TOKEN: $FIOS_XSRF" \
    https://192.168.1.1/api/dhcp/clients |\
  jq '.[] | {
  name: .name, 
  id: .id,
  ip: .ipAddress, 
  mac: .mac, 
  staticIp: .staticIp, 
}'

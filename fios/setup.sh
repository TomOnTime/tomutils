#!/bin/bash

set -e

# I couldn't get this to work.  Maybe the value of "password" (which
# you have to get from F12) isn't right?

#export FIOS_OUT=$(curl -k -s -i -d '{"password":"$1"}' https://192.168.1.1/api/login) 
#export FIOS_SESSION=$(echo "$FIOS_OUT" | grep "^Set-Cookie: Session=" | grep -oE '[0-9]+') 
#export FIOS_XSRF=$(echo "$FIOS_OUT" | grep "^Set-Cookie: XSRF-TOKEN" | grep -oE '=[0-9a-z]+;' | tr -d "=" | tr -d ";")

echo export FIOS_SESSION=\'"$FIOS_SESSION"\'
echo export FIOS_XSRF=\'"$FIOS_XSRF"\'

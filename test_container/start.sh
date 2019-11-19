#!/usr/bin/env bash

/root/start.sh &
echo "Waiting for service startup"
#until $(curl -o /dev/null -s -w "%{http_code}\n" -k https://omdadmin:omdadmin@localhost/demo/thruk/r/config/objects) eq "200"; do
#  echo "STARTING..."
#  sleep 1;
#done;

timeout 30 bash -c 'while [[ "$(curl --insecure -s -o /dev/null -w ''%{http_code}'' https://omdadmin:omdadmin@localhost/demo/thruk/r/config/objects)" != "200" ]]; do sleep 1; done'
if [ $? -eq 0 ]
then
   echo "Thruk server is up and running"
fi
wait
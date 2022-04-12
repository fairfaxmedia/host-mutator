#!/usr/bin/env sh

export ACTION=${ACTION:-host-mutator}

echo "ACTION: ${ACTION}"
echo "CWD:    $(pwd)"

if [ ${ACTION} = "host-mutator" ]; then
  /usr/local/bin/host-mutator
elif [ ${ACTION} = "certificates" ]; then
  /usr/local/bin/certificates
fi

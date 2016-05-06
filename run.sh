#!/bin/bash -e

if [ -z "$1" ]; then
	export CONFIG="artemis.conf"
else
	export CONFIG=$1
fi

./bin/artemisd -config $CONFIG
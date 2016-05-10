#!/bin/bash -e

scp ./bin/artemisd root@188.166.133.162:/usr/bin
scp ./artemis.conf.live root@188.166.133.162:/etc/artemis/artemis.conf
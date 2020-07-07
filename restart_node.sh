#!/bin/bash

PRONAME="node"

BIN=/usr/local/src/guard/node
STDLOG=/usr/local/src/guard/output.log


chmod u+x $BIN

ID=$(/usr/sbin/pidof "$BIN")
if [ "$ID" ] ; then
        echo "kill -SIGINT $ID"
        kill -2 $ID
fi

while :
do
        ID=$(/usr/sbin/pidof "$BIN")
        if [ "$ID" ] ; then
                echo "$PRONAME still running...wait"
                sleep 0.1
        else
                echo "$PRONAME service was not started"
                echo "Starting service..."

                nohup $BIN > $STDLOG 2>&1 &
                break
        fi
done

ps -aux | grep $BIN
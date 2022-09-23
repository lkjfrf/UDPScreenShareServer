#!/bin/sh

if [ "$1" = "start" ]; then
	sudo nohup ./main &
elif [ "$1" = "stop" ]; then
	pkill -9 -ef ./main
elif [ "$1" = "status" ]; then
	ps -ef | grep ./main
elif [ "$1" = "log" ]; then
	tail -f Screen_Log.log
fi

exit 0

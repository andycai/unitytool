#!/bin/bash

SERVICE_NAME=$1
if [ -z "$1" ]; then
	SERVICE_NAME="unitool_serve_linux"
fi

ps -ef | grep $SERVICE_NAME | awk '{print $2}' | xargs kill -15

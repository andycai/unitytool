#!/bin/bash

SERVICE_NAME=$1
if [ -z "$1" ]; then
	SERVICE_NAME="stats_serve_linux"
fi

ps -ef | grep $1 | awk '{print $2}' | xargs kill -15

#!/bin/bash

nohup ./stats_serve_linux -port 8080 > output.log 2>&1 &

#!/bin/bash

nohup ./myhealthserver \
        -a 192.168.100.100:8080 \
        -u USER_ID \
        -c ./server.crt \
        -k ./server.key \
        -s ./filestorage \
        -d ./myhealth.db &

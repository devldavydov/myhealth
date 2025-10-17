#!/bin/bash

nohup ./myhealthserver \
        -a 192.168.100.100:8080 \
        -u USER_ID \
        -d ./myhealth.db &

#!/bin/bash

nohup ./myhealthbot \
        -t BOT_TOKEN \
        -u USER_ID \
        -d ./myhealth.db &

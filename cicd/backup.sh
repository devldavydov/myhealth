#!/bin/bash

# Add to cron
# 0 2 * * * /root/MyProjects/myhealth/backup.sh

cd /root/MyProjects/myhealth

sqlite3 myhealth.db ".backup 'myhealth_$(date +%Y%m%d).db'"
find . -type f -mtime +5 -name 'myhealth_*.db' -execdir rm -- '{}' \;

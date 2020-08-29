#!/bin/sh

TraceID=$1
echo "Start Logging ($TraceID)"

vmstat 1 90 > /tmp/vmstat.log 2>/dev/null  &

sleep 90

sudo pkill -HUP isucon9final

ACCESSLOG=/var/log/nginx/access.log.`date +%Y%m%d%H%M%S`
sudo mv /var/log/nginx/access.log $ACCESSLOG
sudo systemctl restart nginx
cp $ACCESSLOG /tmp/access.log
chmod 644 /tmp/access.log

sleep 1

mkdir -p /tmp/isucon/logs/$TraceID
cp /tmp/access.log /tmp/isucon/logs/$TraceID/access.log
cp /tmp/cpu.pprof /tmp/isucon/logs/$TraceID/cpu.pprof
cp /tmp/perf.log /tmp/isucon/logs/$TraceID/perf.log
cp /tmp/sql.log /tmp/isucon/logs/$TraceID/sql.log
cp /tmp/vmstat.log /tmp/isucon/logs/$TraceID/vmstat.log
rm /tmp/isucon/logs/latest
ln -sf $TraceID /tmp/isucon/logs/latest

echo "End Logging  ($TraceID)"

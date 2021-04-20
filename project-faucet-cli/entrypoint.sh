#!/bin/sh

if pidof $CLI_NAME; then
    export PATH=/proc/$(pidof $CLI_NAME)/root/usr/local/sbin:/proc/$(pidof $CLI_NAME)/root/usr/local/bin:/proc/$(pidof $CLI_NAME)/root/usr/sbin:/proc/$(pidof $CLI_NAME)/root/usr/bin:/proc/$(pidof $CLI_NAME)/root/sbin:/proc/$(pidof $CLI_NAME)/root/bin:$PATH
fi

exec "$@"

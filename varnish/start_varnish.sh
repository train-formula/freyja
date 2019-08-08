#!/bin/bash

VARNISH_VCL_TEMPLATE=`realpath /varnish_template.vcl`
VARNISH_WORKING=`realpath /varnish`

cp $VARNISH_VCL_TEMPLATE /varnish.vcl

envs=`printenv`

for env in $envs
do
    IFS== read name value <<< "$env"

    sed -i "s|\${${name}}|${value}|g" /varnish.vcl
done

VARNISH_VCL=`realpath /varnish.vcl`

varnishd -a 0.0.0.0:$VARNISH_PORT -s malloc,$VARNISH_MALLOC_SIZE, -n $VARNISH_WORKING -f $VARNISH_VCL
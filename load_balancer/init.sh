#! /user/bin/env bash
# Configures the nginx server upon startup.

# Assign configuration values to enviroment values.
export HOST_IP=$1
export RANK_PORT=$2
export REVERSE_INDEX_PORT=$3

# Produce nginx configuration from template
/usr/bin/envsubst '$HOST_IP $RANK_PORT $REVERSE_INDEX_PORT' < nginx.conf.template > /etc/nginx/nginx.conf

# Show generated config for the sake of debugging.
cat /etc/nginx/nginx.conf

# Start nginx.
nginx -g 'daemon off;' 2> start_error.log
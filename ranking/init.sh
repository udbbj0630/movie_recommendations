#! /user/bin/env bash
# Configures the nginx server upon startup.
export HOST_IP=$1
# Assign configuration values to enviroment values.
mkdir /data/
cp movies_metadata.csv /data/movies_metadata.csv
echo 'copy ok $?'

# Start nginx.
go run serve_ranking.go -tfx_url=http://$HOST_IP:8006
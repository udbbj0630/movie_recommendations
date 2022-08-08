#!/bin/bash
while getopts "b" arg
do
  case $arg in
    b)
      docker build -f Dockerfile -t movie_recommendation/ranking:0.1 .
      ;;
    ?)
      echo "Unknown args $arg... exit..."
      exit 1
      ;;
  esac
done
docker stop ranking_dev
docker run -it -d --rm --name ranking_dev -v /data:/data /Users/shuangyueli/movie_recommendations/movie_recommendations/ranking:/ranking movie_recommendation/ranking:0.1 /bin/bash

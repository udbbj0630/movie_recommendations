#!/bin/bash
while getopts "b" arg
do
  case $arg in
    b)
      docker rmi movie_recommendation/ranking:0.0
      docker build -f Dockerfile -t movie_recommendation/ranking:0.0 .
      ;;
    ?)
      echo "Unknown args $arg... exit..."
      exit 1
      ;;
  esac
done
docker stop ranking_dev
docker run -it -d --rm --name ranking_dev -v /Users/shuangyueli/movie_recommendations/movie_recommendations/ranking:/ranking movie_recommendation/ranking:0.0 /bin/bash

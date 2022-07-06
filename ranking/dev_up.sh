#!/bin/bash
docker build -f Dockerfile -t movie_recommendation/ranking:0.0 .
docker run -it --rm -v /Users/shuangyueli/movie_recommendations/movie_recommendations/ranking:/ranking movie_recommendation/ranking:0.0 /bin/bash

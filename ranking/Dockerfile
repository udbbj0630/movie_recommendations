FROM  golang:1.19-rc

RUN ["apt-get", "update"]
RUN ["apt-get", "install", "-y", "vim"]

EXPOSE 80

VOLUME ["/ranking", "/data"]

WORKDIR /ranking

CMD go run serve_ranking.go

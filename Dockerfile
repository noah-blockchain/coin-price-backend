# docker build --no-cache -t coin-price-backend:latest -f Dockerfile .
# docker build -t coin-price-backend:latest -f Dockerfile .
# docker run -d -p 127.0.0.1:10500:10500 --restart=always coin-price-backend:latest

FROM golang:1.12-buster as builder

ENV APP_PATH /home/coin-price-backend

COPY . ${APP_PATH}

WORKDIR ${APP_PATH}

RUN make create_vendor && make build

FROM debian:buster-slim as executor

COPY --from=builder /home/coin-price-backend/build/coin-history /usr/local/bin/coin-history

COPY --from=builder /home/coin-price-backend/migrations /migrations

EXPOSE 10500
CMD ["coin-history"]
STOPSIGNAL SIGTERM

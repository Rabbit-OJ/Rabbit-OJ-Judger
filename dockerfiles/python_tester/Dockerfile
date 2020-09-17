FROM python:3

ENV Role="Tester"

WORKDIR /app
COPY ./tester /app/server
RUN chmod +x /app/server

RUN mkdir -p /result && mkdir -p /case

ENTRYPOINT ["/app/server"]
FROM debian:buster-slim

RUN apt-get update && apt-get install ca-certificates -y

WORKDIR /app

COPY api ./

EXPOSE 7344

CMD ["./api"]

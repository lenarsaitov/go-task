FROM ubuntu:21.10 as server

RUN apt-get -y update
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get install -y postgresql

USER root

RUN apt-get install -y golang git

COPY ./ ./

RUN go mod download
RUN go build -o go_server ./cmd/main.go

EXPOSE 10000
EXPOSE 5432

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql -d docker -a -f ./sql/link.sql &&\
    /etc/init.d/postgresql stop

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

CMD service postgresql start && ./go_server -db postgres
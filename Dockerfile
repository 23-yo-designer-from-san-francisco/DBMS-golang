FROM ubuntu:latest

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt -y update
RUN apt install -y wget gnupg lsb-release ca-certificates --no-install-recommends
RUN sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
RUN wget -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -

RUN apt -y update

RUN apt install -y postgresql-14 --no-install-recommends

RUN cd /tmp && wget https://golang.org/dl/go1.17.2.linux-arm64.tar.gz && tar -C /usr/local -xzf go1.17.2.linux-arm64.tar.gz
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

RUN echo "local   all             postgres                                peer\nlocal   all             all                                     md5\nhost    all             all             127.0.0.1/32            scram-sha-256\nhost    all             all             0.0.0.0/0               md5" > /etc/postgresql/14/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/14/main/postgresql.conf
RUN echo "synchronous_commit = off\nfsync = off\nshared_buffers = 1GB\nunix_socket_directories = '/var/run/postgresql'\nunix_socket_permissions = 0777" >> /etc/postgresql/14/main/postgresql.conf

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .

EXPOSE 5000

ENV DBPORT 5432
ENV DBNAME docker
ENV DBUSER docker
ENV DBPASS docker
ENV PGPASSWORD docker

CMD service postgresql start && psql -h localhost -d docker -U docker -p 5432 -a -q -f ./db/init.sql && /usr/local/go/bin/go run main.go

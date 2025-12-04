# A Postgres image pre-loaded with Chinook.

FROM postgres:latest

RUN apt-get update
RUN apt-get install -y wget

ENV CHINOOK=https://raw.githubusercontent.com/lerocha/chinook-database/9ad881311025c290a911a0461dd9919317b09021/ChinookDatabase/DataSources/Chinook_PostgreSql.sql
RUN wget -O /docker-entrypoint-initdb.d/init.sql $CHINOOK

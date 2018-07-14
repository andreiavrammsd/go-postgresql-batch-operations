FROM postgres:9.6-alpine
COPY schema.sql /docker-entrypoint-initdb.d/

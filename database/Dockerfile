
FROM postgres:latest
ENV POSTGRES_USER docker
ENV POSTGRES_PASSWORD docker
ENV POSTGRES_DB filehashes
COPY init.sql /docker-entrypoint-initdb.d/
EXPOSE 5432

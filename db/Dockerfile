FROM library/postgres
ENV POSTGRES_USER postgres-dev
ENV POSTGRES_PASSWORD s3cr3tp4ssw0rd
ENV POSTGRES_DB dev
COPY init.sql /docker-entrypoint-initdb.d/
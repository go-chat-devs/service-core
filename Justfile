create-val-postgres:
    docker volume create postgres-data

postgres-run:
    docker run --name postgres \
    -e POSTGRES_PASSWORD=admin \
    -p 7765:5432  \
    -v postgres-data:/var/lib/postgresql/data \
    -d postgres:18-bookworm
# Highload Architect - Social Network

## How-to Run
### Build
```bash
make build-docker
```
### Prepare environment
You need to copy `.env.docker.example` and modify(or keep default) envrionment variables.
```bash
cp .env.docker.example .env.docker
vim .env.docker
```
### Run stack
```bash
cd deploy
docker-compose up -d
```

## API
Postman collection is located in `api/` folder.

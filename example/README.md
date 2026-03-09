# Example Deploy (Docker Hub Images)

## 1) Prepare env

```bash
cp .env.sample .env
# edit .env with your Docker Hub username and tag
```

## 2) Run

```bash
docker compose up -d
```

## 3) Check logs

```bash
docker compose logs -f
```

## 4) Stop

```bash
docker compose down
```

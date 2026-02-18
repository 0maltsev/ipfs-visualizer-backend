# IPFS Cluster Topology Orchestrator

Оркестратор IPFS-кластера с визуальным редактором топологии на canvas.

## Возможности

- **Canvas-редактор**: перетаскивание узлов, соединение их рёбрами
- **Топология**: ребро A → B означает, что узел A bootstraps к узлу B (B — bootstrap)
- **Деплой в Kubernetes**: топология разворачивается как StatefulSet IPFS-кластер (по образцу `ipfs-cluster-deployment/`)

## Требования

- Go 1.25+
- PostgreSQL
- Доступ к Kubernetes (kubeconfig)
- Node.js 18+ (для frontend)

## Запуск

### Backend

```bash
# Создать .env (см. build/.env.example)
cp build/.env.example .env

# Запустить
go run ./cmd/app
```

Сервер: `http://localhost:3001`

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Фронтенд: `http://localhost:5173` (проксирует API на :3001)

### Переменные окружения

| Переменная | Описание |
|------------|----------|
| SERVER_ADDRESS_PORT | Порт сервера (default: 3001) |
| POSTGRE_SQL_* | Подключение к PostgreSQL |
| KUBE_CONFIG_PATH | Путь к kubeconfig |
| MANUAL_KUBE_CONFIG_FLAG | true — использовать файл kubeconfig |

## API

См. `docs/openapi.yml`

Основные эндпоинты:
- `GET /v1/topologies` — список топологий
- `POST /v1/topologies` — создать
- `GET /v1/topologies/{id}` — получить
- `PUT /v1/topologies/{id}` — обновить (узлы и рёбра)
- `DELETE /v1/topologies/{id}` — удалить
- `POST /v1/topologies/{id}/deploy` — задеплоить в K8s
- `POST /v1/topologies/{id}/undeploy` — удалить из K8s
- `GET /v1/topologies/{id}/status` — статус деплоя

## Папка ipfs-cluster-deployment

Содержит эталонные манифесты Kubernetes для IPFS-кластера:
- `ipfs-cluster.yaml` — StatefulSet
- `ipfs-service.yaml` — Service
- `bootstrap-script-cm.yaml` — ConfigMap со скриптами
- `secret.yaml` — ConfigMap и Secret

Оркестратор генерирует аналогичные ресурсы по топологии, созданной пользователем.

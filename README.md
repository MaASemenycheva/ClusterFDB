

# FoundationDb (role model + docker-контейнеры + Prometheus-метрики) + fdb-кластер


* `FDB_API_VERSION`
* `FDB_CLUSTER_FILE`
* `FDB_CREATE_CLUSTER_FILE`
* `FDB_EXPORT_WORKLOAD`
* `FDB_METRICS_LISTEN`
* `FDB_METRICS_EVERY`


```bash
cd deployment/docker-compose
docker-compose up --build

# Metrics will be available at
curl localhost:8081/metrics | findstr fdb | findstr -v "#"
```



![Alt text](./img/photo_2022-11-21_04-39-36.jpg?raw=true "Замеры времени")


1. Процессы с мертками скорости: запись по 10000 строк каждого типа, чтениие;
2. Ролевая модель - storage, log, только запись, только чтение;
3. Провести проверки для всех ролей;
4. Продемонстрировать результаты, описанине кейсов и код.

<h2>Замеры времени по каждой операции (в нано и миллисекундах):</h2>

![Alt text](./img/photo_2022-10-23_22-58-48.jpg?raw=true "Замеры времени по каждой операции (в нано и миллисекундах)")

<h2>Поднятие кластера FoundationDb:</h2>

![Alt text](./img/photo_2022-10-23_22-58-48.jpg?raw=true "Поднятие кластера FoundationDb")

![Alt text](./img/photo_2022-10-23_23-01-53.jpg?raw=true "Поднятие кластера FoundationDb")
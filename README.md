1. build base image

```
make base
```

2. start 

```
docker-compose up
```

3. exec in producer container

```
curl localhost:8080/ -v
```

watch docker-compose logs, counter from consumer should be logged

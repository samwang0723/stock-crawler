# stock-crawler
Taiwan stock crawling service

## Setup Redis

### Start Docker container
```

$ docker-compose -p redis -f build/docker/redis/docker-compose.yml up
```

#### Customize Sentinel
Remember to change the image location inside redis/docker-compose.yml
```
$ cd build/docker/redis
$ docker build -t {Repo}/redis-sentinel -f Dockerfile .
```

## Setup Kafka

### Start Kafka container
```
$ docker-compose -p kafka -f build/docker/kafka/docker-compose.yml up
```

### Create topics
```
$ docker exec -it kafka-1 bash

./bin/kafka-topics.sh --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 --create --topic stakeconcentration-v1 --replication-factor 2 --partitions 3
./bin/kafka-topics.sh --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 --create --topic dailycloses-v1 --replication-factor 2 --partitions 3
./bin/kafka-topics.sh --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 --create --topic stocks-v1 --replication-factor 2 --partitions 3
./bin/kafka-topics.sh --bootstrap-server kafka-1:9092,kafka-2:9092,kafka-3:9092 --create --topic threeprimary-v1 --replication-factor 2 --partitions 3
```

## Start Application

### Build image for Mac M1
```
$ make docker-m1
```

### Start stock-crawler container
```
$ docker-compose -p stock-crawler -f build/docker/app/docker-compose.yml up
```

### Environment configuration

Please configure `.env` under project root folder
1. Concentration data crawling need to use proxy to prevent rate limiting from source website, recommend to use https://app.webscrapingapi.com,
can set `WEB_SCRAPING={API_KEY}`, or https://proxycrawl.com, set `PROXY_CRAWL={API_KEY}`
2. If want to test the functionality, use `INCLUDE_WEEKEND=true`, otherwise please set to `false`, 
enabling the testing mode allows you to not skip weekend dates

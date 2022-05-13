FROM quay.io/strimzi/kafka:latest-kafka-3.1.0-amd64
COPY ./server.properties /opt/kafka/config/kraft/server.properties
COPY ./create-topic.sh /opt/kafka/bin/create-topic.sh
COPY ./start.sh /opt/kafka/bin/start.sh
ENTRYPOINT ["bin/start.sh"]

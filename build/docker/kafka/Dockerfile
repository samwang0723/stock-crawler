FROM quay.io/strimzi/kafka:latest-kafka-3.1.0-amd64
COPY ./server.properties /opt/kafka/config/kraft/server.properties
COPY ./create-topic.sh /opt/kafka/bin/create-topic.sh
USER root
RUN ["chmod", "+x", "/opt/kafka/bin/create-topic.sh"]
USER kafka
CMD bash -c  "bin/create-topic.sh & bin/kafka-storage.sh format -t $(bin/kafka-storage.sh random-uuid) -c /opt/kafka/config/kraft/server.properties && bin/kafka-server-start.sh /opt/kafka/config/kraft/server.properties"

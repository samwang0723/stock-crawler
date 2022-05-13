#!/bin/bash
bin/kafka-storage.sh format -t $(bin/kafka-storage.sh random-uuid) -c /opt/kafka/config/kraft/server.properties && bin/kafka-server-start.sh /opt/kafka/config/kraft/server.properties &
bin/create-topic.sh
# Wait for any process to exit
wait -n
        
# Exit with status of process that exited first
exit $?

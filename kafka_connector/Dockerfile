FROM confluentinc/cp-kafka-connect:7.6.0

USER root

RUN yum install -y gettext && yum clean all

COPY connectors /connectors
COPY init.sh /init.sh
COPY plugins /usr/share/confluent-hub-components

RUN chmod +x /init.sh

CMD ["/bin/bash", "-c", "/etc/confluent/docker/run & sleep 20 && /init.sh && wait"]
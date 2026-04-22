package com.example.transactional_orchestrator.config;

import com.example.transactional_orchestrator.dto.TransactionMessage;
import com.example.transactional_orchestrator.dto.TransactionalFailedMessage;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.ConcurrentKafkaListenerContainerFactory;
import org.springframework.kafka.core.ConsumerFactory;
import org.springframework.kafka.core.DefaultKafkaConsumerFactory;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.listener.DefaultErrorHandler;
import org.springframework.kafka.support.serializer.JsonDeserializer;
import org.springframework.util.backoff.FixedBackOff;

import java.util.Map;

@Configuration
public class KafkaConfig {

    @Bean
    public ConsumerFactory<String, TransactionMessage> transactionConsumerFactory(KafkaProperty kafkaProperty) {
        return  getDefaultKafkaConsumerFactory(kafkaProperty, TransactionMessage.class);
    }

    @Bean
    public ConsumerFactory<String, TransactionalFailedMessage> failedConsumerFactory(KafkaProperty kafkaProperty) {
        return getDefaultKafkaConsumerFactory(kafkaProperty, TransactionalFailedMessage.class);
    }

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, TransactionMessage>
    transactionKafkaListenerContainerFactory(ConsumerFactory<String, TransactionMessage> consumerFactory, DefaultErrorHandler errorHandler) {
        ConcurrentKafkaListenerContainerFactory<String, TransactionMessage> factory =
                new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(consumerFactory);
        factory.setCommonErrorHandler(errorHandler);
        return factory;
    }

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, TransactionalFailedMessage>
    failedKafkaListenerContainerFactory(ConsumerFactory<String, TransactionalFailedMessage> consumerFactory, DefaultErrorHandler errorHandler) {
        ConcurrentKafkaListenerContainerFactory<String, TransactionalFailedMessage> factory =
                new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(consumerFactory);
        factory.setCommonErrorHandler(errorHandler);
        return factory;
    }

    @Bean
    public DefaultErrorHandler errorHandler() {
        return new DefaultErrorHandler(
                new FixedBackOff(200L, 5)
        );
    }

    private static <MT> DefaultKafkaConsumerFactory<String, MT> getDefaultKafkaConsumerFactory(KafkaProperty kafkaProperty, Class<MT> messageClass) {
        Map<String, Object> props = kafkaProperty.baseProps();
        return new DefaultKafkaConsumerFactory<>(
                props,
                new StringDeserializer(),
                new JsonDeserializer<>(messageClass, false)
        );
    }
}
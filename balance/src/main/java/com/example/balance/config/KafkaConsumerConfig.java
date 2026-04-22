package com.example.balance.config;

import com.example.balance.dto.TransactionalFailedMessage;
import com.example.balance.dto.TransactionalMessage;
import com.example.balance.mapper.MessageMapper;
import com.example.balance.service.OutboxService;
import lombok.RequiredArgsConstructor;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.TopicPartition;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.ConcurrentKafkaListenerContainerFactory;
import org.springframework.kafka.core.ConsumerFactory;
import org.springframework.kafka.core.DefaultKafkaConsumerFactory;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.listener.DefaultErrorHandler;
import org.springframework.kafka.listener.DeadLetterPublishingRecoverer;
import org.springframework.kafka.support.serializer.JsonDeserializer;
import org.springframework.util.backoff.FixedBackOff;

import java.util.HashMap;
import java.util.Map;

@Configuration
@RequiredArgsConstructor
public class KafkaConsumerConfig {


    @Bean
    public ConsumerFactory<String, TransactionalMessage> transactionConsumerFactory(KafkaProperty kafkaProperty) {
        return  getDefaultKafkaConsumerFactory(kafkaProperty, TransactionalMessage.class);
    }

    @Bean
    public ConsumerFactory<String, TransactionalFailedMessage> failedConsumerFactory(KafkaProperty kafkaProperty) {
        return getDefaultKafkaConsumerFactory(kafkaProperty, TransactionalFailedMessage.class);
    }

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, TransactionalMessage>
    transactionKafkaListenerContainerFactory(ConsumerFactory<String, TransactionalMessage> consumerFactory, DefaultErrorHandler errorHandler) {
        ConcurrentKafkaListenerContainerFactory<String, TransactionalMessage> factory =
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

    private static <MT> DefaultKafkaConsumerFactory<String, MT> getDefaultKafkaConsumerFactory(KafkaProperty kafkaProperty, Class<MT> messageClass) {
        Map<String, Object> props = kafkaProperty.baseProps();
        return new DefaultKafkaConsumerFactory<>(
                props,
                new StringDeserializer(),
                new JsonDeserializer<>(messageClass, false)
        );
    }

    @Bean
    public DefaultErrorHandler errorHandler(MessageMapper mapper,
                                            OutboxService outboxService,
                                            KafkaTemplate<Object, Object> kafkaTemplate) {
        FixedBackOff backOff = new FixedBackOff(500L, 5);
        return new DefaultErrorHandler((record, ex) -> {

            TransactionalFailedMessage failedMessage = mapper.toFailed((TransactionalMessage) record.value());
            failedMessage.setMessageError(ex.getCause().getMessage());
            outboxService.createTransactionalFailedMessage(failedMessage);
        }, backOff);
    }
}
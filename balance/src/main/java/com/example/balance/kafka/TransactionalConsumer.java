package com.example.balance.kafka;

import com.example.balance.dto.TransactionalFailedMessage;
import com.example.balance.dto.TransactionalMessage;
import com.example.balance.exception.OperationNotSupportedException;
import com.example.balance.service.TransactionService;
import lombok.RequiredArgsConstructor;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class TransactionalConsumer {

    private final TransactionService service;

    @KafkaListener(
            topics = "transaction",
            containerFactory = "transactionKafkaListenerContainerFactory"
    )
    public void listen(TransactionalMessage message) {
            switch (message.getEventType()){
                case CREATED -> service.handleCreated(message);
                case PROCESSED -> service.handleProcessed(message);
                case CHANGED, RESERVED -> {}
                case null -> throw new OperationNotSupportedException();
            }
    }

    @KafkaListener(
            topics = "failed",
            containerFactory = "failedKafkaListenerContainerFactory"
    )
    public void listenError(TransactionalFailedMessage message) {
        service.handleFailed(message);
    }
}
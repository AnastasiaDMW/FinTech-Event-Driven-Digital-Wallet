package com.example.transactional_orchestrator.kafka;

import com.example.transactional_orchestrator.dto.TransactionMessage;
import com.example.transactional_orchestrator.dto.TransactionalFailedMessage;
import com.example.transactional_orchestrator.service.TransactionService;
import lombok.RequiredArgsConstructor;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

import static com.example.transactional_orchestrator.dto.EventType.CREATED;

@Service
@RequiredArgsConstructor
public class TransactionalConsumer {

    private final TransactionService service;

    @KafkaListener(
            topics = "transaction",
            containerFactory = "transactionKafkaListenerContainerFactory"
    )
    public void listen(TransactionMessage message) {
        if (CREATED.equals(message.getEventType())) return;
        service.changeStatus(message);
    }

    @KafkaListener(
            topics = "failed",
            containerFactory = "failedKafkaListenerContainerFactory"
    )
    public void listenError(TransactionalFailedMessage message) {
        service.setFailedStatus(message);
    }
}
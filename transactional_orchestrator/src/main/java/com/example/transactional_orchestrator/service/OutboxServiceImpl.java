package com.example.transactional_orchestrator.service;

import com.example.transactional_orchestrator.dto.NotificationMessage;
import com.example.transactional_orchestrator.dto.TransactionMessage;
import com.example.transactional_orchestrator.model.Category;
import com.example.transactional_orchestrator.model.Outbox;
import com.example.transactional_orchestrator.repository.OutboxRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import static com.example.transactional_orchestrator.dto.EventType.CREATED;
import static com.example.transactional_orchestrator.model.Category.NOTIFICATION;
import static com.example.transactional_orchestrator.model.Category.TRANSACTION;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class OutboxServiceImpl implements OutboxService {

    private final OutboxRepository repo;
    private final ObjectMapper mapper;

    @Override
    public void createNotificationMessage(NotificationMessage message) {
        create(NOTIFICATION, message.getId().toString(), message);
    }

    @Override
    public void sendTransactionCreatedMessage(TransactionMessage message) {
        message.setEventType(CREATED);
        create(TRANSACTION, message.getId().toString(), message);
    }

    private void create(Category category, String eventKey, Object payload) {
        var outbox = Outbox.builder()
                .category(category)
                .eventKey(eventKey)
                .payload(mapper.valueToTree(payload))
                .build();
        repo.save(outbox);
    }
}

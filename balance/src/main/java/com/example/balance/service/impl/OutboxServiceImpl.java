package com.example.balance.service.impl;

import com.example.balance.dto.TransactionalFailedMessage;
import com.example.balance.dto.TransactionalMessage;
import com.example.balance.model.Category;
import com.example.balance.model.Outbox;
import com.example.balance.repository.OutboxRepository;
import com.example.balance.service.OutboxService;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import static com.example.balance.dto.EventType.CHANGED;
import static com.example.balance.dto.EventType.RESERVED;
import static com.example.balance.model.Category.FAILED;
import static com.example.balance.model.Category.TRANSACTION;

@Service
@RequiredArgsConstructor
@Transactional
public class OutboxServiceImpl implements OutboxService {

    private final OutboxRepository repo;
    private final ObjectMapper mapper;

    @Override
    public void createTransactionalReservedMessage(TransactionalMessage message) {
        message.setEventType(RESERVED);
        create(TRANSACTION, message.getId().toString(), message);
    }

    @Override
    public void createTransactionalChangedMessage(TransactionalMessage message) {
        message.setEventType(CHANGED);
        create(TRANSACTION, message.getId().toString(), message);
    }

    @Override
    public void createTransactionalFailedMessage(TransactionalFailedMessage message) {
        create(FAILED, message.getId().toString(), message);
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

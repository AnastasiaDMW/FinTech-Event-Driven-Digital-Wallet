package com.example.balance.service.impl;

import com.example.balance.dto.TransactionalChangedMessage;
import com.example.balance.dto.TransactionalReservedMessage;
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
import static com.example.balance.model.Category.TRANSACTIONAL;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class OutboxServiceImpl implements OutboxService {

    private final OutboxRepository repo;
    private final ObjectMapper mapper;

    @Override
    public void createTransactionalReservedMessage(TransactionalReservedMessage message) {
        message.setEventType(RESERVED);
        create(TRANSACTIONAL, message.getId().toString(), message);
    }

    @Override
    public void createTransactionalChangedMessage(TransactionalChangedMessage message) {
        message.setEventType(CHANGED);
        create(TRANSACTIONAL, message.getId().toString(), message);
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

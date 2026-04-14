package com.example.transactional_orchestrator.service;

import com.example.transactional_orchestrator.dto.NotificationMessage;
import com.example.transactional_orchestrator.dto.TransactionCreatedMessage;

public interface OutboxService {
    void createNotificationMessage(NotificationMessage message);

    void createTransactionalCreatedMessage(TransactionCreatedMessage message);
}

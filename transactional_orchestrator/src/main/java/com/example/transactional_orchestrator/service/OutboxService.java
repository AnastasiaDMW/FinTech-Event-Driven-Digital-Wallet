package com.example.transactional_orchestrator.service;

import com.example.transactional_orchestrator.dto.NotificationMessage;
import com.example.transactional_orchestrator.dto.TransactionMessage;

public interface OutboxService {
    void createNotificationMessage(NotificationMessage message);

    void sendTransactionCreatedMessage(TransactionMessage message);
}

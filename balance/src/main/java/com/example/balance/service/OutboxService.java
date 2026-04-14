package com.example.balance.service;

import com.example.balance.dto.TransactionalChangedMessage;
import com.example.balance.dto.TransactionalReservedMessage;
import com.example.transactional_orchestrator.dto.NotificationMessage;
import com.example.transactional_orchestrator.dto.TransactionCreatedMessage;

public interface OutboxService {

    void createTransactionalReservedMessage(TransactionalReservedMessage message);

    void createTransactionalChangedMessage(TransactionalChangedMessage message);
}

package com.example.balance.service;

import com.example.balance.dto.TransactionalFailedMessage;
import com.example.balance.dto.TransactionalMessage;
import org.springframework.transaction.annotation.Transactional;

public interface TransactionService {
    @Transactional
    void handleCreated(TransactionalMessage message);

    @Transactional
    void handleProcessed(TransactionalMessage message);

    @Transactional
    void handleFailed(TransactionalFailedMessage message);
}

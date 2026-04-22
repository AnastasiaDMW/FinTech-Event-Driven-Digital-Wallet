package com.example.balance.service;

import com.example.balance.dto.TransactionalFailedMessage;
import com.example.balance.dto.TransactionalMessage;

public interface OutboxService {

    void createTransactionalReservedMessage(TransactionalMessage message);

    void createTransactionalChangedMessage(TransactionalMessage message);

    void createTransactionalFailedMessage(TransactionalFailedMessage message);
}

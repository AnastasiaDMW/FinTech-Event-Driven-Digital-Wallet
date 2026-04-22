package com.example.transactional_orchestrator.service;

import com.example.transactional_orchestrator.dto.TransactionMessage;
import com.example.transactional_orchestrator.dto.TransactionRequest;
import com.example.transactional_orchestrator.dto.TransactionResponse;
import com.example.transactional_orchestrator.dto.TransactionalFailedMessage;

public interface TransactionService {
    TransactionResponse createTransaction(TransactionRequest request);

    TransactionResponse getTransactionInfo(Long transactionalId);

    void changeStatus(TransactionMessage message);

    void setFailedStatus(TransactionalFailedMessage message);
}

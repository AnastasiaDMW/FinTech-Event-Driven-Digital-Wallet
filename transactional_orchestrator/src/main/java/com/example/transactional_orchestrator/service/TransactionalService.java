package com.example.transactional_orchestrator.service;

import com.example.transactional_orchestrator.dto.TransactionRequest;
import com.example.transactional_orchestrator.dto.TransactionResponse;

public interface TransactionalService {
    TransactionResponse createTransaction(TransactionRequest request);

    TransactionResponse getTransactionInfo(Long transactionalId);
}

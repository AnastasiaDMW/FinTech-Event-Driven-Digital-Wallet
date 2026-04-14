package com.example.transactional_orchestrator.controller;

import com.example.transactional_orchestrator.api.TransactionApi;
import com.example.transactional_orchestrator.dto.TransactionRequest;
import com.example.transactional_orchestrator.dto.TransactionResponse;
import com.example.transactional_orchestrator.service.TransactionalService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import static com.example.transactional_orchestrator.api.TransactionApi.API_V1_TRANSACTION;

@RequestMapping(API_V1_TRANSACTION)
@RestController
@RequiredArgsConstructor
public class TransactionController implements TransactionApi {

    private final TransactionalService service;

    @Override
    public TransactionResponse createTransaction(TransactionRequest request) {
        return service.createTransaction(request);
    }

    @Override
    public TransactionResponse getTransactionInfo(Long transactionalId) {
        return service.getTransactionInfo(transactionalId);
    }
}

package com.example.transactional_orchestrator.api;

import com.example.transactional_orchestrator.dto.TransactionRequest;
import com.example.transactional_orchestrator.dto.TransactionResponse;
import jakarta.validation.Valid;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

public interface TransactionApi {

    String API_V1 = "/api/v1";
    String API_V1_TRANSACTION = API_V1 + "/transaction";
    String VARIABLE_ID = "/{id}";

    @PostMapping
    TransactionResponse createTransaction(@RequestBody @Valid TransactionRequest request);

    @GetMapping(VARIABLE_ID)
    TransactionResponse getTransactionInfo(@PathVariable(value = "id") Long transactionalId);
}

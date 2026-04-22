package com.example.transactional_orchestrator.dto;

import com.example.transactional_orchestrator.model.Status;
import com.example.transactional_orchestrator.model.TransactionType;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.UUID;

@Getter
@Setter
@NoArgsConstructor
public class TransactionResponse {

    private Long id;
    private UUID accountFrom;
    private UUID accountTo;
    private BigInteger amount;
    private UUID idempotent;
    private TransactionType type;
    private Status status;
}

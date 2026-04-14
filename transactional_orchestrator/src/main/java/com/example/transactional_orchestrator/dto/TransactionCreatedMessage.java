package com.example.transactional_orchestrator.dto;

import com.example.transactional_orchestrator.model.TransactionType;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.math.BigDecimal;
import java.util.UUID;

@Getter
@Setter
@NoArgsConstructor
public class TransactionCreatedMessage {

    private Long id;
    private UUID fromAccountId;
    private UUID toAccountId;
    private BigDecimal amount;
    private UUID idempotent;
    private TransactionType type;
    private EventType eventType;
}

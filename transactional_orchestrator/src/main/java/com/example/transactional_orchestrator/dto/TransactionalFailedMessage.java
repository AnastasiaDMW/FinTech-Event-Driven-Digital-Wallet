package com.example.transactional_orchestrator.dto;

import com.example.transactional_orchestrator.model.TransactionType;
import lombok.*;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.UUID;

@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class TransactionalFailedMessage {
    private Long id;
    private UUID accountFrom;
    private UUID accountTo;
    private BigInteger amount;
    private UUID idempotent;
    private TransactionType type;
    private EventType eventType;
    private String messageError;
}

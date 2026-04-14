package com.example.balance.dto;

import lombok.Getter;
import lombok.Setter;

import java.math.BigDecimal;
import java.util.UUID;

import static com.example.balance.dto.EventType.PROCESSED;

@Getter
@Setter
public class TransactionalProcessedMessage {

    private Long id;
    private UUID accountFrom;
    private UUID accountTo;
    private BigDecimal amount;
    private UUID idempotent;
    private TransactionType type;
    private EventType eventType = PROCESSED;
}

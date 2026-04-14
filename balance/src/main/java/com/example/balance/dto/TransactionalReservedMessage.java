package com.example.balance.dto;

import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

import java.math.BigDecimal;
import java.util.UUID;

import static com.example.balance.dto.EventType.RESERVED;

@Getter
@Setter
@Builder
public class TransactionalReservedMessage {

    private Long id;
    private UUID accountFrom;
    private UUID accountTo;
    private BigDecimal amount;
    private UUID idempotent;
    private TransactionType type;
    private EventType eventType = RESERVED;
}

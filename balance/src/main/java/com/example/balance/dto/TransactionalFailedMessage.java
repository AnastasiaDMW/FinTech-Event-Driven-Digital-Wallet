package com.example.balance.dto;

import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

import java.math.BigDecimal;
import java.util.UUID;

import static com.example.balance.dto.EventType.FAILED;

@Getter
@Setter
@Builder
public class TransactionalFailedMessage {
    private Long id;
    private UUID accountFrom;
    private UUID accountTo;
    private BigDecimal amount;
    private UUID idempotent;
    private TransactionType type;
    private EventType eventType = FAILED;
    private String messageError;
}

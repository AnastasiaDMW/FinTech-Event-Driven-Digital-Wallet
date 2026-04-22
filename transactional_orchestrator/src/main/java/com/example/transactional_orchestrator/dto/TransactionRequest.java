package com.example.transactional_orchestrator.dto;

import com.example.transactional_orchestrator.model.TransactionType;
import jakarta.validation.constraints.AssertTrue;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.UUID;

import static java.util.Objects.nonNull;

@Getter
@Setter
@NoArgsConstructor
public class TransactionRequest {

    @NotNull
    private UUID idempotent;
    private UUID accountFrom;
    private UUID accountTo;
    @NotNull
    @Positive
    private BigInteger amount;
    @NotNull
    private TransactionType type;

    @AssertTrue(message = "Не корректный ввод счетов")
    public boolean isValidAccount() {
        return switch (type) {
            case DEPOSIT -> nonNull(accountTo);
            case TRANSFER -> nonNull(accountTo) && nonNull(accountFrom);
            case WITHDRAW -> nonNull(accountFrom);
            case null, default -> false;
        };
    }
}

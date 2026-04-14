package com.example.transactional_orchestrator.dto;

import com.example.transactional_orchestrator.model.TransactionType;
import jakarta.validation.constraints.AssertTrue;
import jakarta.validation.constraints.NotNull;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.math.BigDecimal;
import java.util.UUID;

import static java.util.Objects.nonNull;

@Getter
@Setter
@NoArgsConstructor
public class TransactionRequest {

    @NotNull
    private UUID idempotent;
    private UUID fromAccountId;
    private UUID toAccountId;
    @NotNull
    private BigDecimal amount;
    @NotNull
    private TransactionType type;

    @AssertTrue(message = "Не корректный ввод счетов")
    public boolean isValidAccount() {
        return switch (type){
            case DEPOSIT -> nonNull(toAccountId);
            case TRANSFER -> nonNull(toAccountId) && nonNull(fromAccountId);
            case WITHDRAW -> nonNull(fromAccountId);
            case null, default -> false;
        };
    }
}

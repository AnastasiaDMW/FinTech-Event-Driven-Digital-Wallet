package com.example.transactional_orchestrator.dto;

import jakarta.validation.constraints.NotNull;
import lombok.Getter;
import lombok.Setter;

import java.util.UUID;

@Getter
@Setter
public class AccountRequest {
    private UUID fromAccountId;
    private UUID toAccountId;

}

package com.example.transactional_orchestrator.dto;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.Setter;

import java.util.UUID;

@Getter
@Setter
@AllArgsConstructor
public class AccountRequest {
    private UUID accountFrom;
    private UUID accountTo;
}

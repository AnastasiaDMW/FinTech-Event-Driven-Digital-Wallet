package com.example.transactional_orchestrator.dto;

import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
public class AccountResponse {
    private Boolean isValidAccountFrom;
    private Boolean isValidAccountTo;
}

package com.example.transactional_orchestrator.exception;

import java.util.UUID;

public class NotFoundAccountException extends RuntimeException {
    public NotFoundAccountException(UUID account) {
        super("Счет \"%s\" не валидный".formatted(account));
    }
}

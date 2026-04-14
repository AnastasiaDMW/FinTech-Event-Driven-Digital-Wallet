package com.example.balance.exception;

import java.util.UUID;

public class InsufficientFundsException extends RuntimeException {
    public InsufficientFundsException(UUID account) {
        super("На счете \"%s\" не достаточно средств");
    }
}

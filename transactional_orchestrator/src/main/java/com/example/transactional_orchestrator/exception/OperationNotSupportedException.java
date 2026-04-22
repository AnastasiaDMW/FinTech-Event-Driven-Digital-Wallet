package com.example.transactional_orchestrator.exception;

public class OperationNotSupportedException extends RuntimeException {
    public OperationNotSupportedException() {
        super("Данная операция не поддерживается");
    }
}

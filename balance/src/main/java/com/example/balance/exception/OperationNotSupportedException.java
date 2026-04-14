package com.example.balance.exception;

public class OperationNotSupportedException extends RuntimeException {
    public OperationNotSupportedException() {
        super("Данная операция не поддерживается");
    }
}

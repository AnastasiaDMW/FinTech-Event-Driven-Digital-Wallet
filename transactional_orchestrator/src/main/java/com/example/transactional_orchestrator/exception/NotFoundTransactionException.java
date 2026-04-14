package com.example.transactional_orchestrator.exception;

public class NotFoundTransactionException extends RuntimeException {
  public NotFoundTransactionException(Long transactionId) {
    super("Транзакция №%s не найдена".formatted(transactionId));
  }
}

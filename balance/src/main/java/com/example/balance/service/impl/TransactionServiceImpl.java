package com.example.balance.service.impl;

import com.example.balance.dto.TransactionalFailedMessage;
import com.example.balance.dto.TransactionalMessage;
import com.example.balance.exception.OperationNotSupportedException;
import com.example.balance.model.Reserved;
import com.example.balance.service.BalanceService;
import com.example.balance.service.ReservedService;
import com.example.balance.service.TransactionService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class TransactionServiceImpl implements TransactionService {

    private final BalanceService balanceService;
    private final ReservedService reservedService;
    private final OutboxServiceImpl outboxService;

    @Override
    @Transactional
    public void handleCreated(TransactionalMessage message) {
        if (reservedService.isNotIdempotentRequest(message.getIdempotent())) return;
        switch (message.getType()) {
            case TRANSFER -> transfer(message.getAccountTo(), message.getAccountFrom(), message.getAmount(), message.getIdempotent());
            case DEPOSIT -> deposit(message.getAccountTo(), message.getAmount(), message.getIdempotent());
            case WITHDRAW -> withdraw(message.getAccountFrom(), message.getAmount(), message.getIdempotent());
            case null, default -> throw new OperationNotSupportedException();
        }
        outboxService.createTransactionalReservedMessage(message);
    }

    private void transfer(UUID accountTo, UUID accountFrom, BigInteger amount, UUID idempotent) {
        balanceService.checkBalance(accountFrom, amount.negate());
        withdraw(accountFrom, amount, idempotent);
        deposit(accountTo, amount, idempotent);
    }

    @Override
    @Transactional
    public void handleProcessed(TransactionalMessage message) {
        reservedService.readByIdempotent(message.getIdempotent())
                .stream()
                .filter(reserv -> reserv.getAmount().compareTo(BigInteger.ZERO) > 0)
                .forEach(item -> balanceService.deposit(item.getAccount().getAccount(), item.getAmount()));
        reservedService.confirm(message.getIdempotent());
        outboxService.createTransactionalChangedMessage(message);
    }

    @Override
    @Transactional
    public void handleFailed(TransactionalFailedMessage message) {
        var reserved = reservedService.readByIdempotent(message.getIdempotent());
        reserved.stream()
                .filter(reserv -> reserv.getAmount().compareTo(BigInteger.ZERO) < 0)
                .forEach(item -> balanceService.withdraw(item.getAccount().getAccount(), item.getAmount()));
        reservedService.removeAll(reserved);
    }

    private void withdraw(UUID accountFrom, BigInteger amount, UUID idempotent) {
        balanceService.checkBalance(accountFrom, amount.negate());
        reservedService.create(
                Reserved.builder()
                        .account(balanceService.withdraw(accountFrom, amount))
                        .amount(amount.negate())
                        .idempotent(idempotent)
                        .build()
        );
    }

    private void deposit(UUID accountTo, BigInteger amount, UUID idempotent) {
        reservedService.create(
                Reserved.builder()
                        .account(balanceService.getBalance(accountTo))
                        .amount(amount)
                        .idempotent(idempotent)
                        .build()
        );
    }
}
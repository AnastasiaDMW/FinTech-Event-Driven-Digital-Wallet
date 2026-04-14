package com.example.balance.service.impl;

import com.example.balance.dto.TransactionalCreatedMessage;
import com.example.balance.dto.TransactionalFailedMessage;
import com.example.balance.dto.TransactionalProcessedMessage;
import com.example.balance.exception.OperationNotSupportedException;
import com.example.balance.mapper.MessageMapper;
import com.example.balance.model.Reserved;
import com.example.balance.service.BalanceService;
import com.example.balance.service.ReservedService;
import com.example.balance.service.TransactionService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class TransactionServiceImpl implements TransactionService {

    private final BalanceService balanceService;
    private final ReservedService reservedService;
    private final OutboxServiceImpl outboxService;
    private final MessageMapper messageMapper;

    @Transactional
    public void handleCreated(TransactionalCreatedMessage message) {
        if (!reservedService.isIdempotentRequest(message.getIdempotent())) return;
        switch (message.getType()) {
            case TRANSFER -> {
                withdraw(message.getAccountFrom(), message.getAmount(), message.getIdempotent());
                deposit(message.getAccountTo(), message.getAmount(), message.getIdempotent());
            }
            case DEPOSIT -> deposit(message.getAccountTo(), message.getAmount(), message.getIdempotent());
            case WITHDRAW -> withdraw(message.getAccountFrom(), message.getAmount(), message.getIdempotent());
            case null, default -> throw new OperationNotSupportedException();
        }
        outboxService.createTransactionalReservedMessage(messageMapper.toReserved(message));
    }

    private void withdraw(UUID accountFrom, BigDecimal amount, UUID idempotent) {
        reservedService.create(
                Reserved.builder()
                        .account(balanceService.deposit(accountFrom, amount.negate()))
                        .amount(amount)
                        .idempotent(idempotent)
                        .build()
        );
    }

    private void deposit(UUID accountTo, BigDecimal amount, UUID idempotent) {
        reservedService.create(
                Reserved.builder()
                        .account(balanceService.deposit(accountTo, amount))
                        .amount(amount)
                        .idempotent(idempotent)
                        .build()
        );
    }

    @Transactional
    public void handleProcessed(TransactionalProcessedMessage message) {
        reservedService.confirm(message.getIdempotent());
        outboxService.createTransactionalChangedMessage(messageMapper.toChanged(message));
    }

    @Transactional
    public void handleFailed(TransactionalFailedMessage message) {
        var reserved = reservedService.readByIdempotent(message.getIdempotent());
        reserved.forEach(item -> {
            balanceService.deposit(item.getAccount().getAccount(), item.getAmount());
        });
        reservedService.removeAll(reserved);
    }
}
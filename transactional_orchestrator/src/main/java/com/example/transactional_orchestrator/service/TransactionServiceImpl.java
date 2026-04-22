package com.example.transactional_orchestrator.service;

import com.example.transactional_orchestrator.client.AccountClient;
import com.example.transactional_orchestrator.dto.*;
import com.example.transactional_orchestrator.exception.NotFoundTransactionException;
import com.example.transactional_orchestrator.mapper.TransactionMapper;
import com.example.transactional_orchestrator.model.Transaction;
import com.example.transactional_orchestrator.repository.TransactionRepository;
import com.example.transactional_orchestrator.util.AuthUtils;
import lombok.RequiredArgsConstructor;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;

import static com.example.transactional_orchestrator.dto.EventType.CHANGED;
import static com.example.transactional_orchestrator.model.Status.FAILED;
import static java.lang.Boolean.FALSE;


@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class TransactionServiceImpl implements TransactionService {

    private final TransactionMapper mapper;
    private final TransactionRepository repo;
    private final OutboxService outboxService;
    private final AccountClient accountClient;


    @Override
    @Transactional
    public TransactionResponse createTransaction(TransactionRequest request) {
        var optionalTransaction = repo.findByIdempotent(request.getIdempotent());
        if (optionalTransaction.isPresent()) {
            return mapper.toDTO(optionalTransaction.get());
        }
        return mapper.toDTO(create(request));
    }

    @Override
    public TransactionResponse getTransactionInfo(Long transactionalId) {
        var transaction = getById(transactionalId);
        throwExceptionIfNotOwner(transaction.getOwnerId());
        return mapper.toDTO(transaction);
    }

    @Override
    @Transactional
    public void changeStatus(TransactionMessage message) {
        Transaction transaction = getById(message.getId());
        transaction.setStatus(mapper.toStatus(message.getEventType()));
        if (CHANGED.equals(message.getEventType())) {
            sendNotification(transaction);
            transaction.setExecutedAt(Instant.now());
        }
        repo.save(transaction);
    }

    @Override
    @Transactional
    public void setFailedStatus(TransactionalFailedMessage message) {
        Transaction transaction = getById(message.getId());
        transaction.setStatus(FAILED);
        transaction.setMessageError(message.getMessageError());
        repo.save(transaction);
        sendNotification(transaction);
    }

    private Transaction create(TransactionRequest request) {
        throwExceptionIfInvalidAccount(request);
        var transaction = mapper.toModel(request);
        transaction.setOwnerId(AuthUtils.getUserId());
        transaction = repo.save(transaction);
        sendTransactionCreated(transaction);
        return transaction;
    }

    private Transaction getById(Long transactionalId) {
        return repo.findById(transactionalId)
                .orElseThrow(() ->
                        new NotFoundTransactionException(transactionalId)
                );
    }

    private void sendNotification(Transaction transaction) {
        outboxService.createNotificationMessage(mapper.toNotificationMessage(transaction));
    }

    private void sendTransactionCreated(Transaction transaction) {
        var message = mapper.toTransactionMessage(transaction);
        outboxService.sendTransactionCreatedMessage(message);
    }

    private void throwExceptionIfInvalidAccount(TransactionRequest request) {
        AccountResponse validAccount = accountClient.isValidAccount(
                new AccountRequest(
                        request.getAccountFrom(),
                        request.getAccountTo()
                )
        );
        if (FALSE.equals(validAccount.getIsValidAccountFrom()) || FALSE.equals(validAccount.getIsValidAccountTo())) {
            throw new AccessDeniedException("Не валидные счета");
        }
    }

    private void throwExceptionIfNotOwner(Long ownerId) {
        if (!ownerId.equals(AuthUtils.getUserId()))
            throw new AccessDeniedException("Нет доступа к этим данным");
    }
}

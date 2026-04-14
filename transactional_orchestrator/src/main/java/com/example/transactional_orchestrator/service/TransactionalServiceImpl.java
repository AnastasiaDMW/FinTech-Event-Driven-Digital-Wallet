package com.example.transactional_orchestrator.service;

import com.example.transactional_orchestrator.client.AccountClient;
import com.example.transactional_orchestrator.dto.TransactionRequest;
import com.example.transactional_orchestrator.dto.TransactionResponse;
import com.example.transactional_orchestrator.exception.NotFoundAccountException;
import com.example.transactional_orchestrator.exception.NotFoundTransactionException;
import com.example.transactional_orchestrator.mapper.TransactionMapper;
import com.example.transactional_orchestrator.model.Transaction;
import com.example.transactional_orchestrator.repository.TransactionRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import static java.util.Objects.nonNull;


@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class TransactionalServiceImpl implements TransactionalService {

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

    private Transaction create(TransactionRequest request) {
        throwExceptionIfInvalidAccount(request);
        var transaction = mapper.toModel(request);
        transaction.setOwnerId(0L); //TODO: ЗАПОЛНИТЬ ПОЛЬЗОВАТЕЛЕМ ПОКА ЧТО НЕТ ИНФЫ
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
        var message = mapper.toNotificationMessage(transaction);
        message.setUserId(0L); //TODO: ЗАПОЛНИТЬ ПОЛЬЗОВАТЕЛЕМ ПОКА ЧТО НЕТ ИНФЫ
        outboxService.createNotificationMessage(message);
    }

    private void sendTransactionCreated(Transaction transaction) {
        var message = mapper.toTransactionCreatedMessage(transaction);
        outboxService.createTransactionalCreatedMessage(message);
    }

    private void throwExceptionIfInvalidAccount(TransactionRequest request) {
        if (nonNull(request.getToAccountId()) && !accountClient.isValidAccount(request.getToAccountId()))
            throw new NotFoundAccountException(request.getToAccountId());
        if (nonNull(request.getToAccountId()) && !accountClient.isValidAccount(request.getFromAccountId()))
            throw new NotFoundAccountException(request.getFromAccountId());
    }

    private void throwExceptionIfNotOwner(Long ownerId) {
//        if(!Long.getLong("0").equals(ownerId))
//            throw new AccessDeniedException("Нет доступа к этим данным");
    }
}

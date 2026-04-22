package com.example.transactional_orchestrator.mapper;

import com.example.transactional_orchestrator.dto.*;
import com.example.transactional_orchestrator.model.Status;
import com.example.transactional_orchestrator.model.Transaction;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface TransactionMapper {

    Transaction toModel(TransactionRequest request);

    TransactionResponse toDTO(Transaction model);

    @Mapping(source = "ownerId", target = "userId")
    NotificationMessage toNotificationMessage(Transaction model);

    TransactionMessage toTransactionMessage(Transaction model);

    AccountRequest toAccountRequest(TransactionRequest request);

    Status toStatus(EventType eventType);
}

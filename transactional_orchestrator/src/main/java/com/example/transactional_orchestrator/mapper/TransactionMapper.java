package com.example.transactional_orchestrator.mapper;

import com.example.transactional_orchestrator.dto.*;
import com.example.transactional_orchestrator.model.Transaction;
import org.mapstruct.Mapper;
import org.springframework.web.bind.annotation.Mapping;

@Mapper(componentModel = "spring")
public interface TransactionMapper {

    Transaction toModel(TransactionRequest request);

    TransactionResponse toDTO(Transaction model);

    NotificationMessage toNotificationMessage(Transaction model);

    TransactionCreatedMessage toTransactionCreatedMessage(Transaction model);

    AccountRequest toAccountRequest(TransactionRequest request);
}

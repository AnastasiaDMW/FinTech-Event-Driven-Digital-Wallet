package com.example.balance.mapper;

import com.example.balance.dto.*;

@Mapper(componentModel = "spring")
public interface MessageMapper {
    @Ma
    TransactionalReservedMessage toReserved(TransactionalCreatedMessage message);

    @Ma
    TransactionalChangedMessage toChanged(TransactionalProcessedMessage message);

    @Ma
    TransactionalFailedMessage toFailed(TransactionalProcessedMessage message);
}

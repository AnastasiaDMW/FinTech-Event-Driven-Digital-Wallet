package com.example.balance.mapper;

import com.example.balance.dto.*;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface MessageMapper {

    TransactionalFailedMessage toFailed(TransactionalMessage message);
}

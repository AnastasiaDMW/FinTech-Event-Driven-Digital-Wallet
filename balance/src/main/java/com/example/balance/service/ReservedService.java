package com.example.balance.service;

import com.example.balance.model.Balance;
import com.example.balance.model.Reserved;

import java.math.BigDecimal;
import java.util.List;
import java.util.UUID;

public interface ReservedService {

    boolean isIdempotentRequest(UUID idempotent);

    List<Reserved> readByIdempotent(UUID idempotent);

    void create(Reserved reserved);

    void confirm(UUID idempotent);

    void removeAll(List<Reserved> reserved);
}

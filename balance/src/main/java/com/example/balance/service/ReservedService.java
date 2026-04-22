package com.example.balance.service;

import com.example.balance.model.Reserved;

import java.util.List;
import java.util.UUID;

public interface ReservedService {

    boolean isNotIdempotentRequest(UUID idempotent);

    List<Reserved> readByIdempotent(UUID idempotent);

    void create(Reserved reserved);

    void confirm(UUID idempotent);

    void removeAll(List<Reserved> reserved);
}

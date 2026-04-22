package com.example.balance.service;

import com.example.balance.dto.BalanceResponse;
import com.example.balance.model.Balance;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.UUID;

public interface BalanceService {

    Balance withdraw(UUID fromAccountId, BigInteger amount);

    Balance deposit(UUID fromAccountId, BigInteger amount);

    BalanceResponse getBalanceResponce(UUID request);

    void checkBalance(UUID account, BigInteger amount);

    @Transactional
    Balance getBalance(UUID account);
}

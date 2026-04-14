package com.example.balance.service;

import com.example.balance.dto.BalanceResponse;
import com.example.balance.model.Balance;

import java.math.BigDecimal;
import java.util.UUID;

public interface BalanceService {

    Balance withdraw(UUID fromAccountId, BigDecimal amount);

    Balance deposit(UUID fromAccountId, BigDecimal amount);

    BalanceResponse getBalance(UUID request);
}

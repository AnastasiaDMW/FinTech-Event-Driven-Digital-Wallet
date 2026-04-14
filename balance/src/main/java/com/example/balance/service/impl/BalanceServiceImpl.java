package com.example.balance.service.impl;

import com.example.balance.dto.BalanceResponse;
import com.example.balance.exception.InsufficientFundsException;
import com.example.balance.model.Balance;
import com.example.balance.repository.BalanceRepository;
import com.example.balance.service.BalanceService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class BalanceServiceImpl implements BalanceService {

    private final BalanceRepository repo;

    @Override
    @Transactional
    public Balance withdraw(UUID fromAccountId, BigDecimal amount) {
        return reserve(fromAccountId, amount.negate());
    }

    @Override
    @Transactional
    public Balance deposit(UUID fromAccountId, BigDecimal amount) {
        return reserve(fromAccountId, amount);
    }

    @Override
    @Transactional
    public BalanceResponse getBalance(UUID account) {
        Balance balance = getByAccount(account);
        return new BalanceResponse(balance.getBalance());
    }

    private Balance getByAccount(UUID account) {
        return repo.findByAccount(account).orElseGet(() -> create(account));
    }

    private Balance create(UUID account) {
        return repo.save(
                Balance.builder()
                        .account(account)
                        .build()
        );
    }

    private Balance reserve(UUID account, BigDecimal amount) {
        Balance balance = getByAccount(account);
        balance.setBalance(balance.getBalance().subtract(amount));
        throwExceptionIfInsufficientFunds(balance);
        return repo.save(balance);
    }

    private void throwExceptionIfInsufficientFunds(Balance balance) {
        if (balance.getBalance().compareTo(BigDecimal.ZERO) < 0) {
            throw new InsufficientFundsException(balance.getAccount());
        }
    }
}

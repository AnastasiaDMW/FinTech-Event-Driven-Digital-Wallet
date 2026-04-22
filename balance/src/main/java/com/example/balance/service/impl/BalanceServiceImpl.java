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
import java.math.BigInteger;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class BalanceServiceImpl implements BalanceService {

    private final BalanceRepository repo;

    @Override
    @Transactional
    public Balance withdraw(UUID fromAccountId, BigInteger amount) {
        return reserve(fromAccountId, amount.negate());
    }

    @Override
    @Transactional
    public Balance deposit(UUID fromAccountId, BigInteger amount) {
        return reserve(fromAccountId, amount);
    }

    @Override
    @Transactional
    public BalanceResponse getBalanceResponce(UUID account) {
        Balance balance = getByAccount(account);
        return new BalanceResponse(balance.getBalance());
    }

    @Override
    public void checkBalance(UUID account, BigInteger amount) {
        Balance balance = getByAccount(account);
        if (balance.getBalance().add(amount).compareTo(BigInteger.ZERO) < 0) {
            throw new InsufficientFundsException(balance.getAccount());
        }
    }

    @Override
    @Transactional
    public Balance getBalance(UUID account) {
        return getByAccount(account);
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

    private Balance reserve(UUID account, BigInteger amount) {
        Balance balance = getByAccount(account);
        balance.setBalance(balance.getBalance().add(amount));
        return repo.save(balance);
    }
}

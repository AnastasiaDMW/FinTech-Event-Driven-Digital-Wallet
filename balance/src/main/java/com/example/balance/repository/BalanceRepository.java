package com.example.balance.repository;

import com.example.balance.model.Balance;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface BalanceRepository extends JpaRepository<Balance, Long> {

    Optional<Balance> findByAccount(UUID account);
}

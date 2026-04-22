package com.example.balance.repository;

import com.example.balance.model.Reserved;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface ReservedRepository extends JpaRepository<Reserved, Long> {

    boolean existsByIdempotent(UUID idempotent);

    void removeByIdempotent(UUID idempotent);

    List<Reserved> findByIdempotent(UUID idempotent);
}

package com.example.balance.service.impl;

import com.example.balance.model.Reserved;
import com.example.balance.repository.ReservedRepository;
import com.example.balance.service.ReservedService;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class ReservedServiceImpl implements ReservedService {

    private final ReservedRepository repo;

    @Override
    public boolean isIdempotentRequest(UUID idempotent) {
        return repo.existsByIdempotent(idempotent);
    }

    @Override
    public List<Reserved> readByIdempotent(UUID idempotent) {
        return repo.findByIdempotent(idempotent);
    }

    @Override
    @Transactional
    public void create(Reserved reserved) {
        repo.save(reserved);
    }

    @Override
    @Transactional
    public void confirm(UUID idempotent) {
        repo.removeByIdempotent(idempotent);
    }

    @Override
    @Transactional
    public void removeAll(List<Reserved> reserved) {
        repo.deleteAll(reserved);
    }
}

package com.example.balance.controller;

import com.example.balance.api.BalanceApi;
import com.example.balance.dto.BalanceResponse;
import com.example.balance.service.BalanceService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.UUID;

import static com.example.balance.api.BalanceApi.API_V1_BALANCE;

@RequestMapping(API_V1_BALANCE)
@RestController
@RequiredArgsConstructor
public class BalanceController implements BalanceApi {

    private final BalanceService service;

    @Override
    public BalanceResponse getBalance(UUID request) {
        return service.getBalance(request);
    }
}

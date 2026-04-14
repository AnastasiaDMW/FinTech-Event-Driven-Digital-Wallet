package com.example.balance.api;

import com.example.balance.dto.BalanceResponse;
import jakarta.validation.Valid;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.UUID;

public interface BalanceApi {

    String API_V1 = "/api/v1";
    String API_V1_BALANCE = API_V1 + "/balance";
    String PATH_VARIABLE_ACCOUNT = "/{account}";

    @GetMapping(PATH_VARIABLE_ACCOUNT)
    BalanceResponse getBalance(@PathVariable("account") UUID request);
}

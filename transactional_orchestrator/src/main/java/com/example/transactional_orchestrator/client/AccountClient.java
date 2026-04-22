package com.example.transactional_orchestrator.client;

import com.example.transactional_orchestrator.dto.AccountRequest;
import com.example.transactional_orchestrator.dto.AccountResponse;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

@FeignClient(name = "accountService", url = "${account_service.url}")
public interface AccountClient {

    @PostMapping("/api/v1/accounts/valid")
    AccountResponse isValidAccount(@RequestBody AccountRequest request);
}

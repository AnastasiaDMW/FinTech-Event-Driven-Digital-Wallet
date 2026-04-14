package com.example.transactional_orchestrator.client;

import com.example.transactional_orchestrator.dto.AccountRequest;
import org.springframework.cloud.openfeign.FeignClient;

import java.util.UUID;

@FeignClient(name = "accountService", url = "${account_service.url}")
public interface AccountClient {

    default boolean isValidAccount(UUID account){
        return true;//todo
    }
}

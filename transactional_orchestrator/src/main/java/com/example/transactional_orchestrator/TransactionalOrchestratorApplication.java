package com.example.transactional_orchestrator;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cloud.openfeign.EnableFeignClients;

@EnableFeignClients
@SpringBootApplication
public class TransactionalOrchestratorApplication {

	public static void main(String[] args) {
		SpringApplication.run(TransactionalOrchestratorApplication.class, args);
	}

}

package com.example.transactional_orchestrator.dto;

import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
public class NotificationMessage {

    private Long id;
    private Long amount;
    private String status;
    private Long userId;
    private EventType eventType;
}

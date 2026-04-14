package com.example.balance.model;

import com.fasterxml.jackson.databind.JsonNode;
import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import java.time.Instant;
import java.util.UUID;

@Entity
@Table(
        name = "tbl_outbox",
        indexes = {
                @Index(name = "idx_outbox_event_key", columnList = "event_key"),
                @Index(name = "idx_outbox_category", columnList = "category"),
                @Index(name = "idx_outbox_created_at", columnList = "created_at"),
                @Index(name = "idx_outbox_event_type", columnList = "event_type")
        }
)
@NoArgsConstructor
@Getter
@Setter
@Builder
@AllArgsConstructor
public class Outbox {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    @Column(name = "id", nullable = false)
    private UUID id;

    @Column(name = "event_key", nullable = false)
    private String eventKey;

    @Enumerated(EnumType.STRING)
    @Column(name = "category", nullable = false)
    private Category category;

    @JdbcTypeCode(SqlTypes.JSON)
    @Column(name = "payload", columnDefinition = "jsonb", nullable = false)
    private JsonNode payload;

    @Column(name = "created_at", updatable = false)
    private Instant createdAt;

    @PrePersist
    public void prePersist() {
        createdAt = Instant.now();
    }
}

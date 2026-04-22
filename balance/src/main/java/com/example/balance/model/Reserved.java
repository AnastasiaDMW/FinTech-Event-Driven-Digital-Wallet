package com.example.balance.model;

import jakarta.persistence.*;
import lombok.*;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.UUID;

@Entity
@Table(
        name = "tbl_reserved",
        indexes = {
                @Index(name = "idx_reserved_idempotent_account_id", columnList = "idempotent, account_id")
        }
)
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class Reserved {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(optional = false)
    @JoinColumn(name = "account_id", nullable = false)
    private Balance account;

    @Column(name = "amount", nullable = false, precision = 19, scale = 2)
    private BigInteger amount;

    @Column(name = "idempotent", nullable = false, updatable = false)
    private UUID idempotent;
}
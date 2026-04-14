package com.example.gateway.property;

import lombok.Getter;
import lombok.Setter;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Getter
@Setter
@Component
@ConfigurationProperties("gateway")
public class GatewayProperty {

    private int ttl;
    private String issuer;
}

package com.example.gateway.config;

import org.springframework.cloud.gateway.filter.ratelimit.KeyResolver;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import reactor.core.publisher.Mono;

@Configuration
public class ResolverConfig {

    @Bean
    public KeyResolver smartKeyResolver() {
        return exchange ->
                exchange.getPrincipal()
                        .filter(p -> p instanceof JwtAuthenticationToken)
                        .cast(JwtAuthenticationToken.class)
                        .map(t -> "USER:" + t.getToken().getClaimAsString("userId"))
                        .switchIfEmpty(Mono.fromSupplier(() -> {
                            var ip = exchange.getRequest().getRemoteAddress();
                            return "IP:" + (ip != null ? ip.getAddress().getHostAddress() : "unknown");
                        }));
    }
}
package com.example.gateway.config;

import org.springframework.context.annotation.Bean;
import org.springframework.security.config.annotation.web.reactive.EnableWebFluxSecurity;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;

@EnableWebFluxSecurity
public class SecurityConfig {

    @Bean
    public SecurityWebFilterChain securityWebFilterChain(ServerHttpSecurity http) {
        http.authorizeExchange(exchanges ->
                        exchanges.pathMatchers("/login", "/refresh", "/registration")
                                .permitAll()
                                .anyExchange()
                                .authenticated()
                ).oauth2ResourceServer(oauth2 -> oauth2.jwt(jwtSpec -> {}));

        return http.build();
    }
}

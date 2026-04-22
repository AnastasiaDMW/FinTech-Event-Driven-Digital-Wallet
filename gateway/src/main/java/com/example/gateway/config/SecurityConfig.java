package com.example.gateway.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.reactive.EnableWebFluxSecurity;
import org.springframework.security.config.web.server.ServerHttpSecurity;
import org.springframework.security.web.server.SecurityWebFilterChain;

@Configuration
@EnableWebFluxSecurity
public class SecurityConfig {

    @Bean
    public SecurityWebFilterChain securityWebFilterChain(ServerHttpSecurity http) {
        http
                .csrf(ServerHttpSecurity.CsrfSpec::disable)
                .authorizeExchange(exchanges ->
                        exchanges.pathMatchers("/auth/api/v1/login", "/auth/api/v1/signup", "/auth/api/v1/refresh", "/auth/api/v1/logout")
                                .permitAll()
                                .anyExchange()
                                .authenticated()
                ).oauth2ResourceServer(oauth2 -> oauth2.jwt(jwtSpec -> {
                }));
        return http.build();
    }
}

package com.example.gateway.fileter;

import com.example.gateway.property.GatewayProperty;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import org.springframework.security.oauth2.jwt.JwsHeader;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.jwt.JwtClaimsSet;
import org.springframework.security.oauth2.jwt.JwtEncoder;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import org.springframework.web.server.WebFilter;
import org.springframework.web.server.WebFilterChain;
import reactor.core.publisher.Mono;

import java.time.Instant;

import static org.springframework.security.oauth2.jwt.JwtEncoderParameters.from;

@Component
@RequiredArgsConstructor
public class JwtAuthenticationFilter implements WebFilter {

    private final JwtEncoder jwtEncoder;
    private final GatewayProperty gatewayProperty;

    @Override
    @NonNull
    public Mono<Void> filter(ServerWebExchange exchange, @NonNull WebFilterChain chain) {
        String path = exchange.getRequest().getURI().getPath();
        if (path.equals("/auth/api/v1/login") ||
                path.equals("/auth/api/v1/signup") ||
                path.equals("/auth/api/v1/refresh") ||
                path.equals("/auth/api/v1/logout")) {

            return chain.filter(exchange);
        }
        return exchange.getPrincipal()
                .flatMap(principal -> {
                    if (!(principal instanceof JwtAuthenticationToken jwt)) {
                        return Mono.error(new RuntimeException("Unauthorized"));
                    }
                    String token = generateInternalToken(jwt.getToken());
                    ServerWebExchange mutatedExchange = exchange.mutate()
                            .request(exchange.getRequest()
                                    .mutate()
                                    .header("Authorization", "Bearer " + token)
                                    .build())
                            .build();
                    return chain.filter(mutatedExchange);
                });
    }

    private String generateInternalToken(Jwt jwt) {
        Instant now = Instant.now();
        JwtClaimsSet claims = JwtClaimsSet.builder()
                .issuer(gatewayProperty.getIssuer())
                .issuedAt(now)
                .expiresAt(now.plusSeconds(gatewayProperty.getTtl()))
                .claim("tokenType", jwt.getClaim("tokenType"))
                .claim("scope", "internal")
                .claim("userId", jwt.getClaim("userId"))
                .claim("email", jwt.getClaim("email"))
                .claim("role", jwt.getClaim("role"))
                .build();
        JwsHeader header = JwsHeader.with(() -> "RS256").type("JWT").build();
        return jwtEncoder.encode(
                from(header, claims)
        ).getTokenValue();
    }
}
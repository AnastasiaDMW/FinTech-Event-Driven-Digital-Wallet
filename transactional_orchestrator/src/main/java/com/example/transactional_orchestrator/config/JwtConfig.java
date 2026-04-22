package com.example.transactional_orchestrator.config;

import lombok.SneakyThrows;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.oauth2.jwt.JwtDecoder;
import org.springframework.security.oauth2.jwt.NimbusJwtDecoder;

import java.security.KeyFactory;
import java.security.NoSuchAlgorithmException;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;

@Configuration
public class JwtConfig {

    @Bean
    public JwtDecoder jwtDecoder(@Value("${gateway.public-key}") String publicKey) {
        return NimbusJwtDecoder.withPublicKey(buildPublicKey(publicKey)).build();
    }

    @SneakyThrows({NoSuchAlgorithmException.class, InvalidKeySpecException.class})
    private RSAPublicKey buildPublicKey(String publicKey) {
        byte[] decode = Base64.getDecoder().decode(publicKey);
        return (RSAPublicKey) KeyFactory.getInstance("RSA")
                .generatePublic(new X509EncodedKeySpec(decode));
    }
}
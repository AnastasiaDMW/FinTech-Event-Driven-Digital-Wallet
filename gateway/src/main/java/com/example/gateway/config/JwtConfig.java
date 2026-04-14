package com.example.gateway.config;

import com.nimbusds.jose.jwk.JWK;
import com.nimbusds.jose.jwk.JWKSet;
import com.nimbusds.jose.jwk.RSAKey;
import com.nimbusds.jose.jwk.source.ImmutableJWKSet;
import com.nimbusds.jose.jwk.source.JWKSource;
import com.nimbusds.jose.proc.SecurityContext;
import lombok.SneakyThrows;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.oauth2.jwt.JwtEncoder;
import org.springframework.security.oauth2.jwt.NimbusJwtEncoder;
import org.springframework.security.oauth2.jwt.NimbusReactiveJwtDecoder;
import org.springframework.security.oauth2.jwt.ReactiveJwtDecoder;

import java.security.KeyFactory;
import java.security.NoSuchAlgorithmException;
import java.security.interfaces.RSAPrivateKey;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;

@Configuration
public class JwtConfig {

    @Bean
    public ReactiveJwtDecoder jwtDecoder(@Value("${auth.public-key}") String publicKey) {
        return NimbusReactiveJwtDecoder.withPublicKey(buildPublicKey(publicKey)).build();
    }

    @Bean
    public JwtEncoder jwtEncoder(@Value("${gateway.private-key}") String privateKey,
                                 @Value("${gateway.public-key}") String publicKey) {
        JWK jwk = new RSAKey.Builder(buildPublicKey(publicKey))
                .privateKey(buildPrivateKey(privateKey))
                .build();
        JWKSource<SecurityContext> jwkSource = new ImmutableJWKSet<>(new JWKSet(jwk));
        return new NimbusJwtEncoder(jwkSource);
    }

    @SneakyThrows({NoSuchAlgorithmException.class, InvalidKeySpecException.class})
    private RSAPrivateKey buildPrivateKey(String privateKey) {
        byte[] decode = Base64.getDecoder().decode(privateKey);
        return (RSAPrivateKey) KeyFactory.getInstance("RSA")
                .generatePrivate(new PKCS8EncodedKeySpec(decode));
    }

    @SneakyThrows({NoSuchAlgorithmException.class, InvalidKeySpecException.class})
    private RSAPublicKey buildPublicKey(String publicKey) {
        byte[] decode = Base64.getDecoder().decode(publicKey);
        return (RSAPublicKey) KeyFactory.getInstance("RSA")
                .generatePublic(new X509EncodedKeySpec(decode));
    }
}
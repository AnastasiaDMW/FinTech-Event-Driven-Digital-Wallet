package com.example.transactional_orchestrator.util;

import lombok.experimental.UtilityClass;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.jwt.Jwt;

import static java.lang.Long.valueOf;

@UtilityClass
public class AuthUtils {

    public static Long getUserId() {
        var auth = (Jwt) SecurityContextHolder.getContext()
                .getAuthentication()
                .getPrincipal();
        return valueOf(auth.getClaim("userId"));
    }
}

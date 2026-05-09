package com.onepick.idem

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.context.annotation.Configuration
import org.springframework.web.servlet.config.annotation.InterceptorRegistry
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer

@SpringBootApplication
class IdempotencyApplication

@Configuration
class WebConfig(
    private val interceptor: IdempotencyInterceptor,
) : WebMvcConfigurer {
    override fun addInterceptors(registry: InterceptorRegistry) {
        registry.addInterceptor(interceptor)
            .addPathPatterns("/payments/**")  // 결제 경로에만 적용
    }
}

fun main(args: Array<String>) {
    runApplication<IdempotencyApplication>(*args)
}

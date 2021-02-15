package com.example.mykotlinwebflux.router

import com.example.mykotlinwebflux.domain.FooRepository
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.web.reactive.function.server.router

@Configuration
class RouterConfiguration {

	@Bean
	fun router(repository: FooRepository) = router {
		"/router".nest {
			GET("/foo") { ok().syncBody(repository.foo()) }
		}
	}
}
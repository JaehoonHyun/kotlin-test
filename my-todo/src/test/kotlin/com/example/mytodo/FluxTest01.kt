package com.example.mytodo

import io.kotlintest.should
import io.kotlintest.specs.FunSpec
import org.assertj.core.api.Assertions
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.web.reactive.WebFluxTest
import org.springframework.test.web.reactive.server.WebTestClient
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import reactor.core.publisher.Signal.subscribe
import reactor.core.scheduler.Schedulers
import reactor.test.StepVerifier
import java.lang.RuntimeException


class FluxTest01:FunSpec({

    test("mono") {
        Mono.just(1)
            .subscribeOn(Schedulers.single())
            .subscribe{item -> Assertions.assertThat(item).isEqualTo(1)}
    }

    //실패해야하는데 성공...
    test("monoFail") {
        Mono.just(1)
                .subscribeOn(Schedulers.single())
                .subscribe{item -> Assertions.assertThat(item).isEqualTo(2)}
    }

    test("stepVerify") {
        val mono = Mono.just(1).subscribeOn(Schedulers.single())
        StepVerifier.create(mono)
                    .expectNext(2)
                    .verifyComplete()

    }

    test("stepVerfiyFlux") {
        val flux = Flux.just(1,2,3).subscribeOn(Schedulers.single())
        StepVerifier.create(flux)
                .expectNext(1 )
                .expectNext(2)
                .expectNext(3)
                .verifyComplete()
    }

    test("stepVerfiyFluxError") {
        val flux = Flux.just(1,2,3).subscribeOn(Schedulers.single())
        val expectError = StepVerifier.create(flux)
                .expectNext(1)
                .expectNext(2)
                .expectNext(3)
                .expectError()
                .verify()
    }



    test("map") {
        Flux.just("a", "bc", "def", "wxyz")
                .map { str -> str.length }
                .subscribe{ len -> println(len)}
    }

})
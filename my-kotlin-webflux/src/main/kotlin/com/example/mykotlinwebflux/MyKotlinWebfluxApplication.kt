package com.example.mykotlinwebflux

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class MyKotlinWebfluxApplication

fun main(args: Array<String>) {
	runApplication<MyKotlinWebfluxApplication>(*args)
}

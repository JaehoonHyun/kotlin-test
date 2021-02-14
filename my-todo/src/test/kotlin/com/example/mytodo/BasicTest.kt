package com.example.mytodo

import io.kotlintest.shouldBe
import io.kotlintest.specs.StringSpec

class BasicTest: StringSpec() {
    init {
        "strings.length should return size of string" {
            "hello".length shouldBe 5
        }
    }

}
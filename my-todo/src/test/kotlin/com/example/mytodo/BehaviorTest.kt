package com.example.mytodo

import io.kotlintest.shouldBe
import io.kotlintest.specs.BehaviorSpec

class BehaviorTest : BehaviorSpec(
    {
        given("a boomstick"){
            `when`("I sit on it") {
                then("should be on the fly") {
                    "hello".length shouldBe 5
                }
            }
        }
    }
)
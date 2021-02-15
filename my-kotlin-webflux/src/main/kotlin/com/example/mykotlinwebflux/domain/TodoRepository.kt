package com.example.mytodo.domain

import com.example.mykotlinwebflux.domain.Todo
import org.springframework.data.jpa.repository.JpaRepository

interface TodoRepository : JpaRepository<Todo, Long>


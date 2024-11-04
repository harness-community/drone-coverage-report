package com.example.package3;

public class Divider {
    public int divide(int a, int b) {
        if (b == 0) {
            throw new ArithmeticException("Division by zero is not allowed");
        } else if (b < 0) {
            throw new IllegalArgumentException("Negative divisor is not supported");
        }
        return a / b;
    }
}

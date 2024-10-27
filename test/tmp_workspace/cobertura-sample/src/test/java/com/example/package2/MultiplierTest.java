package com.example.package2;

import org.junit.Test;
import static org.junit.Assert.*;

public class MultiplierTest {
    @Test
    public void testMultiply() {
        Multiplier multiplier = new Multiplier();
        assertEquals(12, multiplier.multiply(3, 4));
    }
}

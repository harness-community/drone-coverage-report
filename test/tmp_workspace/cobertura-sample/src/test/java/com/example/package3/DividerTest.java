package com.example.package3;

import org.junit.Test;
import static org.junit.Assert.*;

public class DividerTest {
    @Test
    public void testDivide() {
        Divider divider = new Divider();
        assertEquals(2, divider.divide(8, 4));
    }
}

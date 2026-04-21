package works.weave.socks.orders.controllers;

import org.junit.Test;
import works.weave.socks.orders.entities.Item;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

public class UnitOrdersController {

    private final OrdersController controller = new OrdersController();

    @Test
    public void calculateTotalIncludesShipping() {
        Item item = new Item(null, "item1", 2, 10.0F);
        float total = controller.calculateTotal(Collections.singletonList(item));
        assertEquals(24.99F, total, 0.001F);
    }

    @Test
    public void calculateTotalRejectsNegativeUnitPrice() {
        Item item = new Item(null, "item1", 1, -5.0F);
        try {
            controller.calculateTotal(Collections.singletonList(item));
            fail("Expected InvalidOrderException for negative unit price");
        } catch (OrdersController.InvalidOrderException e) {
            // expected
        }
    }

    @Test
    public void calculateTotalRejectsMixedNegativeAndPositivePrices() {
        List<Item> items = Arrays.asList(
                new Item(null, "item1", 2, 10.0F),
                new Item(null, "item2", 1, -20.0F)
        );
        try {
            controller.calculateTotal(items);
            fail("Expected InvalidOrderException for item with negative unit price");
        } catch (OrdersController.InvalidOrderException e) {
            // expected
        }
    }

    @Test
    public void calculateTotalAllowsZeroUnitPrice() {
        Item item = new Item(null, "item1", 3, 0.0F);
        float total = controller.calculateTotal(Collections.singletonList(item));
        assertEquals(4.99F, total, 0.001F);
    }
}

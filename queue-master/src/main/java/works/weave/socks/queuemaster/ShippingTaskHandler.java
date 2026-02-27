package works.weave.socks.queuemaster;

import org.springframework.stereotype.Component;
import works.weave.socks.shipping.entities.Shipment;

@Component
public class ShippingTaskHandler {

	public void handleMessage(Shipment shipment) {
		System.out.println("Received shipment task: " + shipment.getName());
	}
}

package works.weave.socks.queuemaster;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.springframework.stereotype.Component;
import works.weave.socks.shipping.entities.Shipment;

@Component
public class ShippingTaskHandler {

	private static final Logger log = LoggerFactory.getLogger(ShippingTaskHandler.class);

	public void handleMessage(Shipment shipment) {
		String tid = MDC.get("traceId");
		String sid = MDC.get("spanId");
		if (tid != null && !tid.isEmpty() && sid != null && !sid.isEmpty()) {
			log.info("Received shipment task: {} | traceId={} spanId={}", shipment.getName(), tid, sid);
		} else {
			log.info("Received shipment task: {}", shipment.getName());
		}
	}
}

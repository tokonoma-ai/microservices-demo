package works.weave.socks.queuemaster.configuration;

import org.junit.jupiter.api.Test;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.core.MessageProperties;
import org.springframework.amqp.rabbit.listener.SimpleMessageListenerContainer;
import org.springframework.amqp.support.converter.MessageConverter;
import org.springframework.test.util.ReflectionTestUtils;
import works.weave.socks.queuemaster.ShippingTaskHandler;
import works.weave.socks.shipping.entities.Shipment;

import java.nio.charset.StandardCharsets;

import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.CoreMatchers.instanceOf;
import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.MatcherAssert.assertThat;

/**
 * Plain unit tests (no Spring context, no broker) for the way queue-master
 * converts the messages the shipping service publishes to "shipping-task".
 */
public class UnitShippingTaskMessageConversion {

    // Type id the shipping service's Jackson2JsonMessageConverter stamps on every message.
    private static final String SHIPPING_TYPE_ID = "works.weave.socks.shipping.entities.Shipment";
    private static final String ID = "59b4d895-7c1a-4f0e-9b2d-3e4f5a6b7c8d";
    private static final String NAME = "6ae3ae051291e6aef482e07f";
    private static final String BODY = "{\"id\":\"" + ID + "\",\"name\":\"" + NAME + "\"}";

    private final ShippingConsumerConfiguration configuration = new ShippingConsumerConfiguration();
    private final MessageConverter converter = configuration.jsonMessageConverter();

    private static Message jsonMessage(String typeId) {
        MessageProperties properties = new MessageProperties();
        properties.setContentType(MessageProperties.CONTENT_TYPE_JSON);
        properties.setContentEncoding(StandardCharsets.UTF_8.name());
        if (typeId != null) {
            properties.setHeader("__TypeId__", typeId);
        }
        return new Message(BODY.getBytes(StandardCharsets.UTF_8), properties);
    }

    @Test
    public void convertsMessageStampedWithShippingServiceTypeId() {
        Object converted = converter.fromMessage(jsonMessage(SHIPPING_TYPE_ID));

        assertThat(converted, is(instanceOf(Shipment.class)));
        Shipment shipment = (Shipment) converted;
        assertThat(shipment.getId(), is(equalTo(ID)));
        assertThat(shipment.getName(), is(equalTo(NAME)));
    }

    @Test
    public void convertsMessageWithoutTypeIdUsingDefaultType() {
        Object converted = converter.fromMessage(jsonMessage(null));

        assertThat(converted, is(instanceOf(Shipment.class)));
        Shipment shipment = (Shipment) converted;
        assertThat(shipment.getId(), is(equalTo(ID)));
        assertThat(shipment.getName(), is(equalTo(NAME)));
    }

    @Test
    public void listenerContainerDoesNotRequeueRejectedMessages() {
        ReflectionTestUtils.setField(configuration, "shippingTaskHandler", new ShippingTaskHandler());

        SimpleMessageListenerContainer container = configuration.listenerContainer();

        Boolean requeueRejected = ReflectionTestUtils.invokeMethod(container, "isDefaultRequeueRejected");
        assertThat(requeueRejected, is(false));
    }
}

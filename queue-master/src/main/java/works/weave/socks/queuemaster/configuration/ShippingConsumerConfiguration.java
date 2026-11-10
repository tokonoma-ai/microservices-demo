package works.weave.socks.queuemaster.configuration;

import org.springframework.amqp.core.Queue;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.amqp.rabbit.listener.SimpleMessageListenerContainer;
import org.springframework.amqp.rabbit.listener.adapter.MessageListenerAdapter;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import works.weave.socks.queuemaster.ShippingTaskHandler;

@Configuration
public class ShippingConsumerConfiguration extends RabbitMqConfiguration
{
	protected final String queueName = "shipping-task";

    @Autowired
    private ShippingTaskHandler shippingTaskHandler;

	@Bean
	public RabbitTemplate rabbitTemplate() {
		RabbitTemplate template = new RabbitTemplate(connectionFactory());
		template.setDefaultReceiveQueue(this.queueName);
        template.setMessageConverter(jsonMessageConverter());
		return template;
	}

    @Bean
	public Queue queueName() {
		return new Queue(this.queueName, false);
	}

	@Bean
	public SimpleMessageListenerContainer listenerContainer() {
		SimpleMessageListenerContainer container = new SimpleMessageListenerContainer();
		container.setConnectionFactory(connectionFactory());
		container.setQueueNames(this.queueName);
		container.setMessageListener(messageListenerAdapter());
		// There is no dead-letter queue: a message the listener cannot handle is dropped after one
		// logged failure rather than requeued and redelivered forever.
		container.setDefaultRequeueRejected(false);

		return container;
	}

    @Bean
    public MessageListenerAdapter messageListenerAdapter() {
        return new MessageListenerAdapter(shippingTaskHandler, jsonMessageConverter());
    }
}

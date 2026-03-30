package works.weave.socks.orders.controllers;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.data.rest.webmvc.RepositoryRestController;
import org.springframework.hateoas.EntityModel;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.client.HttpServerErrorException;
import works.weave.socks.orders.config.OrdersConfigurationProperties;
import works.weave.socks.orders.entities.*;
import works.weave.socks.orders.repositories.CustomerOrderRepository;
import works.weave.socks.orders.resources.NewOrderResource;
import works.weave.socks.orders.services.AsyncGetService;
import works.weave.socks.orders.values.PaymentRequest;
import works.weave.socks.orders.values.PaymentResponse;

import java.io.IOException;
import java.util.Calendar;
import java.util.List;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;


@RepositoryRestController
public class OrdersController {
    private final Logger LOG = LoggerFactory.getLogger(getClass());

    @Autowired
    private OrdersConfigurationProperties config;

    @Autowired
    private AsyncGetService asyncGetService;

    @Autowired
    private CustomerOrderRepository customerOrderRepository;

    @Value(value = "${http.timeout:5}")
    private long timeout;

    @ResponseStatus(HttpStatus.CREATED)
    @RequestMapping(path = "/orders", consumes = MediaType.APPLICATION_JSON_VALUE, method = RequestMethod.POST)
    public
    @ResponseBody
    CustomerOrder newOrder(@RequestBody NewOrderResource item) {
        try {

            if (item.address == null || item.customer == null || item.card == null || item.items == null) {
                throw new InvalidOrderException("Invalid order request. Order requires customer, address, card and items.");
            }


            LOG.debug("Starting calls");
            Future<EntityModel<Address>> addressFuture = asyncGetService.getResource(item.address, new
                    ParameterizedTypeReference<EntityModel<Address>>() {
            });
            Future<EntityModel<Customer>> customerFuture = asyncGetService.getResource(item.customer, new
                    ParameterizedTypeReference<EntityModel<Customer>>() {
            });
            Future<EntityModel<Card>> cardFuture = asyncGetService.getResource(item.card, new
                    ParameterizedTypeReference<EntityModel<Card>>() {
            });
            Future<List<Item>> itemsFuture = asyncGetService.getDataList(item.items, new
                    ParameterizedTypeReference<List<Item>>() {
            });
            LOG.debug("End of calls.");

            float amount = calculateTotal(itemsFuture.get(timeout, TimeUnit.SECONDS));
            if (amount <= 0) {
                throw new InvalidOrderException(
                        "Order total must be greater than zero; computed amount: " + amount +
                        ". One or more items in the cart may have an invalid price.");
            }

            // Call payment service to make sure they've paid
            PaymentRequest paymentRequest = new PaymentRequest(
                    addressFuture.get(timeout, TimeUnit.SECONDS).getContent(),
                    cardFuture.get(timeout, TimeUnit.SECONDS).getContent(),
                    customerFuture.get(timeout, TimeUnit.SECONDS).getContent(),
                    amount);
            LOG.info("action=new_order customer_url={} amount={} status=payment_requested", item.customer, amount);
            Future<PaymentResponse> paymentFuture = asyncGetService.postResource(
                    config.getPaymentUri(),
                    paymentRequest,
                    new ParameterizedTypeReference<PaymentResponse>() {
                    });
            PaymentResponse paymentResponse = paymentFuture.get(timeout, TimeUnit.SECONDS);
            LOG.info("action=new_order customer_url={} amount={} payment_authorised={} payment_message={}",
                    item.customer, amount, paymentResponse != null ? paymentResponse.isAuthorised() : "null",
                    paymentResponse != null ? paymentResponse.getMessage() : "null");
            if (paymentResponse == null) {
                throw new PaymentDeclinedException("Unable to parse authorisation packet");
            }
            if (!paymentResponse.isAuthorised()) {
                throw new PaymentDeclinedException(paymentResponse.getMessage());
            }

            // Ship
            String customerId = parseId(customerFuture.get(timeout, TimeUnit.SECONDS).getRequiredLink("self").getHref());
            Future<Shipment> shipmentFuture = asyncGetService.postResource(config.getShippingUri(), new Shipment
                    (customerId), new ParameterizedTypeReference<Shipment>() {
            });

            CustomerOrder order = new CustomerOrder(
                    null,
                    customerId,
                    customerFuture.get(timeout, TimeUnit.SECONDS).getContent(),
                    addressFuture.get(timeout, TimeUnit.SECONDS).getContent(),
                    cardFuture.get(timeout, TimeUnit.SECONDS).getContent(),
                    itemsFuture.get(timeout, TimeUnit.SECONDS),
                    shipmentFuture.get(timeout, TimeUnit.SECONDS),
                    Calendar.getInstance().getTime(),
                    amount);
            LOG.debug("Received data: " + order.toString());

            CustomerOrder savedOrder = customerOrderRepository.save(order);
            LOG.info("action=new_order customer_id={} order_id={} amount={} status=created",
                    customerId, savedOrder.getId(), amount);

            return savedOrder;
        } catch (TimeoutException e) {
            LOG.error("action=new_order customer_url={} status=failed err=timeout", item.customer, e);
            throw new IllegalStateException("Unable to create order due to timeout from one of the services.", e);
        } catch (ExecutionException e) {
            Throwable cause = e.getCause();
            if (cause instanceof HttpServerErrorException) {
                HttpServerErrorException httpEx = (HttpServerErrorException) cause;
                String responseBody = httpEx.getResponseBodyAsString();
                LOG.error("action=new_order customer_url={} status=failed err={} payment_response={}",
                        item.customer, httpEx.getMessage(), responseBody);
                throw new PaymentDeclinedException("Payment service error: " + httpEx.getStatusCode() +
                        (responseBody != null && !responseBody.isEmpty() ? " - " + responseBody : ""));
            }
            LOG.error("action=new_order customer_url={} status=failed err={}", item.customer, e.getMessage(), e);
            throw new IllegalStateException("Unable to create order due to unspecified IO error.", e);
        } catch (InterruptedException | IOException e) {
            LOG.error("action=new_order customer_url={} status=failed err={}", item.customer, e.getMessage(), e);
            throw new IllegalStateException("Unable to create order due to unspecified IO error.", e);
        }
    }

    private String parseId(String href) {
        Pattern idPattern = Pattern.compile("[\\w-]+$");
        Matcher matcher = idPattern.matcher(href);
        if (!matcher.find()) {
            throw new IllegalStateException("Could not parse user ID from: " + href);
        }
        return matcher.group(0);
    }

    private float calculateTotal(List<Item> items) {
        float amount = 0F;
        float shipping = 4.99F;
        amount += items.stream().mapToDouble(i -> i.getQuantity() * i.getUnitPrice()).sum();
        amount += shipping;
        return amount;
    }

    @ResponseStatus(value = HttpStatus.NOT_ACCEPTABLE)
    public class PaymentDeclinedException extends IllegalStateException {
        public PaymentDeclinedException(String s) {
            super(s);
        }
    }

    @ResponseStatus(value = HttpStatus.NOT_ACCEPTABLE)
    public class InvalidOrderException extends IllegalStateException {
        public InvalidOrderException(String s) {
            super(s);
        }
    }
}

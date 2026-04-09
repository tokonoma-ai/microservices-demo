package works.weave.socks.orders.configuration;

import brave.Tracing;
import brave.context.slf4j.MDCCurrentTraceContext;
import brave.http.HttpTracing;
import brave.sampler.CountingSampler;
import brave.sampler.Sampler;
import brave.servlet.TracingFilter;
import javax.servlet.Filter;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.Ordered;
import zipkin2.reporter.Sender;
import zipkin2.reporter.brave.AsyncZipkinSpanHandler;
import zipkin2.reporter.okhttp3.OkHttpSender;

@Configuration
@ConditionalOnProperty(name = "management.tracing.enabled", havingValue = "true")
public class BraveTracingMdcConfiguration {

    @Bean(destroyMethod = "close")
    public Sender zipkinSender(@Value("${management.zipkin.tracing.endpoint}") String url) {
        return OkHttpSender.create(url);
    }

    @Bean(destroyMethod = "close")
    public AsyncZipkinSpanHandler zipkinSpanHandler(Sender sender) {
        return AsyncZipkinSpanHandler.create(sender);
    }

    @Bean(destroyMethod = "close")
    public Tracing tracing(
            @Value("${spring.application.name}") String serviceName,
            AsyncZipkinSpanHandler zipkinSpanHandler,
            @Value("${management.tracing.sampling.probability:1.0}") double probability) {
        Sampler sampler =
                probability >= 1.0 ? Sampler.ALWAYS_SAMPLE : CountingSampler.create((float) probability);
        return Tracing.newBuilder()
                .localServiceName(serviceName)
                .addSpanHandler(zipkinSpanHandler)
                .sampler(sampler)
                .currentTraceContext(MDCCurrentTraceContext.create())
                .build();
    }

    @Bean
    public HttpTracing httpTracing(Tracing tracing) {
        return HttpTracing.create(tracing);
    }

    @Bean
    public FilterRegistrationBean<Filter> tracingFilter(HttpTracing httpTracing) {
        FilterRegistrationBean<Filter> reg = new FilterRegistrationBean<>();
        reg.setFilter(TracingFilter.create(httpTracing));
        reg.setOrder(Ordered.HIGHEST_PRECEDENCE + 5);
        reg.addUrlPatterns("/*");
        return reg;
    }
}

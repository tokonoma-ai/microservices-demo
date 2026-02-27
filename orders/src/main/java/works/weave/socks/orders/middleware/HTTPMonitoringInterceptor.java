package works.weave.socks.orders.middleware;

import io.prometheus.client.Histogram;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

public class HTTPMonitoringInterceptor implements HandlerInterceptor {

    private static final Logger log = LoggerFactory.getLogger(HTTPMonitoringInterceptor.class);

    static final Histogram requestLatency = Histogram.build()
            .name("http_request_duration_seconds")
            .help("Request duration in seconds.")
            .labelNames("service", "method", "path", "status_code")
            .register();

    private static final String startTimeKey = "startTime";

    @Value("${spring.application.name:orders}")
    private String serviceName;

    @Override
    public boolean preHandle(HttpServletRequest httpServletRequest, HttpServletResponse
            httpServletResponse, Object o) throws Exception {
        httpServletRequest.setAttribute(startTimeKey, System.nanoTime());
        return true;
    }

    @Override
    public void postHandle(HttpServletRequest httpServletRequest, HttpServletResponse
            httpServletResponse, Object o, ModelAndView modelAndView) throws Exception {
    }

    @Override
    public void afterCompletion(HttpServletRequest httpServletRequest, HttpServletResponse
            httpServletResponse, Object o, Exception e) throws Exception {
        if (httpServletRequest.getAttribute(startTimeKey) == null) {
            return;
        }
        String path = httpServletRequest.getRequestURI();
        if (!path.startsWith("/orders") && !path.equals("/health")) {
            return;
        }
        long start = (long) httpServletRequest.getAttribute(startTimeKey);
        long elapsedNs = System.nanoTime() - start;
        long elapsedMs = elapsedNs / 1_000_000;
        String method = httpServletRequest.getMethod();
        int status = httpServletResponse.getStatus();
        log.info("{} {} {} {}ms", method, path, status, elapsedMs);
        try {
            requestLatency.labels(
                    serviceName,
                    method,
                    path,
                    Integer.toString(status)
            ).observe(elapsedNs / 1_000_000_000.0);
        } catch (Exception ex) {
            log.debug("Could not record request metric: {}", ex.getMessage());
        }
    }
}

FROM eclipse-temurin:17-jre

WORKDIR /usr/src/app
COPY ./target/carts.jar ./app.jar

# Kubernetes sets tracing via JAVA_OPTS (-Dzipkin_*). exec-form ENTRYPOINT must invoke a shell
# so those JVM flags are actually applied (see deploy/kubernetes/manifests/*-dep.yaml).
ENTRYPOINT ["sh", "-c", "exec java ${JAVA_OPTS:-} -Djava.security.egd=file:/dev/urandom -jar ./app.jar --port=80"]

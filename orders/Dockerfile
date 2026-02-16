FROM eclipse-temurin:8-jre

WORKDIR /usr/src/app
COPY ./target/orders.jar ./app.jar

ENTRYPOINT ["java","-Djava.security.egd=file:/dev/urandom","-jar","./app.jar", "--port=80"]

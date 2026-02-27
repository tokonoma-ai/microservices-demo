FROM eclipse-temurin:17-jre

WORKDIR /usr/src/app
COPY ./target/queue-master.jar ./app.jar

ENTRYPOINT ["java","-Djava.security.egd=file:/dev/urandom","-jar","./app.jar", "--port=80"]

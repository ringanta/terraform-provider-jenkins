FROM jenkins/jenkins:lts

ENV CASC_JENKINS_CONFIG /usr/share/jenkins/casc.yaml

RUN /usr/local/bin/install-plugins.sh git matrix-auth configuration-as-code

COPY casc.yaml /usr/share/jenkins/casc.yaml

HEALTHCHECK --interval=4s --start-period=5s --retries=30 CMD [ "curl", "-f", "http://localhost:8080/login" ]
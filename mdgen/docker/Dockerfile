FROM govsvc/cloudhsm-client-test:0.0.1560968513
COPY ./sandbox-customerCA.crt /opt/cloudhsm/etc/customerCA.crt
RUN apt-get update && apt-get install dnsutils -y
RUN /opt/cloudhsm/bin/configure -a $(dig +short a88bb4c07943b11e9bbf30ae9bf7a1ac-aa921f50a00f6c2a.elb.eu-west-2.amazonaws.com)
RUN wget https://s3.amazonaws.com/cloudhsmv2-software/CloudHsmClient/Xenial/cloudhsm-client_latest_amd64.deb
RUN wget https://s3.amazonaws.com/cloudhsmv2-software/CloudHsmClient/Xenial/cloudhsm-client-jce_latest_amd64.deb
RUN apt-get install openjdk-8-jre -y
RUN dpkg -i cloudhsm-client_latest_amd64.deb
RUN dpkg -i cloudhsm-client-jce_latest_amd64.deb
# don't provide HSM_USER=cu9 variable - pass it in with the -e HSM_USER when running the image
# don't provide HSM_PASSWORD variable - pass it in with the -e HSM_PASSWORD when running the image
# when running the image, use -v to bind mount the folder with test java and test results

#CMD /opt/cloudhsm/bin/cloudhsm_client /opt/cloudhsm/etc/cloudhsm_client.cfg > cloudhsm.log
COPY ./start.bash /start.bash
ENTRYPOINT [ "/bin/bash", "/start.bash" ]

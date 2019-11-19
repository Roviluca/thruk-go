FROM consol/omd-labs-debian:latest

RUN htpasswd -b /omd/sites/demo/etc/htpasswd omdadmin omdadmin

ADD ./start.sh /root/start_and_wait.sh

RUN chmod +x /root/start_and_wait.sh

CMD ["/root/start_and_wait.sh"]

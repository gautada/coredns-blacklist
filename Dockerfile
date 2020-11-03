FROM development:build

EXPOSE 53/tcp 53/udp
EXPOSE 9153/tcp

RUN useradd -ms /bin/bash coredns \
 && echo 'coredns:coredns' | chpasswd \
 && usermod -aG sudo coredns 

RUN apt-get install -y --no-install-recommends dnsutils \
 && git clone https://github.com/coredns/coredns.git coredns 

COPY plugin.cfg /coredns/plugin.cfg

WORKDIR /coredns 
                                                                                                                                                        
RUN go generate 

RUN ls 

# COPY blacklist/ /coredns/plugin/blacklist/
# RUN ls -l /coredns/plugin/blacklist/
RUN go build                                                                                                          

RUN chmod u+s /coredns/coredns 

RUN echo "All Plugins:" 

RUN /coredns/coredns -plugins

RUN echo "/coredns/coredns -conf /etc/coredns/Corefile" > /home/coredns/.bash_history                                                                                     
                                                                                           
COPY blacklist.py /home/coredns/blacklist.py

# RUN mkdir /home/coredns/blocked \
#  && cd /home/coredns/blocked \
#  && git clone https://github.com/StevenBlack/hosts.git sb-hosts
# ADD http://mirror1.malwaredomains.com/files/spywaredomains.zones /home/coredns/blocked/malware.zone

RUN chown -R coredns:coredns /home/coredns 

COPY etc/hosts /etc/coredns/hosts
COPY etc/whitelist /etc/coredns/whitelist
COPY etc/blacklist /etc/coredns/blacklist
COPY etc/Corefile /etc/coredns/Corefile

USER coredns

CMD ["/coredns/coredns", "-conf", "/etc/coredns/Corefile"]

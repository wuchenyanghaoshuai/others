FROM 192.168.3.103:8888/base/jre-sky-agent:alpine_v1
COPY opt  CFCA_OCA1.cer CFCA_CS_CA.cer /opt/
COPY  lieyun.cer /usr/lib/jvm/java-1.8-openjdk/jre/lib/
RUN sh -c '/bin/echo -e "yes"|keytool -import -trustcacerts -keystore "/usr/lib/jvm/java-1.8-openjdk/jre/lib/security/cacerts" -storepass changeit -alias nealnet -file  /usr/lib/jvm/java-1.8-openjdk/jre/lib/lieyun.cer' &&  sh -c '/bin/echo -e "yes"|keytool -import -trustcacerts -keystore "/usr/lib/jvm/java-1.8-openjdk/jre/lib/security/cacerts" -storepass changeit -alias nealnet1 -file  /opt/CFCA_CS_CA.cer' &&  sh -c '/bin/echo -e "yes" |keytool -import -trustcacerts -keystore "/usr/lib/jvm/java-1.8-openjdk/jre/lib/security/cacerts" -storepass changeit -alias nealnet2 -file  /opt/CFCA_OCA1.cer'

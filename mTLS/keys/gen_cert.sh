#! /bin/sh

for NAME in server client intermediate root
do
    if [ ! -f $NAME.pem ]; then
        openssl req -new -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 \
            -x509 -nodes -days 365 -out $NAME.cert.pem -keyout $NAME.pem \
            -subj "/C=US/ST=TN/L=Nashville/O=FakeOrg/CN=localhost" \
            -addext "certificatePolicies = 1.2.3.4" \
            -addext "subjectAltName = DNS:localhost, DNS:127.0.0.1, DNS:::1"

        chmod 444 "$NAME.pem"
    fi
done



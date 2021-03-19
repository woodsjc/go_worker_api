#! /bin/sh

if [ ! -f intermediate.signed.cert.pem ]; then
    openssl req -new -sha256 -key intermediate.pem \
        -subj "/C=US/ST=TN/L=Nashville/O=FakeOrg/CN=localhost" \
        -addext "subjectAltName = DNS:localhost" \
        -out intermediate.request_sign.pem

    if [ $? -ne 0 ]; then
        exit $?
    fi

    openssl ca -extensions v3_intermediate_ca -days 365 -notext \
        -md sha256 -in intermediate.request_sign.pem  \
        -batch \
        -out intermediate.signed.cert.pem -config ./openssl.cnf

    if [ $? -ne 0 ]; then
        exit $?
    fi

    chmod 444 "intermediate.signed.cert.pem"
    rm intermediate.request_sign.pem
fi


## Also make chain file
if [ ! -f ca-chain.cert.pem ]; then
    cat intermediate.signed.cert.pem root.cert.pem > ca-chain.cert.pem
    chmod 444 ca-chain.cert.pem
fi

for NAME in server client
do
    if [ ! -f "$NAME.signed.cert.pem" ]; then
        openssl req -new -sha256 -key "$NAME.pem" \
            -subj "/C=US/ST=TN/L=Nashville/O=FakeOrg/CN=localhost" \
            -addext "subjectAltName = DNS:localhost" \
            -out $NAME.request_sign.pem
        
        if [ $? -ne 0 ]; then
            exit $?
        fi

        openssl ca -config intermediate_openssl.cnf  \
            -extensions "${NAME}_cert" -days 365 -notext \
            -batch \
            -md sha256 -in "$NAME.request_sign.pem" -out "$NAME.signed.cert.pem" 

        if [ $? -ne 0 ]; then
            exit $?
        fi

        chmod 444 "$NAME.signed.cert.pem"
        rm "$NAME.request_sign.pem"
    fi
done


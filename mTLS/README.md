# mTLS - with self signed certs

## [keys](keys)

After setting up intermediate\_openssl.cnf and openssl.cnf run ```./gen_cert.sh``` and ```sign_keys.sh```

Ton going on here because have to create a certificate authority. Then make an intermediary. Then Add to chain and make server/client as well as sign all keys.

## [server](server)

Finally working with curl 

Post example:
```
curl -X POST "https://localhost:55555/command/" --cacert ../keys/ca-chain.cert.pem --cert ../keys/client.signed.cert.pem --key ../keys/client.pem -d '{"Name":"ls","Args":"-al"}'m --cert ../keys/client.signed.cert.pem --key ../keys/client.pem -d '{"Name":"ls","Args":"-al"}'
``

## [client](client) - TODO

Think some lingering basic auth still causing issues

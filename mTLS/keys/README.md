# keys

Need setup
* intermediate\_openssl.cnf
* openssl.cnf

## TODO

those 2 config files need a bit more love. Currently not using SANs and causing go client errors. However works in curl. Need to add sections for SubjectAlternativeName

# Running

```
./gen_cert.sh
```

Builds root cert, intermediate, client, server public and private keys

```
./sign_keys
```

Signs intermediate, client, server 


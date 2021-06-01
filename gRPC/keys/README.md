# keys

Modify - `dir=` line on both config files to current path
* intermediate\_openssl.cnf
* openssl.cnf

# Running

```
./gen_cert.sh
```

Builds root cert, intermediate, client, server public and private keys

```
./sign_keys
```

Signs intermediate, client, server 


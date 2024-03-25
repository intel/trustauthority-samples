# Encryptor
This is a utility program to encrypt the data using the wrapped DEK retrieved from KBS.

## Build

```shell
go build encrypt.go
```

Requires **Go 1.21 or newer**. See https://go.dev/doc/install for installation of Go.

## Usage

encrypt <data_file> <private_key_file> <wrapped_dek_file>

## Data Encryption Steps

### Generate RSA key-pair using openssl

```shell
openssl genrsa -out keypair.pem 2048
```
  
### Extract public key

```shell
openssl rsa -in keypair.pem -pubout -out public.crt
```

### Get DEK from KBS using public key

> [!Note]
> Get the auth token from KBS before requesting dek.

* **URL**
`https://{{kbs-ip}}:9443/kbs/v1/keys/<key_id>`

* **Method:**

`POST`

* **Headers:**

`"Accept" : "application/json"` <br>
`"Content-Type" : "application/x-pem-file"` <br>
`"Authorization": "Bearer <token>"`

* **Data Params:**

  `public-key=[string]`
  
* **Request Body:**
```
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0NTXPU3AojKVDVuKyRDk
QSm7AQF+2NXdOOgOaNQrvY/38gT8YP3HM0tpvRvUey/nBnGMj6MvGeadIzuYZOm8
H2F8nMuPbWy+zw6NgUdyrpvbduUCw3Lkgf+VqdpiQjGWMIQV5TzhUwZRNfz2VSom
sIbG6r4QivDpi9MOWvfreBXCzCRcyQ2y5gyxp4/Cm/WQwF6kKQhL1p/WWkdF9TiB
VxFAxP39G/D+lg/QKKX95rDGOJipn2a0ud0P+YnXbVsSU3BP3sdxHVUF/0Wha+/2
j0uNjOws7Pdxs1heyMB1D4nJOKdwRtS1RyC9fscznq4rlaJ6CYyyE07BCmzbJK2a
NQIDAQAB
-----END PUBLIC KEY-----
```

* **Success Response:**
  * **Code:** 200 <br>
    **Content:**
    ``` json
      {"wrapped_key":"HJ51f0EIAx1xKvXnC++NBCjAE/5RTrs2SSy4e0ZCGj4OuNACRqSkxNG5VSmiLzbp50ONCHUZC9/Opdu8xfx8k1yvzAFf+rTZUKWKGgc52td4oD85oPbWU3Dh9+8C+eCe/n0GyzM9FLRyWp+ykLJYDX51+6s/3V4wDwujdvMNGCcYR/2rrprmzZ/DAvNTej1P7Qz7lkIRnHM0znlk3XfVITpq2WqgUkz9PZOzOmgdqQ2drTvVQvCs3Dw8M7pi4LNEld4vdRD1JY599A13EOef0q+2/Op9XVX4qYUa7dlN7K/c0Fgj00dbwxNKHa2JI2B8TvJiA1su9+Yb1gWgNyaHHw=="}
    ```

* **Error Response:**
  * **Code:** 404 NOT FOUND <br />

Save the wrapped DEK from KBS in a file for running encryption tool later.

### Encrypt data
Execute **encrypt** binary with required args

```shell
./encrypt diabetes-linreg.model keypair.pem wrapped.key
```

## Security Considerations
1. This encryption tool needs to be run in a secure environment. In real world, encryption operation happens on enterprise side.
1. Make sure to remove the data file and private key post running encryptor.
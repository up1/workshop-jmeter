# Demo secure with Go app

## 1. Manage certificates in localhost

### 1.1 Create the Certificate Authority (CA)
* The CA key (ca.key) is used to sign both the server and client certificates, making them mutually trusted. The public part is saved as client_ca.crt
```
# Generates the private key for the CA (4096-bit RSA)
$openssl genrsa -out ca.key 4096

# Generates the self-signed CA certificate, valid for 10 years (3650 days). The -subj flag avoids interactive prompts
$openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out client_ca.crt -subj "/CN=My Root CA"
```

### 1.2 Create the Server Certificate
* The server needs its own private key and a certificate signed by the CA

```
# Generates the private key for the server (2048-bit RSA)
$openssl genrsa -out server.key 2048

# Creates a Certificate Signing Request (CSR). The Common Name (CN) must match the host the server runs on (e.g., localhost)
$openssl req -new -key server.key -out server.csr -subj "/CN=localhost"

# Signs the CSR using the CA's key/certificate, creating the final server certificate. This command also adds a Subject Alternative Name (SAN) extension for localhost, which is required by modern browsers/clients
$openssl x509 -req -in server.csr -CA client_ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256 -extfile <(printf "subjectAltName=DNS:localhost")

# Delete server.csr
$rm server.csr
```

### 1.3 Create a Client Certificate (For JMeter)
```
# Generates the private key for Client 1
$openssl genrsa -out client1.key 2048

# Creates a CSR for the client. The CN should be unique (e.g., JMeterUser1)
$openssl req -new -key client1.key -out client1.csr -subj "/CN=JMeterUser1"

# Signs the CSR using the CA's key/certificate, creating the final client certificate. The extendedKeyUsage=clientAuth flag explicitly marks this as a client certificate
$openssl x509 -req -in client1.csr -CA client_ca.crt -CAkey ca.key -CAcreateserial -out client1.crt -days 365 -sha256 -extfile <(printf "extendedKeyUsage=clientAuth")

# Delete client1.csr
$rm client1.csr
```

### 1.4 Convert PEM to PKCS12 (OpenSSL)
```
$openssl pkcs12 -export -out user1.p12 -inkey client1.key -in client1.crt -name "cert_user_1"
```

### 1.5 Consolidate PKCS12 to JKS (Keytool)
```
$keytool -importkeystore -srckeystore user1.p12 -srcstoretype PKCS12 -destkeystore jmeter_certs.jks -deststoretype JKS -srcalias cert_user_1 -destalias cert_user_1
```

## 2. Start app
```
$go run main.go
```

## 3. Testing
```
$curl https://localhost:8443
```
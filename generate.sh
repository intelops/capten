#.
#├── client
#├── intermediate
#│   ├── ca-cert.pem
#│   ├── ca-cert.srl
#│   ├── ca-key.pem
#│   ├── cert-chain.pem
#│   ├── cluster-ca.csr
#│   ├── intermediate.conf
#│   └── root-cert.pem
#├── root
#│   ├── root-ca.conf
#│   ├── root-cert.csr
#│   ├── root-cert.pem
#│   ├── root-cert.srl
#│   └── root-key.pem
#└── server
#    ├── server.crt
#    ├── server.csr
#    └── server.key




mkdir root intermediate server client

##Root certificate creation

cat <<EOF > root/root-ca.conf
[ req ]
encrypt_key = no
prompt = no
utf8 = yes
default_md = sha256
default_bits = 4096
req_extensions = req_ext
x509_extensions = req_ext
distinguished_name = req_dn
[ req_ext ]
subjectKeyIdentifier = hash
basicConstraints = critical, CA:true
keyUsage = critical, digitalSignature, nonRepudiation, keyEncipherment, keyCertSign
[ req_dn ]
O = Intelops
CN = Root CA
EOF

openssl genrsa -out root/root-key.pem 4096
openssl req -new -key root/root-key.pem -config root/root-ca.conf -out root/root-cert.csr
openssl x509 -req -days 1825 -signkey root/root-key.pem -extensions req_ext -extfile root/root-ca.conf -in root/root-cert.csr -out root/root-cert.pem

##Intermediate Certificate creation
cat <<EOF > intermediate/intermediate.conf
[ req ]
encrypt_key = no
prompt = no
utf8 = yes
default_md = sha256
default_bits = 4096
req_extensions = req_ext
x509_extensions = req_ext
distinguished_name = req_dn
[ req_ext ]
subjectKeyIdentifier = hash
basicConstraints = critical, CA:true, pathlen:0
keyUsage = critical, digitalSignature, nonRepudiation, keyEncipherment, keyCertSign
subjectAltName=@san
[ san ]
DNS.1 = dev.optimizor.app
[ req_dn ]
O = Intelops
CN = Optimizor CA
L = agent
EOF


openssl genrsa -out intermediate/ca-key.pem 4096
openssl req -new -config intermediate/intermediate.conf -key intermediate/ca-key.pem -out intermediate/cluster-ca.csr
openssl x509 -req -days 365 -CA root/root-cert.pem -CAkey root/root-key.pem -CAcreateserial -extensions req_ext -extfile intermediate/intermediate.conf -in intermediate/cluster-ca.csr -out intermediate/ca-cert.pem
cat intermediate/ca-cert.pem root/root-cert.pem > intermediate/cert-chain.pem
cp root/root-cert.pem intermediate/


##Server Certificate Creation
openssl req -newkey rsa:2048 -nodes -keyout server/server.key -subj "/O=Intelops Inc./CN=*.dev.optimizor.app" -out server/server.csr
openssl x509 -req -extfile <(printf "subjectAltName=DNS:captenagent.dev.optimizor.app,DNS:dev.optimizor.app") -days 365 -in server/server.csr -CA intermediate/ca-cert.pem -CAkey intermediate/ca-key.pem -CAcreateserial -out server/server.crt

##K8s secrets creation
#cd ..
#kubectl create secret generic cert-chain --from-file=ca.crt=intermediate/cert-chain.pem -n capten
#kubectl create secret tls capten-agent --cert=server/server.crt --key=server/server.key -n capten


##Client certificate Creation
openssl req -newkey rsa:2048 -nodes -keyout client/client.key -subj "/O=Intelops Inc./CN=*.dev.optimizor.app" -out client/client.csr
openssl x509 -req -extfile <(printf "subjectAltName=DNS:client.dev.optimizor.app") -days 365 -in client/client.csr -CA intermediate/ca-cert.pem -CAkey intermediate/ca-key.pem -CAcreateserial -out client/client.crt




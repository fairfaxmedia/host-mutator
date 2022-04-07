#! /bin/sh
set -eo pipefail

export APP="${1:-host-mutator}"
export NAMESPACE="${2:-default}"
export CSR_NAME="${APP}.${NAMESPACE}.svc"

rm -rf ssl && mkdir ssl && cd ssl

echo "Creating cert.key"
openssl genrsa -out cert.key 2048

echo "Creating csr.conf"
cat > csr.conf <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${APP}
DNS.2 = ${APP}.${NAMESPACE}
DNS.3 = ${CSR_NAME}
DNS.4 = ${CSR_NAME}.cluster.local
EOF
openssl req -new -key cert.key -subj "/CN=system:node:${CSR_NAME} /O=system:nodes" -out cert.csr -config csr.conf

echo "Deleting existing csr, if any"
kubectl delete csr "$CSR_NAME" || true

echo "Creating kubernetes CSR object"
echo "kubectl create -f -"
kubectl create -f - <<EOF
apiVersion: certificates.k8s.io/v1
kind: CertificateSigningRequest
metadata:
  name: ${CSR_NAME}
spec:
  groups:
  - system:authenticated
  request: $(base64 -i cert.csr | tr -d '\n')
  signerName: kubernetes.io/kubelet-serving
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF

SECONDS=0
while true; do
  echo "Waiting for csr to be present in kubernetes"
  if kubectl get csr "$CSR_NAME" > /dev/null 2>&1; then
      break
  fi
  if [[ $SECONDS -ge 60 ]]; then
    echo "Timed out waiting for csr"
    exit 1
  fi
  sleep 2
done

kubectl certificate approve "$CSR_NAME"

SECONDS=0
while true; do
  echo "Waiting for serverCert to be present in kubernetes"
  serverCert=$(kubectl get csr "$CSR_NAME" -o jsonpath='{.status.certificate}')
  if [[ $serverCert != "" ]]; then 
    break
  fi
  if [[ $SECONDS -ge 60 ]]; then
    echo "Timed out waiting for serverCert"
    exit 1
  fi
  sleep 2
done

echo "Creating cert.pem"
echo "$serverCert" | openssl base64 -d -A -out cert.pem
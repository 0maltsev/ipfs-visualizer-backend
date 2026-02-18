# ipfs-cluster-deployment

## Based on https://ipfscluster.io/documentation/guides/k8s/

- Run to generate Cluster Secret 
For Linux:
```
od  -vN 32 -An -tx1 /dev/urandom | tr -d ' \n' | base64 -w 0 - 
```
For Mac:
```
od -vN 32 -An -tx1 /dev/urandom | tr -d ' \n' | base64
```
- Run to generate Bootstrap Peer ID and Private Key (https://github.com/adityajoshi12/ipfs-key/releases/tag/v1.0.0)

- - Install IPFS-key utility

For Linux:
```
wget https://github.com/adityajoshi12/ipfs-key/releases/download/v1.0.0/ipfs-key_1.0.0_linux_amd64.tar.gz
tar -xvf ipfs-key_1.0.0_linux_amd64.tar.gz
sudo mv ipfs-key /usr/local/bin/
ipfs-key
```

For Mac:
```
wget https://github.com/adityajoshi12/ipfs-key/releases/download/v1.0.0/ipfs-key_1.0.0_darwin_arm64.tar.gz
tar -xvf ipfs-key_1.0.0_darwin_arm64.tar.gz
sudo mv ipfs-key /usr/local/bin/
ipfs-key
```

- - Generate Private Key
For Linux:
```
ipfs-key | base64 -w 0
```
For Mac: 
```
ipfs-key | base64
```

- - Encode base64 your Private Key
For Linux:
```
echo "<your_key>" | base64 -w 0
```
For Mac: 
```
echo "<your_key>" | base64
```

- Apply Kube resources
```
kubectl apply -f .
```

- Monitor creation of resources
```
kubectl get sts
kubectl get pod
```
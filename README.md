# Sample Extension Service for Kyverno

This repository contains a sample extension service for Kyverno that performs a namespace check to allow or denying creation of a resource. Requests where the namespace is `default` or empty are not allowed.

The [accompanying policy](/policies/check-namespace.yaml) checks namespaces for `ConfigMap` resources by making a call to the extenstion service. 

## Running

1. Install cert-manager

```sh
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
```

Refer to [cert-manager installation docs](https://cert-manager.io/docs/installation/)


2. Run the service

```sh
kubectl apply -f manifests/install.yaml 
```

3. Check access from a test pod

Run netshoot:

```sh
kubectl run -it netshoot --rm --image nicolaka/netshoot
```

Try accessing the HTTP GET endpoint with a valid namespace:

```sh
curl "http://sample.kyverno-extension/check-namespace?namespace=test"
```
Result:
> {"allowed": true}

Try accessing the HTTPS GET endpoint with an invalid:

```sh
curl "http://sample.kyverno-extension/check-namespace?namespace=default"
```

Result:
> {"allowed": false}

Try accessing the HTTPS POST endpoint with a valid namespace:

```sh
curl -k https://sample.kyverno-extension/check-namespace -X POST --data '{"namespace" : "test"}'
```
Result:
> {"allowed": true}

4. Test with a HTTP call from the policy

Install the policy:

```sh
kubectl apply -f policies/check-namespace.yaml
```

Try to create a `ConfigMap` in the default namespace. It will be blocked:

```sh
kubectl create cm test
```
Result:
> error: failed to create configmap: admission webhook "validate.kyverno.svc-fail" denied the request: 
>
> policy ConfigMap/default/test for resource violation: 
>
> check-namespaces:
> call-extension: namespace default is not allowed


Try to create a `ConfigMap` in some other namespace. It will be allowed:

```sh
kubectl create cm test -n kube-system --dry-run=server
```

Result:
> configmap/test created (server dry run)

5. Test with a HTTPS call from the policy

Get the generated certificate keychain:

```sh
kubectl get secret sample-tls -n kyverno-extension -o json | jq -r '.data."tls.crt"' | base64 -d && kubectl get secret sample-tls -n kyverno-extension -o json | jq -r '.data."ca.crt"' | base64 -d
```

Update the policy to configure your certificate chain:

## Building

Build service executable

```sh
make build
```

Build and push image

```sh
make ko
```

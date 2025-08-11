# Creating a Sample Kubernetes Controller

## Step 1: Define the API Package

Create the API definition files under the directory structure:

```
pkg/api/<groupname>/<version>/
```

Inside this folder, create the following files:

* `types.go` — Define your Custom Resource (CR) Go structs here.
* `doc.go` — Package documentation and API group metadata.
* `register.go` — Register your types with the scheme.

---

## Step 2: Generate Boilerplate Code

Use Kubernetes code generators to generate necessary deepcopy functions, clientsets, listers, and informers.

* Ensure you have the boilerplate files such as `boilerplate.go.txt`, `tools.go`, and the script `update-codegen.sh`.
* Run the code generator script:

```bash
bash hack/update-codegen.sh
```

This generates all the required client and deepcopy code based on your API definitions.

---

## Step 3: Generate the CustomResourceDefinition (CRD)

Use the `controller-gen` tool to generate CRD YAML manifests from your API types.

1. Install `controller-gen` if you haven't:

```bash
go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
```

2. Generate CRDs and output them to the manifests directory:

```bash
go run sigs.k8s.io/controller-tools/cmd/controller-gen crd paths=./pkg/api/... output:crd:dir=./manifests
```

---

## Step 4: Implement the Controller Logic

* Implement the reconciliation logic for your controller.
* Define event handlers and business logic to react to changes in your custom resources.

---

## Step 5: Register the Controller with the Manager

In your `main.go`:

* Configure the Kubernetes client and clientset.
* Create shared informer factories.
* Initialize your controller.
* Register the controller with the controller-runtime manager.

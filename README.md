# Pathfinder

## CRD Development


First of all install kubebuilder 

```bash
os=$(go env GOOS)
arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
curl -L https://go.kubebuilder.io/dl/2.3.1/${os}/${arch} | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
sudo mv /tmp/kubebuilder_2.3.1_${os}_${arch} /usr/local/kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin
```

Create a directory and do 

```bash
kubebuilder init --domain xmbsmdsj
# This will generate an empty project
# Including entry, dependencies, initial kubernetes configs
```

Add an api Pathfinder
```bash
kubebuilder create api --group pathfinder --version v1 PathFinder

# We will need both type and controller for this, so select Yes for both
# This step setup data type and some utils code like deepcopy and Scheme registration for 
# Out CRD
```


Customize spec.

Spec holds the desired state of a resource, in this case, we do not desire anything, but maintain a list of services.
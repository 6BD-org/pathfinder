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


## How to use? 

### Step 1: Deploy CRD to your cluster

```bash
$ make install
```

### Step 2: Build image for pathfinder controller


```bash
$ make docker-push
```

### Step 3: Deploy third party components and pathfinder controller on cluster

```bash
# cert-manager is required for pathfinder to run
$ make deploy-thirdparty
$ make deploy
```

### Step 4: Deploy Sample

```bash

make deploy-sample

```

### Step 5: Setup region for your pathfinder

Just modify the `region` value in `Pathfinder`'s `spec`

### Step N: Make your service discoverable

In order to make service discoverable, you only need to add a few annotations to service object, like this

```yaml

apiVersion: v1
kind: Service
metadata:
  annotations:
    XM-PathFinder-Region: some-region # Region must match existing path-finder's region
    XM-PathFinder-Service: Activated
    XM-PathFinder-ServiceName: my-svc
```
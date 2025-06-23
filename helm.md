# Using hep-sidekick with Helm

`hep-sidekick` is designed to work seamlessly with applications deployed via Helm. By deploying `hep-sidekick` into your cluster, you can dynamically monitor any application, such as FreeSWITCH, Kamailio, or any other Session Border Controller (SBC), without modifying their original Helm charts.

The integration process involves two main steps:
1.  **Labeling Target Pods**: Ensuring the pods you want to monitor have a specific label that `hep-sidekick` can discover.
2.  **Deploying `hep-sidekick`**: Running `hep-sidekick` as a deployment within your cluster, configured to look for the target labels.

## Example: Attaching to a FreeSWITCH Deployment

Let's say you have a FreeSWITCH instance deployed in your cluster using a Helm chart. Here's how you can use `hep-sidekick` to monitor it.

### 1. Labeling the FreeSWITCH Pods

The easiest way to label the pods is by modifying the `values.yaml` file of your FreeSWITCH Helm chart. Most charts provide a way to add custom labels to pods.

For example, your `values.yaml` might look like this:

```yaml
# values.yaml for a hypothetical FreeSWITCH chart
# ... other values ...

# Add a label for hep-sidekick to discover these pods
podLabels:
  hep-sidekick/enabled: "true"

# ... other values ...
```

After adding the label, upgrade your Helm release:
```bash
helm upgrade freeswitch . -f values.yaml
```

If your Helm chart does not support adding pod labels directly, you can use `kubectl patch` to add the labels after the deployment is created:
```bash
kubectl patch deployment <your-freeswitch-deployment> -p '{"spec":{"template":{"metadata":{"labels":{"hep-sidekick/enabled":"true"}}}}}'
```

### 2. Deploying hep-sidekick

Next, you need to deploy `hep-sidekick` to your cluster. This involves creating a `ServiceAccount`, a `ClusterRole` with the necessary permissions, a `ClusterRoleBinding`, and the `Deployment` itself.

Below is a sample manifest (`hep-sidekick-deployment.yaml`) you can use.

**Note:** Remember to change the `--homer-address` to point to your actual HOMER server.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: hep-sidekick-sa
  namespace: default # Or a dedicated namespace for observability tools

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: hep-sidekick-clusterrole
rules:
- apiGroups: [""]
  # Needed to discover target pods and create heplify pods
  resources: ["pods"]
  verbs: ["get", "list", "watch", "create"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: hep-sidekick-clusterrolebinding
subjects:
- kind: ServiceAccount
  name: hep-sidekick-sa
  namespace: default
roleRef:
  kind: ClusterRole
  name: hep-sidekick-clusterrole
  apiGroup: rbac.authorization.k8s.io

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hep-sidekick-deployment
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hep-sidekick
  template:
    metadata:
      labels:
        app: hep-sidekick
    spec:
      serviceAccountName: hep-sidekick-sa
      containers:
      - name: hep-sidekick
        image: ghcr.io/sipcapture/hep-sidekick:latest
        imagePullPolicy: Always
        args:
        # This selector should match the label you added to your FreeSWITCH pods
        - "--selector=hep-sidekick/enabled=true"
        # The address of your HOMER server
        - "--homer-address=YOUR_HOMER_IP:9060"
        # Optional: Custom arguments for the heplify agent
        - "--heplify-args=-i any -m SIP -l 7"
```

Save this manifest to a file (e.g., `hep-sidekick-deployment.yaml`) and apply it to your cluster:

```bash
kubectl apply -f hep-sidekick-deployment.yaml
```

Once applied, `hep-sidekick` will start running, find your FreeSWITCH pods, and launch `heplify` agents to begin capturing traffic. 
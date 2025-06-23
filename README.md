<img src="https://user-images.githubusercontent.com/1423657/55069501-8348c400-5084-11e9-9931-fefe0f9874a7.png" width=200/>

# hep-sidekick

`hep-sidekick` is a Kubernetes tool for the High-Performance Stored Procedure (HEP) stack. It allows you to dynamically attach a `heplify` sniffing agent to existing pods in your cluster, providing real-time packet capture for your HOMER/SIPCAPTURE observability platform.

## Getting Started

### Prerequisites

- Go (version 1.18 or higher recommended)
- Access to a Kubernetes cluster (a valid `kubeconfig` file)

### Building

To build the tool from the source code, run the following command:

```bash
make build
```

This will produce the `hep-sidekick` binary in the root of the project directory.

## Usage

`hep-sidekick` discovers pods to monitor using a Kubernetes label selector. By default, it looks for pods with the label `hep-sidekick/enabled=true`. You can override this with the `--selector` flag.

You must also specify the address of your HOMER server where `heplify` will send the captured HEP packets using the `--homer-address` flag.

```bash
./hep-sidekick [flags]
```

### Flags

- `--selector`: Label selector to find pods to attach to. (default: `hep-sidekick/enabled=true`)
- `--homer-address`: Address of the HOMER server. (default: `127.0.0.1:9060`)
- `--heplify-args`: A string of arguments to pass to the heplify container. (default: `"-i any -t pcap"`)

### Example

To attach to all pods with the label `app=my-voip-app` and send HEP packets to a HOMER instance at `10.0.0.1:9060`, you would run:

```bash
./hep-sidekick --selector="app=my-voip-app" --homer-address="10.0.0.1:9060"
```

### Customizing heplify Settings

The `heplify` agent has many configuration options that can be controlled with command-line arguments. You can pass a custom set of arguments to the `heplify` containers using the `--heplify-args` flag.

The `--homer-address` is handled separately, so you do not need to include the `-hs` flag in your custom arguments string.

For example, to change the sniffing interface to `eth0` and set the capture mode to `SIP`, you could run:
```bash
./hep-sidekick \
  --selector="app=my-voip-app" \
  --homer-address="10.0.0.1:9060" \
  --heplify-args="-i eth0 -m SIP"
```

## How it Works

`hep-sidekick` does not inject any containers into your existing pods. Instead, it uses a less intrusive method inspired by tools like `kubeshark`:

1.  It uses the Kubernetes API to find pods that match the provided label selector.
2.  For each matching target pod, it schedules a new, temporary `heplify` pod on the **same Kubernetes node**.
3.  This new `heplify` pod is configured to use the node's host network (`hostNetwork: true`) and is given privileged access. This allows `heplify` to sniff all traffic on the node's network interfaces, including the traffic for the target pod, without altering the target pod's specification.

## Future Work

- **Lifecycle Management**: Watch for target pod termination to automatically clean up associated `heplify` pods.
- **Containerization**: Package `hep-sidekick` into a container for in-cluster deployment.
- **Enhanced Error Handling**: Improve resilience with more robust error handling and retry logic.

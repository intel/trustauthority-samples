# Build and Deployment of Intel® Key Broker Service and Intel TDX Demo App

These instructions describe how to build and deploy the two Docker containers comprising the Intel® Trust Domain Extensions (Intel TDX) demonstration application. For more information about the secure key release workflow and the demo app, see [Key release and workload for Intel TDX](https://docs.trustauthority.intel.com/articles/tutorial-tdx-workload.html).

There are four components required to build the demo app:

- Dockerfile to build the Intel® Key Broker Service (Intel KBS) relying party container.
- **kbs.env** environment file to configure Intel KBS. You must modify this file to suit your environment.
- Dockerfile to build the Intel TDX workload container. The Intel TDX workload build process supports Intel SGX DCAP, Azure confidential VMs with Intel TDX, and Google Cloud Platform (GCP) Confidential VMs.
- **workload.env** environment file to supply the workload with Intel KBS and Intel Trust Authority connection info.

## Prerequisites
- [Docker engine >= 20.10](https://docs.docker.com/engine/install/)

## Intel KBS

This section describes how to build and deploy the Intel KBS relying party container. For this demo, Intel KBS uses [PyKMIP](https://github.com/OpenKMIP/PyKMIP) for the key management system (KMS). PyKMIP is a relatively simple, standalone KMS that is useful for development and testing, but it is not a secure, production-quality KMS. Don't use PyKMIP for storing actual secrets.

### Build
1. Download the [Dockerfile](./kbs/Dockerfile) for building the container.
   ```bash
   wget https://raw.githubusercontent.com/intel/trustauthority-samples/main/deployment/kbs/Dockerfile
   ```
2. Run the following command to build the Docker image for Intel KBS.
   ```bash
   sudo docker build -t relying-party .
   ```

### Deploy
1. Download the [kbs.env](./kbs/kbs.env) file and update the missing configurations.
   ```bash
   wget https://raw.githubusercontent.com/intel/trustauthority-samples/main/deployment/kbs/kbs.env
   ```

   When the relying-party container is run, a startup script configures Intel KBS to use the `ADMIN_USERNAME` and `ADMIN_PASSWORD` provided in the **kbs.env** file. You'll also need to provide an attestation API key for `TRUSTAUTHORITY_API_KEY`, which you can get from the Intel Trust Authority [portal](https://portal.trustauthority.intel.com). You might also need to configure an HTTPS proxy, depending on your network configuration. In the following sample kbs.env file, you should replace the parameters marked with angle brackets, for example `<admin-username>`.

   ```
   LOG_LEVEL=INFO
   KEY_MANAGER=kmip
   ADMIN_USERNAME=<admin-username>
   ADMIN_PASSWORD=<admin-password>
   KMIP_VERSION=2.0
   KMIP_SERVER_IP=127.0.0.1
   KMIP_SERVER_PORT=5696
   KMIP_CLIENT_KEY_PATH=/etc/pykmip/client_key.pem
   KMIP_CLIENT_CERT_PATH=/etc/pykmip/client_certificate.pem
   KMIP_ROOT_CERT_PATH=/etc/pykmip/ca_certificate.pem
   PYKMIP_LOG_LEVEL=info
   HTTPS_PROXY=<proxy url if required>
   TRUSTAUTHORITY_BASE_URL=https://portal.trustauthority.intel.com
   TRUSTAUTHORITY_API_URL=https://api.trustauthority.intel.com
   TRUSTAUTHORITY_API_KEY=<API Key>
   ```
   > [!NOTE]
   > The kbs.env file contains secrets, such as the Intel KBS administrator password, the KMS password or access token, and an Intel Trust Authority attestation API key. For a secure production system, you should delete the kbs.env file after initial configuration to avoid compromising the configuration secrets. However, for this demo it's not necessary to delete the kbs.env, unless you want to protect the API key.

2. Run the following command to start Intel KBS as a Docker container.
   ```bash
   sudo docker run --name relying-party -d --env-file kbs.env --restart=always -p 9443:9443 relying-party:latest
   ```

   9443 is the port Intel KBS will listen on. The port can be changed, if necessary.

## Intel TDX Demo workload

This section describes how to build and deploy the demonstration workload. The workload is an application that runs a machine learning (ML) model in an Intel TDX trust domain. The ML model is assumed to contain confidential data and parameters, therefore it is encrypted when first deployed. The demo workload (attester) must obtain an Intel Trust Authority attestation token and include the token with a request to Intel KBS to retrieve the key needed to decrypt the ML model. This workflow implements a secure key release use case in passport attestation mode.

### Build
1. Download the [Dockerfile](./sample-workload/Dockerfile) for building the Intel TDX demo workload.
   ```bash
   wget https://raw.githubusercontent.com/intel/trustauthority-samples/main/deployment/sample-workload/Dockerfile
   ```
2. Run the following command to build the Docker image for the TDX demo workload.
   ```bash
   sudo docker build --build-arg PLATFORM=dcap -t trustauthority-demo .
   ```
   The demo workload relies on the Intel Trust Authority CLI for Intel TDX to obtain a quote and request an attestation token. The CLI requires an adapter that is specific to the platform on which the demo app is run. The 'PLATFORM' argument specifies which adapter to use:

   | PLATFORM | Description |
   |---|---|
   | `dcap` | Installs the adapter for native Intel TDX TDs that use Intel SGX DCAP to collect a quote for attestation.|
   | `azure` | Installs the adapter for Microsoft Azure confidential VMs with Intel TDX. |
   | `gcp` | Installs the adapter for Google Cloud Platform confidential VMs. |

### Deploy
1. Download the [workload.env](./sample-workload/workload.env) file and update the missing configurations.
   ```bash
   wget https://raw.githubusercontent.com/intel/trustauthority-samples/main/deployment/sample-workload/workload.env
   ```
   The **workload.env** file contains the settings needed for the sample workload to communicate with Intel KBS and Intel Trust Authority. This file contains secrets. Embedding secrets in an unencrypted file is _not_ a recommended best practice for secure systems. It's done here to make the demo easy to configure.

   The `KBS_ADMIN` and `KBS_PASSWORD` must match the `ADMIN_USERNAME` and `ADMIN_PASSWORD` used in KBS.env.

   The `KBS_URL` is of the form `https://<IP_address>:9443/kbs/v1`, where IP_address is the IP address of the host machine where the relying-party (KBS) container is running. If the KBS is running on same TDVM as the workload, then provide `172.17.0.1` as IP_address. Port 9443 is the default listening port for Intel KBS. If you change the port number for the Intel KBS container, you'll need to make a corresponding change to `KBS_URL`.

   `SKIP_TLS_VERIFICATION` if set to **true** will skip the TLS server certificate verification. For this demo, set it to **true** since we are using self-signed TLS certificates for KBS.

   `TRUSTAUTHORITY_API_KEY` is the same as used in **kbs.env**.

   ```
   KBS_ADMIN=<KBS-admin-username>
   KBS_PASSWORD=<KBS-admin-password>
   KBS_URL=https://<IP_address>:9443/kbs/v1
   SKIP_TLS_VERIFICATION=<true | false>
   HTTPS_PROXY=<proxy URL if required>
   TRUSTAUTHORITY_API_URL=https://api.trustauthority.intel.com
   TRUSTAUTHORITY_API_KEY=<API Key>
   ```
2. Run the following command to start the demo workload as a Docker container.
   ```bash
   sudo docker run --name ita-demo -d --env-file workload.env --device=/dev/tdx_guest -p 12780:12780 --user 0 trustauthority-demo:latest
   ```
   On successful run, the container will generate a `execute_workload_flow.env` file in the `/tmp/` folder. This env file will be used later for executing the workload flow script. Update the env variables in the `execute_workload_flow.env` file if you used custom settings.

> [!NOTE]
> If running on Azure, replace `--device=/dev/tdx_guest` with `--device=/dev/tpmrm0`  
> The command above will start the container as root user. This is needed because the Intel TDX demo workload container requires elevated privileges to access `/dev/tdx_guest` or `/dev/tpmrm0` device for quote collection.  
> If you don't want to run the container as root user, make sure the `/dev/tdx_guest` or `/dev/tpmrm0` device is accessible to the container. One of the ways to do this is to pass `gid` owning the `tdx_guest` or `tpmrm0` device to container.
> ```bash
> sudo docker run --name ita-demo -d --env-file workload.env --device=/dev/tdx_guest -p 12780:12780 --group-add $(getent group <user-group> | cut -d: -f3) trustauthority-demo:latest
> ```
> Where `<user-group>` is the group owning the `tdx_guest` or `tpmrm0` device
> ```bash
> crw-rw---- 1 root <user-group> 10, 123 /dev/tdx_guest
> ```

### Execute Workload Flow
1. Download the [execute_workload_flow.sh](./sample-workload/execute_workload_flow.sh) script for executing secure key release workflow.
   ```bash
   wget https://raw.githubusercontent.com/intel/trustauthority-samples/main/deployment/sample-workload/execute_workload_flow.sh
   ```
2. Run the following command to execute the workload flow script.
   ```bash
   bash execute_workload_flow.sh
   ```

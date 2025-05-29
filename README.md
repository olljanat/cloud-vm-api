# Cloud VM API
Unified API for manage virtual machines across multiple cloud providers.


Architecture
* The API is stateless and handles all backend calls synchronously
* Authentication is handled via Bearer token containing base64-encoded credentials which are passed directly to backend
* The server supports AWS, Azure, ESXi, Nutanix, and Proxmox
* Utilizing [CloudMuX](https://github.com/yunionio/cloudmux/) library from [Cloudpods](https://www.cloudpods.org) project to handle most of the cloud specific logic.
  * We have also [our own](https://github.com/olljanat/cloudmux) stripped version of CloudMuX code to simplify studying their code.

## Setup
1. Ensure that environments.json and machinetypes.json are updated match to your environment is configured with your cloud environments
2. Run the server:
```bash
go run main.go 
```
See table below how to fill those:
| Cloud   | Access Key               | Secret                     | VPC ID                          | Project             | Network ID | URL                |
| ------- | ------------------------ | -------------------------- | ------------------------------- | ------------------- | ---------- | ------------------ |
| Aws     | Access key               | Secret access key          | `VPC ID`                        | N/A                 | Subnet ID  | N/A                |
| Azure   | Service principal appId  | Service principal password | `<Tenant ID>/<Subscription ID>` | Resource group name | Subnet ID  | `AzurePublicCloud` |
| ESXi    | Service account username | Service account password   | ||||
| Nutanix | Service account username | Service account password   | ||||
| Proxmox | Service account username | Service account password   | `node/<node name>`              | N/A                 | N/A        | Proxmox URL        |

## Usage
Bearer token for all platforms are created with command like:
```bash
echo '<access key>:<secret >' | base64 -w 0
```

### List VMs
```bash
curl -X GET http://localhost:8080/vm?environment=aws-test1 \
  -H "Authorization: Bearer <Bearer token>"
```

### Create VM
Create a VM by sending a POST request to /vm:
```bash
cloud_init=$(cat <<EOF | base64
#cloud-config
users:
- name: dev
  sudo: ALL=(ALL) NOPASSWD:ALL
EOF
)
echo $cloud_init

curl -X POST http://localhost:8080/vm \
  -H "Authorization: Bearer <Bearer token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test1",
    "environment": "<env name>",
    "machinetype": "medium-debian",
    "cloud_init": "<base64 encoded cloud init>"
  }'
```

### Stop VM
```bash
curl -X GET http://localhost:8080/vm/i-048501077e53646ed/stop?environment=aws-test1 \
  -H "Authorization: Bearer <Bearer token>"
```

### Start VM
```bash
curl -X GET http://localhost:8080/vm/i-048501077e53646ed/start?environment=aws-test1 \
  -H "Authorization: Bearer <Bearer token>"
```

### Delete VM
```bash
curl -X DELETE http://localhost:8080/vm/i-048501077e53646ed?environment=aws-test1 \
  -H "Authorization: Bearer <Bearer token>"
```

# Troubleshooting
You can see API calls send to backend by running [mitmproxy](https://mitmproxy.org)
and running API server like this:
```bash
export HTTP_PROXY="http://192.168.8.226:8080"
export HTTPS_PROXY="http://192.168.8.226:8080"
export NO_PROXY="login.microsoftonline.com"
go run main.go
```

> [!NOTE]
> You need add mitmproxy [certificate](https://docs.mitmproxy.org/stable/concepts/certificates/) to trusted authorities list.
Other why API calls will fail.

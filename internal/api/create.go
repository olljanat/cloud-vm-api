package api

import (
	"encoding/json"
	"net/http"

	"github.com/olljanat/cloud-vm-api/internal/config"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

type VMCreateRequest struct {
	Name        string `json:"name"`
	Environment string `json:"environment"`
	MachineType string `json:"machinetype"`
	CloudInit   string `json:"cloud_init,omitempty"`
}

type VMCreateResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func CreateVMHandler(w http.ResponseWriter, r *http.Request) {
	var req VMCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	host, env, status, err := getHost(r, req.Environment)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	mtCloud, mtConfig, err := config.GetMachineTypeConfig(req.MachineType, env.Cloud)
	if err != nil {
		http.Error(w, "Invalid machinetype for this cloud: "+err.Error(), http.StatusBadRequest)
		return
	}

	vmConfig := &cloudprovider.SManagedVMCreateConfig{
		Name:              req.Name,
		Hostname:          req.Name,
		ProjectId:         env.ProjectId,
		ExternalImageId:   mtCloud.Image,
		InstanceType:      mtCloud.InstanceType,
		ExternalNetworkId: env.NetworkId,
		ExternalVpcId:     env.VpcId,
		UserData:          req.CloudInit,
		SysDisk: cloudprovider.SDiskInfo{
			StorageExternalId: env.StorageExternalId,
			SizeGB:            mtConfig.OsDiskSizeGB,
		},
		OsType:         mtConfig.OsType,
		OsDistribution: mtConfig.OsDistribution,
	}

	if env.Cloud == "Azure" {
		vmConfig.NameEn = env.Name
	}

	if env.Cloud == "Proxmox" {
		cpu, memory := getProxmoxSpec(mtCloud.InstanceType)
		vmConfig.Cpu = cpu
		vmConfig.MemoryMB = memory
	}

	vm, err := host.CreateVM(vmConfig)
	if err != nil {
		http.Error(w, "Create VM failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := VMCreateResponse{
		ID:     vm.GetGlobalId(),
		Name:   vm.GetName(),
		Status: vm.GetStatus(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

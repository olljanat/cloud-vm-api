package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/olljanat/cloud-vm-api/internal/auth"
	"github.com/olljanat/cloud-vm-api/internal/cloud"
	"github.com/olljanat/cloud-vm-api/internal/config"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func getHost(r *http.Request, envName string) (cloudprovider.ICloudHost, *config.Environment, int, error) {
	if envName == "" {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("missing environment parameter")
	}
	env, err := config.GetEnvironment(envName)
	if err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("environment with name %s not found", envName)
	}

	token := r.Header.Get("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	} else {
		return nil, nil, http.StatusUnauthorized, fmt.Errorf("missing or invalid Authorization header")
	}
	creds, err := auth.DecodeCredentials(token)
	if err != nil {
		return nil, nil, http.StatusUnauthorized, fmt.Errorf("unable parsing credentials. Error: %s", err)
	}

	provider, err := cloud.NewCloudProvider(env, creds)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to initialize cloud provider: %s", err)
	}
	regionId := env.Region
	region, err := provider.GetIRegionById(regionId)
	if err != nil {
		regions, _ := provider.GetIRegions()
		regionIds := []string{}
		for _, region := range regions {
			regionIds = append(regionIds, region.GetId())
		}
		return nil, nil, http.StatusInternalServerError, fmt.Errorf("failed to get region: %s, Error: %s , Available regions: %s", regionId, err, regionIds)
	}

	var host cloudprovider.ICloudHost
	if env.Cloud == "Proxmox" {
		host, err = region.GetIHostById(env.VpcId)
		if err != nil {
			return nil, nil, http.StatusInternalServerError, fmt.Errorf("no host with name pve found: %s", err)
		}
	} else {
		hosts, err := region.GetIHosts()
		if err != nil || len(hosts) == 0 {
			return nil, nil, http.StatusInternalServerError, fmt.Errorf("no hosts found in region: %s", env.Region)
		}

		// FixMe: We might want to select other host than just first one?
		host = hosts[0]
	}

	return host, env, http.StatusOK, nil
}

func getVMByID(host cloudprovider.ICloudHost, env *config.Environment, vmID string) (cloudprovider.ICloudVM, error) {
	if env.Cloud == "Azure" {
		parts := strings.Split(env.VpcId, "/")
		vmID = fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/microsoft.compute/virtualmachines/%s", parts[1], env.ProjectId, vmID)
	}

	vms, err := host.GetIVMs()
	if err != nil {
		return nil, err
	}
	for _, vm := range vms {
		if vm.GetGlobalId() == vmID {
			return vm, nil
		}
	}
	return nil, fmt.Errorf("VM not found")
}

func getProxmoxSpec(t string) (int, int) {
	cpu := 0
	ramMB := 0

	re := regexp.MustCompile(`c(\d+)m(\d+)`)
	matches := re.FindStringSubmatch(t)
	if len(matches) == 3 {
		cpuStr := matches[1]
		ramStr := matches[2]

		cpu, err1 := strconv.Atoi(cpuStr)
		ram, err2 := strconv.Atoi(ramStr)
		ramMB = ram * 1024

		if err1 == nil && err2 == nil {
			fmt.Printf("CPU: %d\n", cpu)
			fmt.Printf("RAM: %d GB\n", ram)
		} else {
			fmt.Println("Error converting to integer")
		}
	} else {
		fmt.Println("Pattern not found in input string")
	}
	return cpu, ramMB
}

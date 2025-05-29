package cloud

import (
	"net/http"
	"net/url"
	"os"

	"github.com/olljanat/cloud-vm-api/internal/auth"
	"github.com/olljanat/cloud-vm-api/internal/config"
	"golang.org/x/net/http/httpproxy"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	_ "yunion.io/x/cloudmux/pkg/multicloud/aws/provider"
	_ "yunion.io/x/cloudmux/pkg/multicloud/azure/provider"
	_ "yunion.io/x/cloudmux/pkg/multicloud/esxi/provider"
	_ "yunion.io/x/cloudmux/pkg/multicloud/nutanix/provider"
	_ "yunion.io/x/cloudmux/pkg/multicloud/proxmox/provider"
)

func NewCloudProvider(env *config.Environment, creds *auth.Credentials) (cloudprovider.ICloudProvider, error) {

	proxyCfg := &httpproxy.Config{
		HTTPProxy:  os.Getenv("HTTP_PROXY"),
		HTTPSProxy: os.Getenv("HTTPS_PROXY"),
		NoProxy:    os.Getenv("NO_PROXY"),
	}
	cfgProxyFunc := proxyCfg.ProxyFunc()
	proxyFunc := func(req *http.Request) (*url.URL, error) {
		return cfgProxyFunc(req.URL)
	}

	// cfg.Vendor must be exactly "AWS", "Azure", etc. (case-sensitive)

	cfg := cloudprovider.ProviderConfig{
		Id:        "1",
		Name:      env.Name,
		Vendor:    env.Cloud,
		URL:       env.Url,
		ProxyFunc: proxyFunc,
	}

	if env.Cloud == "Azure" {
		cfg.Account = env.VpcId
		cfg.Secret = creds.AccessKey + "/" + creds.Secret
	} else {
		cfg.Account = creds.AccessKey
		cfg.Secret = creds.Secret
	}

	return cloudprovider.GetProvider(cfg)
}

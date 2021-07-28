package cmd

import (
	"context"
	"fmt"
	"github.com/up9inc/mizu/cli/kubernetes"
	"github.com/up9inc/mizu/cli/mizu"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
)

func runMizuView(mizuViewOptions *MizuViewOptions) {
	var resourcesNamespace string
	if mizuViewOptions.MizuNamespace != "" {
		resourcesNamespace = mizuViewOptions.MizuNamespace
	} else {
		resourcesNamespace = mizu.ResourcesDefaultNamespace
	}

	kubernetesProvider, err := kubernetes.NewProvider(mizuViewOptions.KubeConfigPath)
	if err != nil {
		if clientcmd.IsEmptyConfig(err) {
			fmt.Printf(mizu.Red, "Couldn't find the kube config file, or file is empty. Try adding '--kube-config=<path to kube config file>'\n")
			return
		}
		if clientcmd.IsConfigurationInvalid(err) {
			fmt.Printf(mizu.Red, "Invalid kube config file. Try using a different config with '--kube-config=<path to kube config file>'\n")
			return
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exists, err := kubernetesProvider.DoesServicesExist(ctx, resourcesNamespace, mizu.ApiServerPodName)
	if err != nil {
		panic(err)
	}
	if !exists {
		fmt.Printf("The %s service not found\n", mizu.ApiServerPodName)
		return
	}

	mizuProxiedUrl := kubernetes.GetMizuApiServerProxiedHostAndPath(mizuViewOptions.GuiPort)
	_, err = http.Get(fmt.Sprintf("http://%s/", mizuProxiedUrl))
	if err == nil {
		fmt.Printf("Found a running service %s and open port %d\n", mizu.ApiServerPodName, mizuViewOptions.GuiPort)
		return
	}
	fmt.Printf("Found service %s, creating k8s proxy\n", mizu.ApiServerPodName)

	fmt.Printf("Mizu is available at  http://%s\n", kubernetes.GetMizuApiServerProxiedHostAndPath(mizuViewOptions.GuiPort))
	err = kubernetes.StartProxy(kubernetesProvider, mizuViewOptions.GuiPort, resourcesNamespace, mizu.ApiServerPodName)
	if err != nil {
		fmt.Printf("Error occured while running k8s proxy %v\n", err)
	}
}

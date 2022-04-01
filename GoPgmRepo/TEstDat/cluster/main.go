package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	esxiAddr      = "https://%s:%s@%s/sdk"
	hostSystemObj = "HostSystem"
	summaryProp   = "summary"
	runtimeProp   = "runtime"
	configProp    = "config"
	hardwareProp  = "hardware"
	datastoreProp = "datastore"
	apiProtocol   = "https://"
	dataStoreObj  = "Datastore"
	infoProp      = "info"
)

var roleIdNameMapping map[int32]string

// Host represents connection required information
type Host struct {
	RegID     string `json:"regID,omitempty"`
	Address   string `json:"address"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	UUID      string
	GUID      string
	TargetIPs []string `json:"targetIPs"` // DEPRECATED
	VCenterIP string   `json:"vcenterIP"`
}

func getAPIPath(host, username, password string) (*url.URL, error) {
	userInfo := url.UserPassword(username, password)
	raw := apiProtocol + userInfo.String() + "@" + host + "/sdk"
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse url of IP %v ,password and username not shown because of security reasons", host)
	}

	return u, nil
}

// NewClient creates session with vmware server
func NewClient(ctx context.Context, h Host) (*govmomi.Client, error) {
	url, err := getAPIPath(h.Address, h.Username, h.Password)
	if err != nil {
		return nil, err
	}

	return govmomi.NewClient(ctx, url, true)
}

// CancelHandler returns cancellation context and function for graceful shutdown
func CancelHandler() (context.Context, context.CancelFunc) {
	ctx, cancelFn := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		cancelFn()
		signal.Stop(signals)
	}()

	return ctx, cancelFn
}

// AuthorizationManager retrieves AuthorizationManager information
func AuthorizationManager(ctx context.Context, c *govmomi.Client) error {
	roleIdNameMapping = make(map[int32]string)
	m := object.NewAuthorizationManager(c.Client)

	list, err := m.RoleList(ctx)
	if err != nil {
		return err
	}

	for _, val := range list {
		roleIdNameMapping[val.RoleId] = val.Name
	}

	return nil
}

func ClusterComputeResource(ctx context.Context, c *govmomi.Client, host []mo.HostSystem) {
	pc := property.DefaultCollector(c.Client)

	for _, objHost := range host {
		fmt.Println("for host", objHost.Summary.Config.Name)
		var hs mo.HostSystem
		pc.RetrieveOne(ctx, objHost.Self, []string{"parent"}, &hs)
		var css mo.ClusterComputeResource
		err := pc.RetrieveOne(ctx, *hs.Parent, []string{"parent", "name", "permission", "summary", "configuration", "configurationEx", "host", "datastore", "overallStatus"}, &css)
		if err != nil {
			fmt.Println("Error in retiriving cluster   ", err)
		}
		var folder mo.Folder
		err = pc.RetrieveOne(ctx, *css.Parent, []string{"parent"}, &folder)
		if err != nil {
			fmt.Println("Error in retiriving parent from cluster")
		}
		if folder.Parent.Type == "Datacenter" {
			fmt.Println("type of parent", folder.Parent.Type)
		} else {
			folder = pritnDatacenterName(ctx, pc, folder)
		}
		var datacenter mo.Datacenter
		err = pc.RetrieveOne(ctx, *folder.Parent, []string{"name"}, &datacenter)
		if err != nil {
			fmt.Println("Error in getting datacenter", err)
		}
		var objdatacenter mo.Datacenter
		err = pc.RetrieveOne(ctx, datacenter.Self, []string{"name"}, &objdatacenter)
		if err == nil {
			fmt.Println("\n********************************")
			fmt.Println("Datacenter name", objdatacenter.Name)
		}
		var hostSystems []mo.HostSystem
		err = pc.Retrieve(ctx, css.Host, []string{"summary", "runtime"}, &hostSystems)
		if err != nil {
			fmt.Println("Error in getting hostsystem", err)
		}
		for _, obj := range hostSystems {
			fmt.Println("name of host", obj.Summary.Config.Name)
			fmt.Println("connection state of the host", obj.Runtime.ConnectionState)
		}
		var datastore []mo.Datastore
		err = pc.Retrieve(ctx, css.Datastore, []string{"summary", "info"}, &datastore)
		if err != nil {
			fmt.Println("Error in datastore", err)
		}
		for _, val := range datastore {
			fmt.Println("\nName\t", val.Summary.Name)
			fmt.Println("Type\t", val.Summary.Type)
			fmt.Println("Free Space\t", val.Summary.FreeSpace)
			fmt.Println("Capacity\t", val.Summary.Capacity)
		}
		fmt.Printf("Cluster Name :: %+v\n", css.Name)
		if css.Configuration.DasConfig.Enabled != nil {
			fmt.Printf("HA Enabled :: %t\n", *css.Configuration.DasConfig.Enabled)
		}
		fmt.Println("Cluster status", css.OverallStatus)
		fmt.Println("Host Monitoring: ..", css.Configuration.DasConfig.HostMonitoring)
		fmt.Println("VM Monitoring: ..", css.Configuration.DasConfig.VmMonitoring)
		fmt.Println("Fail over level:- ", css.Configuration.DasConfig.FailoverLevel)

		s := reflect.ValueOf(css.Summary).Elem()
		typeOfT := s.Type()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			if "CurrentEVCModeKey" == typeOfT.Field(i).Name {
				evcModeKey, _ := f.Interface().(string)
				if evcModeKey != "" {
					fmt.Println("EvcModeKey:", strings.TrimSpace(evcModeKey))
				} else {
					fmt.Println("EvcModeKey:", "Disabled")
				}

			}
			if "NumVmotions" == typeOfT.Field(i).Name {
				numVmotions, _ := f.Interface().(int)
				fmt.Println("numVmotions:", numVmotions)
			}
		}
		s = reflect.ValueOf(css.ConfigurationEx).Elem()
		typeOfT = s.Type()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			if "DpmConfigInfo" == typeOfT.Field(i).Name {
				clusterDmpInfo, _ := f.Interface().(*types.ClusterDpmConfigInfo)
				if clusterDmpInfo != nil {
					if clusterDmpInfo.Enabled != nil {
						fmt.Println("dpm config enabled", *clusterDmpInfo.Enabled)
					}
				}
			}
			if "DrsConfig" == typeOfT.Field(i).Name {
				//need to check why it is not type casting with *types.ClusterDrsConfigInfo
				clusterDrsInfo, _ := f.Interface().(types.ClusterDrsConfigInfo)

				if clusterDrsInfo.Enabled != nil {
					fmt.Println("drs config enabled", *clusterDrsInfo.Enabled)
				}
				fmt.Println("Default vmotion rate...", clusterDrsInfo.VmotionRate)
				fmt.Println("Default defaultVmBehavior ...", clusterDrsInfo.DefaultVmBehavior)
			}

		}
	}
}

func pritnDatacenterName(ctx context.Context, pc *property.Collector, folder mo.Folder) mo.Folder {

	var folderNew mo.Folder
	for {
		err := pc.RetrieveOne(ctx, *folder.Parent, []string{"parent"}, &folderNew)
		if err != nil {
			fmt.Println("Error in getting datacenter", err)
		}
		if folderNew.Parent.Type == "Datacenter" {
			break
		} else {
			folder = folderNew
		}
	}
	return folderNew
}

// HostSystems retrieves host systems information
func HostSystems(ctx context.Context, c *govmomi.Client) ([]mo.HostSystem, error) {
	var hostSystems []mo.HostSystem
	hostSystem := []string{hostSystemObj}
	hostSysProps := []string{summaryProp, runtimeProp, configProp, hardwareProp, datastoreProp}

	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, hostSystem, true)
	if err != nil {
		return nil, err
	}

	err = v.Retrieve(ctx, hostSystem, hostSysProps, &hostSystems)
	if err != nil {
		return nil, err
	}

	err = v.Destroy(ctx)
	if err != nil {
		return nil, err
	}

	return hostSystems, nil
}

func main() {
	hostObj := Host{}
	// hostObj.Address = "10.2.40.120"
	// hostObj.Username = "dtitsupport247\\agentteam"
	// hostObj.Password = "agent@123"

	// Esxi Host Info
	hostObj.Address = "10.4.129.10"
	hostObj.Username = "root"
	hostObj.Password = "Admin@123"

	// hostObj.Address = "10.4.129.13"
	// hostObj.Username = "root"
	// hostObj.Password = "Admin@123"

	// Vcenter Host Info
	// hostObj.Address = "10.4.129.11"
	// hostObj.Username = "administrator@vcenter.local"
	// hostObj.Password = "Admin@123"

	ctx, _ := CancelHandler()

	fmt.Println("Connecting to VMWare host")
	conn, err := NewClient(ctx, hostObj)
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully connected to VMWare host")
	fmt.Println("Getting Host Information")
	hostSys, err := HostSystems(ctx, conn)
	if nil != err {
		fmt.Println(err)
		return
	}

	//AuthorizationManager(ctx, conn)
	ClusterComputeResource(ctx, conn, hostSys)

}

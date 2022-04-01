package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
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
)

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

//Engg - 10.2.40.10 , 10.2.40.200
func main() {
	ctx, _ := CancelHandler()
	conn := clientConnection(ctx)
	hs, err := HostSystems(ctx, conn)
	if err != nil || len(hs) == 0 {
		fmt.Println("Host is Nil")
		return
	}
	// fmt.Println(err)
	// h, err := returnHost("10.2.40.200", hs)
	// fmt.Println("Printing Network Info : ")
	// getStandardVirtualSwitch(h)
	// CollectNetworkInfo(ctx, conn, h.Network)
	h, err := returnHost("10.4.129.13", hs)
	fmt.Println("Printing Network Info : ")
	getStandardVirtualSwitch(h)
	CollectNetworkInfo(ctx, conn, h.Network)
}

func clientConnection(ctx context.Context) *govmomi.Client {
	hostObj := Host{}
	// hostObj.Address = "10.2.40.120"
	// hostObj.Username = "dtitsupport247\\agentteam"
	// hostObj.Password = "agent@123"

	hostObj.Address = "10.4.129.13"
	hostObj.Username = "root"
	hostObj.Password = "Admin@123"

	fmt.Println("Connecting to VMWare host")
	conn, err := NewClient(ctx, hostObj)

	if nil != err {
		fmt.Println(err)
		return nil
	}
	fmt.Println("Successfully connected to VMWare host")

	return conn
}

// func Hosts(ctx context.Context, conn *govmomi.Client) []mo.HostSystem {

// 	hostSys, err := HostSystems(ctx, conn)
// 	if nil != err {
// 		fmt.Println(err)
// 		return nil
// 	}

// 	fmt.Println(len(hostSys))
// 	return hostSys
// }

func HostSystems(ctx context.Context, c *govmomi.Client) ([]mo.HostSystem, error) {
	fmt.Println("Getting Host Information")
	var hostSystems []mo.HostSystem
	hostSystem := []string{hostSystemObj}
	hostSysProps := []string{summaryProp, runtimeProp, configProp, hardwareProp, datastoreProp, "network"}

	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, hostSystem, true)
	if err != nil {
		return nil, err
	}

	err = v.Retrieve(ctx, hostSystem, hostSysProps, &hostSystems)
	if err != nil {
		return nil, err
	}
	defer v.Destroy(ctx)
	// err = v.Destroy(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return hostSystems, nil
}

// NewClient creates session with vmware server
func NewClient(ctx context.Context, h Host) (*govmomi.Client, error) {
	userInfo := url.UserPassword(h.Username, h.Password)
	raw := "https://" + userInfo.String() + "@" + h.Address + "/sdk"
	url, err := url.Parse(raw)
	//url, err := soap.ParseURL(fmt.Sprintf(esxiAddr, h.Username, h.Password, h.Address))
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

func returnHost(ip string, sys []mo.HostSystem) (mo.HostSystem, error) {

	for _, hs := range sys {
		ips := GetHostIPs(hs.Config.Network)
		for j := range ips {
			if ips[j] == ip {
				return hs, nil
			}
		}

	}
	return mo.HostSystem{}, fmt.Errorf("Error")
	//return sys[0]
}

//GetHostIPs returns the Managed IP address of the ESXi Host
func GetHostIPs(hostNetworkInfo *types.HostNetworkInfo) []string {

	if nil == hostNetworkInfo {
		return make([]string, 0)
	}
	var ipFound bool
	IPs := getHostAllIps(hostNetworkInfo, &ipFound)

	if !ipFound && len(hostNetworkInfo.ConsoleVnic) > 0 {
		for _, val := range hostNetworkInfo.ConsoleVnic {
			if nil != val.Spec.Ip && "" != val.Spec.Ip.IpAddress {
				IPs = append(IPs, val.Spec.Ip.IpAddress)
			}
		}
	}
	return IPs
}

//getHostAllIps gets all the host Ips along with the managed IP
//sets the IPfound value to true if any host related IP is found
func getHostAllIps(hostNetworkInfo *types.HostNetworkInfo, ipFound *bool) []string {
	var IPs = make([]string, 0)
	if len(hostNetworkInfo.Vnic) > 0 {
		for _, val := range hostNetworkInfo.Vnic {
			if nil != val.Spec.Ip && "" != val.Spec.Ip.IpAddress {
				IPs = append(IPs, val.Spec.Ip.IpAddress)
				*ipFound = true
			}
		}
	}
	return IPs
}

func networks(ctx context.Context, c *vim25.Client) ([]mo.Network, error) {
	fmt.Println()
	// Create a view of Network types
	m := view.NewManager(c)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Network"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.Network.html
	var networks []mo.Network
	err = v.Retrieve(ctx, []string{"Network"}, nil, &networks)
	if err != nil {
		return nil, err
	}

	for _, net := range networks {
		fmt.Println(net.Self.Type)
		fmt.Println(net.Reference().Type)
		if len(net.Host) > 0 {
			fmt.Println(net.Host[0].Value)
			fmt.Println(net.Host[0].String())
			fmt.Println(net.Host[0].Reference().String())
			fmt.Println(net.Host[0].Type)
		}
		fmt.Printf("%s: %s\n", net.Name, net.Reference())
	}

	return networks, nil
}

func CollectNetworkInfo(ctx context.Context, c *govmomi.Client, networkRef []types.ManagedObjectReference) {
	m := view.NewManager(c.Client)
	v, err := m.CreateListView(ctx, networkRef)
	defer v.Destroy(ctx)
	if err != nil {
		return
	}
	dss := make([]mo.Network, 0)

	for _, hd := range networkRef {
		err = v.Properties(ctx, hd.Reference(), []string{"summary", "host", "name"}, &dss)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	// err = v.Destroy(ctx)
	// if err != nil {
	// 	return
	// }
	//fmt.Println("Printing Network Info : ")
	for _, d := range dss {
		if d.Self.Type == "DistributedVirtualPortgroup" {
			//temp := types.DVPortgroupConfigInfo{}
			//fmt.Println("PK config::", d.Config.DefaultPortConfig)
			//fmt.Println("Host Having Distributed Vswitch %+v", d.Host)
			distributedPortGroup := mo.DistributedVirtualPortgroup{}

			err = v.Properties(ctx, d.Reference(), []string{"config"}, &distributedPortGroup)
			//err = v.Properties(ctx, d.Reference(), []string{"config"}, &temp)
			if err != nil {
				fmt.Println(err)
				continue
			}
			testData := distributedPortGroup.Config.DefaultPortConfig
			fmt.Println("Distributed Switch data here")
			type12, ok := testData.(*(types.VMwareDVSPortSetting))
			if ok {
				type34, ok2 := type12.Vlan.(*(types.VmwareDistributedVirtualSwitchVlanIdSpec))
				if ok2 != false {
					fmt.Println("Distributed Switch VlanId", type34.VlanId)
				}
			}

			fmt.Println("###################### Distributed Switch End ##############################")
		}
	}
}

//https://asvignesh.in/vsphere-tagging-feature-identify-uplink-port-in-dswitch/
func containsUplinkTag(tags []types.Tag) bool {
	for _, t := range tags {
		if strings.Contains(t.Key, "UPLINK") {
			return true
		}
	}
	return false
}

//SwitchType Virtual Switch can be of multiple types like Standard,distributed,etc
type SwitchType string

const (
	// DistributedSwitchType provides a centralized interface from which one can configure,
	//monitor and administer virtual machine access switching for the entire data center.
	DistributedSwitchType = SwitchType("DistributedVirtualSwitch")
	// StandardSwitchType provides network connectivity to hosts and virtual machines
	StandardSwitchType = SwitchType("StandardVirtualSwitch")
)

// VirtualSwitch is a software entity to which multiple virtual network
//adapters can connect to create a virtual network. It can also be bridged to a physical network.
type VirtualSwitch struct {
	//Type of Virtual Switch
	Type SwitchType `json:"type,omitempty" cql:"type"`
	//Key is a unique identifier for Virtual Switch
	Key  string `json:"key,omitempty" cql:"key"`
	Name string `json:"name,omitempty" cql:"name"`
	//PhysicalNICs are a set of physical network adapters associated with this bridge.
	PhysicalNICs []PhysicalNIC `json:"physicalNICs,omitempty" cql:"physical_nics"`
	//PortGroups are a list of port groups configured for this virtual switch.
	PortGroups []PortGroup `json:"portGroups,omitempty" cql:"portgroups"`
}

//PortGroup  are used to group virtual network adapters on a virtual switch,
//associating them with networks and network policies.
type PortGroup struct {
	//Key is a unique identifier for Port Group
	Key  string `json:"key,omitempty" cql:"key"`
	Name string `json:"name,omitempty" cql:"name"`
}

//PhysicalNIC describes the physical network adapter
//as seen by the primary operating system.
type PhysicalNIC struct {
	//Key is a unique identifier for PhysicalNIC
	Key  string `json:"key,omitempty" cql:"Key"`
	Name string `json:"name,omitempty" cql:"name"`
	//MAC is the media access control (MAC) address of the physical network adapter.
	MAC string `json:"mac,omitempty" cql:"mac"`
}

type SwitchsWithPortGroups struct {
	VSwitchName   string `json:"vSwitchName,omitempty"`
	PortGroupName string `json:"portGroupName,omitempty"`
}

func getStandardVirtualSwitch(h mo.HostSystem) {
	var portgroups = make(map[string]string, len(h.Config.Network.Portgroup))
	for _, portGroup := range h.Config.Network.Portgroup {
		//portgroups[portGroup.Key] = portGroup.Vswitch
		portgroups["VSwitch Name"] = portGroup.Spec.VswitchName
		portgroups["name"] = portGroup.Spec.Name
		portgroups["VlanID"] = strconv.FormatInt(int64(portGroup.Spec.VlanId), 10)

	}
	fmt.Println("Standard Switch data:  ", portgroups)
	fmt.Println("###################### getStandardVirtualSwitch End ##############################")
}

func getNetworkPortGroups(Network *types.HostNetworkInfo) (vms []SwitchsWithPortGroups) {
	vms = make([]SwitchsWithPortGroups, len(Network.Portgroup))
	for _, portgrp := range Network.Portgroup {
		vms = append(vms, SwitchsWithPortGroups{
			PortGroupName: portgrp.Spec.Name,
			VSwitchName:   portgrp.Spec.VswitchName,
		})
	}
	return
}

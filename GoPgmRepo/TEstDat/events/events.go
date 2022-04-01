package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"gitlab.connectwisedev.com/platform/platform-api-model/clients/model/Golang/resourceModel/virtualization/inventory"
)

const (
	apiProtocol       = "https://"
	virtualMachineObj = "VirtualMachine"
	summaryProp       = "summary"
	runtimeProp       = "runtime"
	configProp        = "config"
	hardwareProp      = "hardware"
	datastoreProp     = "datastore"
	guestProp         = "guest"
	snapshotProp      = "snapshot"
	toolsNotRunning   = "guestToolsNotRunning"
)

var count int

func handleEvent(ref types.ManagedObjectReference, events []types.BaseEvent) (err error) {
	for _, event := range events {
		a, b := eventMan.EventCategory(c, event)
		fmt.Printf("%+v %+v\n", a, b)
		eventType := reflect.TypeOf(event).String()
		fmt.Printf("Event found of type %s\n", eventType)
		fmt.Printf("Event found of type %+v\n", event.GetEvent())
		fmt.Printf("Full message body %+v\n", event.GetEvent().FullFormattedMessage)

		//		fmt.Printf("Host Name %+v\n", event.GetEvent().Host.Name)
		fmt.Printf("Host Name %+v\n", event.GetEvent().Host)
		fmt.Printf("Vm %+v\n", event.GetEvent().Vm)
		fmt.Printf("DynamicData %+v\n", event.GetEvent().DynamicData)
		fmt.Printf("Datacenter %+v\n", event.GetEvent().Datacenter)
		fmt.Printf("ComputeResource %+v\n", event.GetEvent().ComputeResource)
	}

	return nil
}

var eventMan *event.Manager
var c context.Context
var host mo.HostSystem
var client1 *govmomi.Client

func Eventf(ctx context.Context, client *govmomi.Client, hs mo.HostSystem) {
	// Selecting default datacenter
	//finder := find.NewFinder(client.Client, true)
	// dc, err := finder.DefaultDatacenter(ctx)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	// 	os.Exit(1)
	// }
	refs := []types.ManagedObjectReference{hs.Reference()}
	host = hs
	// Setting up the event manager
	eventManager := event.NewManager(client.Client)
	eventMan = eventManager
	c = ctx
	client1 = client
	for {
		err := eventManager.Events(ctx, refs, 1, true, true, handleEventFormatted)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	}
}

func handleEventFormatted(ref types.ManagedObjectReference, events []types.BaseEvent) (err error) {

	for _, event := range events {
		switch event.(type) {
		case *types.VmCreatedEvent:
			e := event.GetEvent()
			printFullFormattedMessage("VmCreatedEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			vms, err := VirtualMachines(c, client1, e.Vm.Vm.Reference())
			if err != nil {
				fmt.Println("Error in fetching vminfo", err)
			} else {
				for _, vm := range vms {
					if vm.Config != nil {
						fmt.Println("Vm name", vm.Config.Name)
						fmt.Println("guest full name", vm.Config.GuestFullName)
						fmt.Println("Vm Instance ID", vm.Config.InstanceUuid)
					}
				}
			}
			vmsArray, err := VirtualMachinesNew(c, client1, host.Vm, host.Config.Network)
			if err != nil {
				fmt.Println("Error in fetching vminfo", err)
			} else {
				for _, vm := range vmsArray {
					if vm.Config != nil {
						fmt.Println("Vm name after created", vm.Config.Name)
					}
				}
			}

		case *types.VmRenamedEvent:
			e := event.GetEvent()
			printFullFormattedMessage("VmCreatedEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			t := e.Vm.Name
			fmt.Println("Vm Name", t)
			fmt.Println("Vm Data", e.Vm)
			vms, _ := VirtualMachines(c, client1, e.Vm.Vm.Reference())
			for _, vm := range vms {
				if vm.Config != nil {
					fmt.Println("Vm name", vm.Config.Name)
					fmt.Println("guest full name", vm.Config.GuestFullName)
					fmt.Println("Vm Instance ID", vm.Config.InstanceUuid)
				}
			}
		case *types.VmRemovedEvent:
			count = count + 1
			e := event.GetEvent()
			printFullFormattedMessage("VmRemovedEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			fmt.Println("Need to send remaing vms data.....", count)
		case *types.VmPoweredOnEvent:
			e := event.GetEvent()
			printFullFormattedMessage("VmPoweredOnEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			vms, _ := VirtualMachines(c, client1, e.Vm.Vm.Reference())
			for _, vm := range vms {
				if vm.Config != nil {
					fmt.Println("Vm Instance ID", vm.Config.InstanceUuid)
					fmt.Println("Vm State", vm.Summary.Runtime.PowerState)
					fmt.Println("VM overall Status", vm.Summary.OverallStatus)
					fmt.Println("guest full name", vm.Config.GuestFullName)
				}
			}
		case *types.VmSuspendedEvent:
			e := event.GetEvent()
			printFullFormattedMessage("VmSuspendedEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			vms, _ := VirtualMachines(c, client1, e.Vm.Vm.Reference())
			for _, vm := range vms {
				if vm.Config != nil {
					fmt.Println("Vm Instance ID", vm.Config.InstanceUuid)
					fmt.Println("Vm State", vm.Summary.Runtime.PowerState)
					fmt.Println("VM overall Status", vm.Summary.OverallStatus)
					fmt.Println("guest full name", vm.Config.GuestFullName)
				}
			}
		case *types.VmResumingEvent:
			e := event.GetEvent()
			printFullFormattedMessage("VmResumingEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			vms, _ := VirtualMachines(c, client1, e.Vm.Vm.Reference())
			for _, vm := range vms {
				if vm.Config != nil {
					fmt.Println("Vm Instance ID", vm.Config.InstanceUuid)
					fmt.Println("Vm State", vm.Summary.Runtime.PowerState)
					fmt.Println("VM overall Status", vm.Summary.OverallStatus)
					fmt.Println("guest full name", vm.Config.GuestFullName)
				}
			}
		case *types.VmPoweredOffEvent:
			e := event.GetEvent()
			printFullFormattedMessage("VmPoweredOffEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			vms, _ := VirtualMachines(c, client1, e.Vm.Vm.Reference())
			for _, vm := range vms {
				if vm.Config != nil {
					fmt.Println("Vm Instance ID", vm.Config.InstanceUuid)
					fmt.Println("Vm State", vm.Summary.Runtime.PowerState)
					fmt.Println("VM overall Status", vm.Summary.OverallStatus)
					fmt.Println("guest full name", vm.Config.GuestFullName)
				}
			}
		case *types.VmResettingEvent:
			e := event.GetEvent()
			printFullFormattedMessage("VmResettingEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)

		case *types.VmReconfiguredEvent:
			e := event.GetEvent()
			printFullFormattedMessage("VmReconfiguredEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			fmt.Println("Its general event if we are changing cpu, memory, network setting etc")

		case *types.VmRelocatedEvent:
			e := event.GetEvent()
			printFullFormattedMessage("VmRelocatedEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			fmt.Println("When we are doing Vmotion. its fromvcenter event")

		case *types.EnteredMaintenanceModeEvent:
			e := event.GetEvent()
			printFullFormattedMessage("EnteredMaintenanceModeEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			fmt.Println("Host entered into maintenanceMode")
		case *types.ExitMaintenanceModeEvent:
			e := event.GetEvent()
			printFullFormattedMessage("ExitMaintenanceModeEvent", e.CreatedTime, e.UserName, e.FullFormattedMessage, e.Host.Name, e.Vm.Name)
			fmt.Println("Host is exit from maintenance mode")
		default:
			fmt.Println("------------------------------------------------------------------------")
			fmt.Printf("Event of type %v Ignored!!", reflect.TypeOf(event).String())
			fmt.Println("\n------------------------------------------------------------------------\n")
			//e.logger.WithFields(logrus.Fields{"type": reflect.TypeOf(event).String()}).Debug("Event ignored")
		}
	}
	return nil
}

func printFullFormattedMessage(eventType string, time time.Time, username string, msg string, host string, vmName string) {
	fmt.Println("Event Type : ", eventType)
	fmt.Println("Created Time : ", time)
	fmt.Println("Username : ", username)
	fmt.Println("Host Name : ", host)
	fmt.Println("VM Name : ", vmName)
	fmt.Println("Full Formatted Message: ", msg)
	fmt.Println("------------------------------------------------------------------------\n")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//uri := "https://administrator@vsphere.local:Admin!23@10.160.164.46"
	//u, err := getAPIPath("10.4.129.10", "root", "Admin@123") //soap.ParseURL(uri)
	//u, err := getAPIPath("10.4.129.13", "root", "Admin@123") //soap.ParseURL(uri)
	u, err := getAPIPath("10.4.129.11", "administrator@vcenter.local", "Admin@123")
	//u, err := getAPIPath("10.2.40.120", "dtitsupport247\\agentteam", "agent@123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Creating client connection")
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Client Connection Success")

	hostSys, err := HostSystems(ctx, c)
	if nil != err {
		fmt.Println(err)
		return
	}
	for _, val := range hostSys {
		fmt.Println("For host", val.Summary.Config.Name)
		VirtualMachinesNew(ctx, c, val.Vm, val.Config.Network)
		//Eventf(ctx, c, val)
	}
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

// HostSystems retrieves host systems information
func HostSystems(ctx context.Context, c *govmomi.Client) ([]mo.HostSystem, error) {
	var hostSystems []mo.HostSystem
	hostSystem := []string{"HostSystem"}
	var hostSysProps []string
	//hostSysProps := []string{summaryProp, runtimeProp, configProp, hardwareProp, datastoreProp}
	hostSysProps = []string{"summary", "vm", "config"}
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

//VirtualMachines retrieves VMs information as MOB
func VirtualMachines(ctx context.Context, c *govmomi.Client, virtualMachinesRef types.ManagedObjectReference) ([]mo.VirtualMachine, error) {
	var virtualMachines []mo.VirtualMachine
	virtualMachinesProps := []string{summaryProp, configProp, guestProp}

	m := view.NewManager(c.Client)
	err := m.Properties(ctx, virtualMachinesRef, virtualMachinesProps, &virtualMachines)
	if err != nil {
		return nil, err
	}
	// v, err := m.CreateListView(ctx, virtualMachinesRef)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, hd := range virtualMachinesRef {
	// 	err = v.Properties(ctx, hd.Reference(), virtualMachinesProps, &virtualMachines)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// err = v.Destroy(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return virtualMachines, nil
}

func VirtualMachinesNew(ctx context.Context, c *govmomi.Client, virtualMachinesRef []types.ManagedObjectReference, Network *types.HostNetworkInfo) ([]mo.VirtualMachine, error) {
	var virtualMachines []mo.VirtualMachine
	virtualMachinesProps := []string{summaryProp, configProp, guestProp, datastoreProp, snapshotProp}

	m := view.NewManager(c.Client)
	v, err := m.CreateListView(ctx, virtualMachinesRef)
	if err != nil {
		return nil, err
	}

	for _, hd := range virtualMachinesRef {
		err = v.Properties(ctx, hd.Reference(), virtualMachinesProps, &virtualMachines)
		if err != nil {
			return nil, err
		}
	}
	err = v.Destroy(ctx)
	if err != nil {
		return nil, err
	}

	npgData := getNetworkPortGroups(Network)

	for _, val := range virtualMachines {
		netData := getNetworks(val.Guest.Net)
		fmt.Println("Guest Full Name", val.Guest.GuestFullName)
		fmt.Println("Vmtools status running", val.Guest.ToolsRunningStatus)

		if netData != nil {
			for _, nets := range netData {
				fmt.Println("Ip Addess", nets.IPv4)
				for _, val2 := range npgData {
					if val2.PortGroupName == nets.PortGroup {
						fmt.Println("vswitch name", val2.VSwitchName)
					}
				}
			}
		}

	}
	return virtualMachines, nil
}

func getNetworks(gni []types.GuestNicInfo) []inventory.Network {
	networks := make([]inventory.Network, 0, len(gni))
	for _, nw := range gni {
		network := inventory.Network{}
		network.IPEnabled = nw.Connected
		network.MacAddress = nw.MacAddress
		network.PortGroup = nw.Network
		if nw.IpConfig != nil {
			for _, ip := range nw.IpConfig.IpAddress {

				if strings.Contains(ip.IpAddress, ":") {
					network.IPv6 = append(network.IPv6, net.ParseIP(ip.IpAddress))
				} else {
					network.IPv4 = append(network.IPv4, net.ParseIP(ip.IpAddress))
				}
			}
			networks = append(networks, network)
		}
	}
	return networks
}

type SwitchsWithPortGroups struct {
	VSwitchName   string `json:"vSwitchName,omitempty"`
	PortGroupName string `json:"portGroupName,omitempty"`
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

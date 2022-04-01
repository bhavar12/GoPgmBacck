package main

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/license"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	esxiAddr               = "https://%s:%s@%s/sdk"
	hostSystemObj          = "HostSystem"
	summaryProp            = "summary"
	runtimeProp            = "runtime"
	configProp             = "config"
	hardwareProp           = "hardware"
	datastoreProp          = "datastore"
	networkProp            = "network"
	licensableResourceProp = "licenses"
	apiProtocol            = "https://"
	//[]string{"licenses"}
)

var (
	networkName string
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

// HostSystems retrieves host systems information
func HostSystems(ctx context.Context, c *govmomi.Client) ([]mo.HostSystem, error) {
	var hostSystems []mo.HostSystem
	hostSystem := []string{hostSystemObj}
	hostSysProps := []string{summaryProp, runtimeProp, configProp, hardwareProp, datastoreProp, networkProp}

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

	// Vcenter Host Info for qa env
	// hostObj.Address = "10.2.40.120"
	// hostObj.Username = "dtitsupport247\\ajay.d"
	// hostObj.Password = "aj@y@098"
	ctx, cancelFn := CancelHandler()
	fmt.Println(cancelFn)
	//fmt.Println("Connecting to VMWare host")
	conn, err := NewClient(ctx, hostObj)
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully connected to VMWare host")
	fmt.Println("Getting Host Information")
	hostSys, err := HostSystems(ctx, conn)
	//fmt.Printf("Lenght Of host Sys %+v", hostSys[0])
	if nil != err {
		fmt.Println(err)
		return
	}
	//DataStoreNew(ctx, conn, hostSys)
	for _, val := range hostSys {
		//printHostSummary(val.Summary)
		//printHostConfigInfo(val.Config)
		fmt.Println("Gor host", val.Hardware.SystemInfo.Uuid)
		printHostHardwareInfo(val.Hardware)
		//printBiosData(val.Runtime, val.Config)
		//fmt.Println("License Info", val.LicensableResource.Resource)
		//printNetwork(ctx, conn, val.Network)
		//printDataStoreInfo(ctx, conn, val.Datastore)
		//DataStore(ctx, conn, val.Datastore)

	}
}

//DataStore will print the Lun info
func DataStore(ctx context.Context, c *govmomi.Client, datastore []types.ManagedObjectReference) {

	pc := property.DefaultCollector(c.Client)
	var dss []mo.Datastore

	pc.Retrieve(ctx, datastore, []string{"info", "summary"}, &dss)

	var vmfs *types.HostVmfsVolume
	for _, ds := range dss {
		fmt.Println("Datastore name\t", ds.Summary.Name)
		s := reflect.ValueOf(ds.Info).Elem()
		typeOfT := s.Type()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)

			//fmt.Println(i, typeOfT.Field(i).Name, f.Type(), f.Interface())
			if "Vmfs" == typeOfT.Field(i).Name {
				vmfs, _ = f.Interface().(*types.HostVmfsVolume)
				if vmfs != nil {
					for _, val := range vmfs.Extent {
						fmt.Println("Number of Partition\t", val.Partition)
						fmt.Println("Disk Name\t", val.DiskName)
					}
					fmt.Println("vmfs UUID\t", vmfs.Uuid)
				}
			}

		}
	}
}

//DataStoreNew will Collect the Lun information against host
func DataStoreNew(ctx context.Context, c *govmomi.Client, host []mo.HostSystem) {

	pc := property.DefaultCollector(c.Client)

	for _, objHost := range host {
		var dss []mo.Datastore
		fmt.Println("For the host", objHost.Name)
		pc.Retrieve(ctx, objHost.Datastore, []string{"info", "summary"}, &dss)

		var vmfs *types.HostVmfsVolume
		for _, ds := range dss {
			fmt.Println("Data Store Name\t", ds.Summary.Name)
			fmt.Println(reflect.TypeOf(ds.Info))
			s := reflect.ValueOf(ds.Info).Elem()

			typeOfT := s.Type()

			for i := 0; i < s.NumField(); i++ {

				f := s.Field(i)
				//fmt.Println(i, typeOfT.Field(i).Name, f.Type(), f.Interface())
				if "Vmfs" == typeOfT.Field(i).Name {
					vmfs, _ = f.Interface().(*types.HostVmfsVolume)
					//fmt.Println(vmfs.Extent)
					if vmfs != nil {
						for _, val := range vmfs.Extent {
							fmt.Println("Number of Partition\t", val.Partition)
							fmt.Println("Disk Name\t", val.DiskName)
						}
						fmt.Println("UUID\t", vmfs.Uuid)
					}
				}

			}
		}
	}
}

func LicenseInfo(ctx context.Context, c *govmomi.Client) {
	licMgr := license.NewManager(c.Client)
	//fmt.Println(licMgr.List(ctx))
	//info, _ := licMgr.List(ctx)
	assmgr, _ := licMgr.AssignmentManager(ctx)
	info, _ := assmgr.QueryAssigned(ctx, "")

	for _, obj := range info {
		fmt.Println("Name", obj.AssignedLicense.Name)
		fmt.Println("Host Name", obj.EntityDisplayName)
		fmt.Println("License Key", obj.AssignedLicense.LicenseKey)
	}

}

func printDataStoreInfo(ctx context.Context, cli *govmomi.Client, datastore []types.ManagedObjectReference) {
	m := view.NewManager(cli.Client)
	v, err := m.CreateListView(ctx, datastore)
	if err != nil {
		fmt.Println("Error in view", err)
	}

	ds := make([]mo.Datastore, 0)
	for _, hd := range datastore {
		err = v.Properties(ctx, hd.Reference(), []string{"summary", "info"}, &ds)
		if err != nil {
			fmt.Println("Error in properties", err)
		}
	}
	err = v.Destroy(ctx)
	if err != nil {
		fmt.Println("Error in destory", err)
	}
	fmt.Println("\n####### Data Store Info #######")
	for _, val := range ds {
		fmt.Println("\nName\t", val.Summary.Name)
		fmt.Println("Type\t", val.Summary.Type)
		fmt.Println("Free Space\t", val.Summary.FreeSpace)
		fmt.Println("Capacity\t", val.Summary.Capacity)
		fmt.Println("usedspace\t", val.Summary.Capacity-val.Summary.FreeSpace)
		data := val.Info.GetDatastoreInfo()
		fmt.Println("Url Data\t", data.Url)
		//fmt.Println("vmfs\t", data.Vmfs)
	}
	// var vmfs *types.HostVmfsVolume
	// for _, ds := range ds {
	// 	fmt.Println("Data Store Name\t", ds.Summary.Name)
	// 	//fmt.Println("Multiple host access", *ds.Summary.MultipleHostAccess)
	// 	fmt.Println(reflect.TypeOf(ds.Info))
	// 	s := reflect.ValueOf(ds.Info).Elem()

	// 	typeOfT := s.Type()

	// 	for i := 0; i < s.NumField(); i++ {

	// 		f := s.Field(i)
	// 		//fmt.Println(i, typeOfT.Field(i).Name, f.Type(), f.Interface())
	// 		if "Vmfs" == typeOfT.Field(i).Name {
	// 			vmfs, _ = f.Interface().(*types.HostVmfsVolume)
	// 			//fmt.Println(vmfs.Extent)
	// 			if vmfs != nil {
	// 				for _, val := range vmfs.Extent {
	// 					fmt.Println("Number of Partition\t", val.Partition)
	// 					fmt.Println("Disk Name\t", val.DiskName)
	// 				}
	// 				fmt.Println("UUID\t", vmfs.Uuid)
	// 			}
	// 		}

	// 	}
	// }

}
func printNetwork(ctx context.Context, cli *govmomi.Client, network []types.ManagedObjectReference) {
	m := view.NewManager(cli.Client)
	v, err := m.CreateListView(ctx, network)
	if err != nil {
		fmt.Println("Error in view", err)
	}

	ds := make([]mo.Network, 0)
	for _, hd := range network {
		err = v.Properties(ctx, hd.Reference(), []string{"summary", "name"}, &ds)
		if err != nil {
			fmt.Println("Error in properties", err)
		}
	}

	err = v.Destroy(ctx)
	if err != nil {
		fmt.Println("Error in destory", err)
	}
	for _, val := range ds {
		fmt.Println("\nNetwork Name\t", val.Name)
	}
}
func printBiosData(hostruntime types.HostRuntimeInfo, config *types.HostConfigInfo) {
	sysBootTimeLayout := "2006-01-02 15:04:05.999999999"
	pastDate := hostruntime.BootTime
	timeZoneInString := config.DateTimeInfo.TimeZone.Key
	fmt.Println("\n######  Boot Time #########")
	fmt.Println("Bios Date\t", pastDate)
	fmt.Println("Time Zone\t", timeZoneInString)
	if timeZoneInString == "UTC" {
		pastDateStringFormat := pastDate.Format(sysBootTimeLayout)
		fmt.Println("Formated date in string \t", pastDateStringFormat)
		loc, err := time.LoadLocation(timeZoneInString)
		if err != nil {
			fmt.Println("Error in Load Location", err)
		}
		newDate, err := time.ParseInLocation(sysBootTimeLayout, pastDateStringFormat, loc)
		if err != nil {
			fmt.Println("Error in parsing", err)
		}
		fmt.Println("Final Date in UTC\t", newDate.UTC())
	}
}
func printHostHardwareInfo(hostHardwareInfo *types.HostHardwareInfo) {
	if nil == hostHardwareInfo {
		return
	}
	//printHostBIOSInfo(hostHardwareInfo.BiosInfo)
	printHostSystemInfo(hostHardwareInfo.SystemInfo)
	//printHostPciInfo(hostHardwareInfo.PciDevice)

}

func printHostPciInfo(pciDevice []types.HostPciDevice) {
	for _, data := range pciDevice {
		fmt.Println("DeviceID", data.DeviceId)
		fmt.Println("DeviceName", data.DeviceName)
		fmt.Println("Vendor Name", data.VendorName)
	}
}
func printHostSystemInfo(hostSystemInfo types.HostSystemInfo) {
	fmt.Println("\n===== Host System Info =====")
	fmt.Println("Vendor\t:", hostSystemInfo.Vendor)
	fmt.Println("Model\t:", hostSystemInfo.Model)
	fmt.Println("Uuid\t:", hostSystemInfo.Uuid)
	fmt.Println("SerialNumber\t:", hostSystemInfo.SerialNumber)
}
func printHostBIOSInfo(hostBIOSInfo *types.HostBIOSInfo) {
	if nil == hostBIOSInfo {
		return
	}
	fmt.Println("\n===== Host BIOS Info =====")
	fmt.Println("BiosVersion\t:", hostBIOSInfo.BiosVersion)
	fmt.Println("ReleaseDate\t:", hostBIOSInfo.ReleaseDate)
	fmt.Println("Vendor\t:", hostBIOSInfo.Vendor)

}

func printHostConfigInfo(hostCfgInfo *types.HostConfigInfo) {
	if nil == hostCfgInfo {
		return
	}
	printHostNetworkInfo(hostCfgInfo.Network)

}

func printHostNetworkInfo(hostNetworkInfo *types.HostNetworkInfo) {

	if nil == hostNetworkInfo {
		return
	}
	fmt.Println("\n===== Host Network Info =====")
	for _, val := range hostNetworkInfo.Pnic {
		fmt.Println("Device\t:", val.Device)
		fmt.Println("Key\t:", val.Key)
		fmt.Println("MAC\t:", val.Mac)
	}
	for _, val := range hostNetworkInfo.Vnic {
		fmt.Println("IP Adress\t:", val.Spec.Ip.IpAddress)
		fmt.Println("SubNet Mask\t:", val.Spec.Ip.SubnetMask)
	}
	testDNS := hostNetworkInfo.DnsConfig.GetHostDnsConfig()
	fmt.Println("Domain Name", testDNS.DomainName)
	defaultGateway := hostNetworkInfo.IpRouteConfig.GetHostIpRouteConfig()
	fmt.Println("Default Gateway", defaultGateway.DefaultGateway)
}

func printHostSummary(hostSummary types.HostListSummary) {
	fmt.Println("\n===== Host List Summary =====")
	fmt.Println("Management Server Ip\t:", hostSummary.ManagementServerIp)
	fmt.Println("Overall Status\t:", hostSummary.OverallStatus)
	fmt.Println("Reboot Required\t:", hostSummary.RebootRequired)

	printHostConfigSummary(hostSummary.Config)
	printHostHardwareSummary(hostSummary.Hardware)
	printHostListSummaryQuickStats(hostSummary.QuickStats)
	printHostRuntimeInfo(hostSummary.Runtime)

}

func printHostConfigSummary(hostCfgSummary types.HostConfigSummary) {
	fmt.Println("\n===== Host Config Summary =====")
	fmt.Println("Name\t:", hostCfgSummary.Name)
	fmt.Println("Fault Tolerance Enabled\t:", *hostCfgSummary.FaultToleranceEnabled)
	fmt.Println("Vmotion Enabled\t:", hostCfgSummary.VmotionEnabled)

	printHostAboutInfo(hostCfgSummary.Product)

}

func printHostAboutInfo(hostAbtInfo *types.AboutInfo) {
	if nil == hostAbtInfo {
		return
	}
	fmt.Println("\n===== Host About Info =====")
	fmt.Println("Name\t:", hostAbtInfo.Name)
	fmt.Println("Full Name\t:", hostAbtInfo.FullName)
	fmt.Println("Vendor\t:", hostAbtInfo.Vendor)
	fmt.Println("Version\t:", hostAbtInfo.Version)
	fmt.Println("Build\t:", hostAbtInfo.Build)
	fmt.Println("Locale Version\t:", hostAbtInfo.LocaleVersion)
	fmt.Println("Locale Build\t:", hostAbtInfo.LocaleBuild)
	fmt.Println("OS Type\t:", hostAbtInfo.OsType)
	fmt.Println("Product Line Id\t:", hostAbtInfo.ProductLineId)
	fmt.Println("Api Type\t:", hostAbtInfo.ApiType)
	fmt.Println("Api Version\t:", hostAbtInfo.ApiVersion)
	fmt.Println("Instance Uuid\t:", hostAbtInfo.InstanceUuid)
	fmt.Println("License Product Name\t:", hostAbtInfo.LicenseProductName)
	fmt.Println("License Product Version\t:", hostAbtInfo.LicenseProductVersion)
	//hostAbtInfo.L
}

func printHostHardwareSummary(hostHardwareSummary *types.HostHardwareSummary) {
	if nil == hostHardwareSummary {
		return
	}

	fmt.Println("\n===== Host Hardware Summary =====")
	fmt.Println("Vendor\t:", hostHardwareSummary.Vendor)
	fmt.Println("Model\t:", hostHardwareSummary.Model)
	fmt.Println("Uuid\t:", hostHardwareSummary.Uuid)
	fmt.Println("MemorySize\t:", hostHardwareSummary.MemorySize)
	fmt.Println("CpuModel\t:", hostHardwareSummary.CpuModel)
	fmt.Println("CpuMhz\t:", hostHardwareSummary.CpuMhz)
	fmt.Println("NumCpuPkgs\t:", hostHardwareSummary.NumCpuPkgs)
	fmt.Println("NumCpuCores\t:", hostHardwareSummary.NumCpuCores)
	fmt.Println("NumCpuThreads\t:", hostHardwareSummary.NumCpuThreads)
	fmt.Println("NumNics\t:", hostHardwareSummary.NumNics)
	fmt.Println("NumHBAs\t:", hostHardwareSummary.NumHBAs)
}

func printHostListSummaryQuickStats(hostListSummaryQuickStats types.HostListSummaryQuickStats) {
	fmt.Println("\n===== Host List Summary Quick Stats =====")
	fmt.Println("Overall Cpu Usage\t:", hostListSummaryQuickStats.OverallCpuUsage)
	fmt.Println("Overall Memory Usage\t:", hostListSummaryQuickStats.OverallMemoryUsage)
	fmt.Println("Uptime\t:", hostListSummaryQuickStats.Uptime)
}

func printHostRuntimeInfo(hostRuntimeInfo *types.HostRuntimeInfo) {
	if nil == hostRuntimeInfo {
		return
	}

	fmt.Println("\n===== Host Run Time Info =====")
	fmt.Println("ConnectionState\t:", hostRuntimeInfo.ConnectionState)
	fmt.Println("PowerState\t:", hostRuntimeInfo.PowerState)
	fmt.Println("StandbyMode\t:", hostRuntimeInfo.StandbyMode)
	fmt.Println("InMaintenanceMode\t:", hostRuntimeInfo.InMaintenanceMode)
	if nil != hostRuntimeInfo.InQuarantineMode {
		fmt.Println("InQuarantineMode\t:", *hostRuntimeInfo.InQuarantineMode)
	}
	fmt.Println("BootTime\t:", *hostRuntimeInfo.BootTime)
	pastDate := *hostRuntimeInfo.BootTime
	now := time.Now().UTC()
	diff := now.Sub(pastDate.UTC())
	days := float64(diff.Hours() / 24)
	fmt.Printf("Diffrence in days : %d days\n", int(math.Round(days)))
}

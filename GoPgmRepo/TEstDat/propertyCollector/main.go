package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	esxiAddr      = "https://%s:%s@%s/sdk"
	hostSystemObj = "HostSystem"
	aboutInfoObj  = "AboutInfo"
	aboutProp     = "about"
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

// NewClient creates session with vmware server
func NewClient(ctx context.Context, h Host) (*govmomi.Client, error) {
	//url, err := soap.ParseURL(fmt.Sprintf(esxiAddr, h.Username, h.Password, h.Address))
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

func PropertyCollector(ctx context.Context, c *govmomi.Client, hostSys []mo.HostSystem) {
	prColl := property.DefaultCollector(c.Client)

	var valTrue = true
	var valFalse = false
	var objectSpec []types.ObjectSpec
	var obj types.ObjectSpec

	var propertySpec []types.PropertySpec
	var pro types.PropertySpec

	//var typ []types.ManagedObjectReference
	for _, val := range hostSys {
		for _, vm := range val.Vm {
			obj.Obj = vm //val.Self
			obj.Skip = &valFalse
			objectSpec = append(objectSpec, obj)
			//	typ = append(typ, vm)
		}

		//obj.SelectSet = []types.BaseSelectionSpec{}
	}

	pro.All = &valTrue
	pro.Type = "VirtualMachine" //VirtualMachine //HostSystem
	//pro.PathSet = []string{}
	propertySpec = append(propertySpec, pro)

	var filType types.CreateFilter
	filType.Spec.ObjectSet = objectSpec
	filType.Spec.PropSet = propertySpec
	filType.Spec.ReportMissingObjectsInResults = &valTrue

	err := prColl.CreateFilter(ctx, filType)
	fmt.Println(err)
	var version string
	update, err := prColl.WaitForUpdates(ctx, version)
	prColl.CancelWaitForUpdates(ctx)
	if nil == update {
		return
	}

	for {
		version = update.Version

		update, err = prColl.WaitForUpdates(ctx, version)
		defer prColl.CancelWaitForUpdates(ctx)
		if nil == update || nil != err {
			return
		}

		if "" == update.Version || 0 == len(update.FilterSet) {
			return
		}

		for _, filtObj := range update.FilterSet {
			printPropertyDataFormatted(filtObj)
			// fmt.Println(reflect.TypeOf(filtObj).String())
			// fmt.Printf("Object Set : %+v", filtObj.ObjectSet)
			// fmt.Println("Type: ", filtObj.Filter.Type)
			// fmt.Println("Data: ", filtObj.Filter.Reference().Type)
			// fmt.Println("Value: ", filtObj.Filter.Value)
		}
	}

}

// HostSystems retrieves host systems information
func HostSystems(ctx context.Context, c *govmomi.Client) ([]mo.HostSystem, error) {
	var hostSystems []mo.HostSystem
	hostSystem := []string{hostSystemObj}
	var hostSysProps []string
	//hostSysProps := []string{summaryProp, runtimeProp, configProp, hardwareProp, datastoreProp}

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
	// hostObj.Address = "10.4.129.11"
	// hostObj.Username = "administrator@vcenter.local"
	// hostObj.Password = "Pass@123"

	hostObj.Address = "10.4.129.13"
	hostObj.Username = "root"
	hostObj.Password = "Admin@123"

	// hostObj.Address = "10.2.40.120"
	// hostObj.Username = "dtitsupport247\\ajay.d"
	// hostObj.Password = "aj@y@098"
	ctx, _ := CancelHandler()

	fmt.Println("Connecting to VMWare VCenter")
	conn, err := NewClient(ctx, hostObj)
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully connected to VMWare VCenter")

	fmt.Println("Getting VCenter Information")
	hostSys, _ := HostSystems(ctx, conn)
	fmt.Println("Started Property Collector")
	PropertyCollector(ctx, conn, hostSys)

}

func printPropertyDataFormatted(filterObj types.PropertyFilterUpdate) {
	fmt.Println("#########################################")
	for _, obj := range filterObj.ObjectSet {
		fmt.Println("Kind : ", obj.Kind)
		fmt.Println("Object Type : ", obj.Obj.Type)
		fmt.Printf(" Object name : %+v", obj.Obj.Value)
		fmt.Println("\n-----------All Obj Data---------------")
		for _, change := range obj.ChangeSet {
			fmt.Println("Name : ", change.Name)
			fmt.Println("Operation Type : ", change.Op)
			fmt.Printf("Change Val : %+v\n", change.Val)
			fmt.Println("*********************************")
		}
		fmt.Println("-----------Object data Ended---------------")
	}

	fmt.Println("###############################################")
	fmt.Println()
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

const (
	apiProtocol = "https://"
)

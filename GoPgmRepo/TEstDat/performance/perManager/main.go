package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
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
	vmProp        = "vm"
	apiProtocol   = "https://"
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
	//var hostSysProps []string
	hostSysProps := []string{summaryProp, runtimeProp, configProp, hardwareProp, datastoreProp, vmProp}

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

func VirtualMachinePerformance(ctx context.Context, c *govmomi.Client, hs []mo.HostSystem) {
	var (
		entities           = make([]types.ManagedObjectReference, 0)
		metrics            = make([]string, 0)
		performanceManager = performance.NewManager(c.Client)
		spec               = types.PerfQuerySpec{MaxSample: 1, MetricId: []types.PerfMetricId{{Instance: "*"}}, Format: string(types.PerfFormatNormal)}
	)

	for _, host := range hs {
		for _, vm := range host.Vm {
			entities = append(entities, vm.Reference())
		}
	}

	metrics = append(metrics, listPerformanceMetrics()...)

	sn, err := performanceManager.SampleByName(ctx, spec, metrics, entities)
	if err != nil {
		return
	}

	result, err := performanceManager.ToMetricSeries(ctx, sn)

	if err != nil {
		return
	}

	// Read result
	for _, metric := range result {
		name := metric.Entity

		for _, v := range metric.Value {

			instance := v.Instance
			if instance == "" { //That means it is showing _Total
				instance = "_Total"
			}
			if len(v.Value) != 0 {
				fmt.Printf("%s\t%s\t%s\t%s\n",
					name, instance, v.Name, v.ValueCSV())
			}
		}
	}

}

func HostPerformance(ctx context.Context, c *govmomi.Client, hs []mo.HostSystem) {
	var (
		entities           = make([]types.ManagedObjectReference, 0)
		metrics            = make([]string, 0)
		performanceManager = performance.NewManager(c.Client)
		spec               = types.PerfQuerySpec{MaxSample: 1, MetricId: []types.PerfMetricId{{Instance: "*"}}}
	)

	for _, host := range hs {
		entities = append(entities, host.Reference())
	}

	metrics = append(metrics, listPerformanceMetrics()...)

	sn, err := performanceManager.SampleByName(ctx, spec, metrics, entities)
	if err != nil {
		return
	}

	result, err := performanceManager.ToMetricSeries(ctx, sn)

	if err != nil {
		return
	}

	// Read result
	for _, metric := range result {
		name := metric.Entity

		for _, v := range metric.Value {

			instance := v.Instance
			if instance == "" { //That means it is showing _Total
				instance = "_Total"
			}
			if len(v.Value) != 0 {
				fmt.Printf("%s\t%s\t%s\t%s\n",
					name, instance, v.Name, v.ValueCSV())
			}
		}
	}

}

func listPerformanceMetrics() []string {
	return []string{
		"cpu.usage.none",
		// "cpu.ready.summation",
		// "net.usage.average",
		// "net.transmitted.average",
		// "mem.usage.none",
		// "mem.active.average",
		// "mem.compressed.average",
		// "mem.consumed.average",
		// "mem.overhead.average",
		// "mem.shared.average",
		// "mem.swapin.average",
		// "mem.swapout.average",
		// "mem.swapused.average",
		// "mem.vmmemctl.average",
		// "mem.state.latest",
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

var interval = flag.Int("i", 20, "Interval ID")

func PerfRun(ctx context.Context, c *vim25.Client) error {

	var (
		entities = make([]types.ManagedObjectReference, 0)
	)

	var hostSystems []mo.HostSystem
	hostSystem := []string{hostSystemObj}
	//var hostSysProps []string
	hostSysProps := []string{summaryProp, runtimeProp, configProp, hardwareProp, datastoreProp, vmProp}

	var (
		metrics = make([]string, 0)
	)
	// Get virtual machines references
	m := view.NewManager(c)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, hostSystem, true)
	if err != nil {
		return err
	}

	err = v.Retrieve(ctx, hostSystem, hostSysProps, &hostSystems)

	if err != nil {
		return err
	}

	defer v.Destroy(ctx)

	for _, host := range hostSystems {
		entities = append(entities, host.Reference())
	}
	//vmsRefs, err := v.Find(ctx, []string{"HostSystem"}, nil)
	if err != nil {
		return err
	}

	// Create a PerfManager
	perfManager := performance.NewManager(c)

	// Retrieve counters name list
	counters, err := perfManager.CounterInfoByName(ctx)
	if err != nil {
		return err
	}

	var names []string
	for name := range counters {
		names = append(names, name)
	}

	// Create PerfQuerySpec
	spec := types.PerfQuerySpec{
		MaxSample:  1,
		MetricId:   []types.PerfMetricId{{Instance: "*"}},
		IntervalId: int32(*interval),
	}

	metrics = append(metrics, listPerformanceMetrics()...)

	// Query metrics
	sample, err := perfManager.SampleByName(ctx, spec, metrics, entities)
	if err != nil {
		return err
	}

	result, err := perfManager.ToMetricSeries(ctx, sample)
	if err != nil {
		return err
	}

	// Read result
	for _, metric := range result {
		name := metric.Entity

		for _, v := range metric.Value {
			counter := counters[v.Name]
			units := counter.UnitInfo.GetElementDescription().Label

			instance := v.Instance
			if instance == "" {
				instance = "_Total"
			}

			if len(v.Value) != 0 {
				fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
					name, instance, v.Name, v.ValueCSV(), units)
			}
		}
	}
	return nil
}

func main() {
	fmt.Println("Hello")
	hostObj := Host{}
	hostObj.Address = "10.4.129.10"
	hostObj.Username = "root"
	hostObj.Password = "Admin@123"

	// hostObj.Address = "10.4.129.11"
	// hostObj.Username = "administrator@vcenter.local"
	// hostObj.Password = "Pass@123"

	// Vcenter Host Info for qa env
	// hostObj.Address = "10.2.40.120"
	// hostObj.Username = "dtitsupport247\\agentteam"
	// hostObj.Password = "agent@123"

	ctx, _ := CancelHandler()

	fmt.Println("Connecting to VMWare VCenter")
	conn, err := NewClient(ctx, hostObj)
	if nil != err {
		fmt.Println(err)
		return
	}

	PerfRun(ctx, conn.Client)

	/*
		fmt.Println("Successfully connected to VMWare VCenter")
		hostSys, _ := HostSystems(ctx, conn)
		fmt.Println("Host Performance Counters")
		HostPerformance(ctx, conn, hostSys)
		fmt.Println("Virtual Machine Performance Counters")
	*/
	//VirtualMachinePerformance(ctx, conn, hostSys)
	//fmt.Println(hostSys)

}

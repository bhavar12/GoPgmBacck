package main

import (
	"errors"
	"fmt"
	"runtime"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

type abbr struct {
	std string
	dst string
}

//go:generate env ZONEINFO=$GOROOT/lib/time/zoneinfo.zip go run genzabbrs.go -output zoneinfo_abbrs_windows.go

// A Location maps time instants to the zone in use at that time.
// Typically, the Location represents the collection of time offsets
// in use in a geographical area, such as CEST and CET for central Europe.
type Location struct {
	name string
	zone []zone
	tx   []zoneTrans

	// Most lookups will be for the current time.
	// To avoid the binary search through tx, keep a
	// static one-element cache that gives the correct
	// zone for the time when the Location was created.
	// if cacheStart <= t < cacheEnd,
	// lookup can return cacheZone.
	// The units for cacheStart and cacheEnd are seconds
	// since January 1, 1970 UTC, to match the argument
	// to lookup.
	cacheStart int64
	cacheEnd   int64
	cacheZone  *zone
}

// A zone represents a single time zone such as CEST or CET.
type zone struct {
	name   string // abbreviated name, "CET"
	offset int    // seconds east of UTC
	isDST  bool   // is this zone Daylight Savings Time?
}

// A zoneTrans represents a single time zone transition.
type zoneTrans struct {
	when         int64 // transition time, in seconds since 1970 GMT
	index        uint8 // the index of the zone that goes into effect at that time
	isstd, isutc bool  // ignored - no idea what these mean
}

var localLoc Location
var zoneSources = []string{
	runtime.GOROOT() + "/lib/time/zoneinfo.zip",
}
var abbrs = map[string]abbr{
	"Egypt Standard Time":             {"EET", "EET"},     // Africa/Cairo
	"Morocco Standard Time":           {"WET", "WEST"},    // Africa/Casablanca
	"South Africa Standard Time":      {"SAST", "SAST"},   // Africa/Johannesburg
	"Sudan Standard Time":             {"CAT", "CAT"},     // Africa/Khartoum
	"W. Central Africa Standard Time": {"WAT", "WAT"},     // Africa/Lagos
	"E. Africa Standard Time":         {"EAT", "EAT"},     // Africa/Nairobi
	"Libya Standard Time":             {"EET", "EET"},     // Africa/Tripoli
	"Namibia Standard Time":           {"CAT", "CAT"},     // Africa/Windhoek
	"Aleutian Standard Time":          {"HST", "HDT"},     // America/Adak
	"Alaskan Standard Time":           {"AKST", "AKDT"},   // America/Anchorage
	"Tocantins Standard Time":         {"-03", "-03"},     // America/Araguaina
	"Paraguay Standard Time":          {"-04", "-03"},     // America/Asuncion
	"Bahia Standard Time":             {"-03", "-03"},     // America/Bahia
	"SA Pacific Standard Time":        {"-05", "-05"},     // America/Bogota
	"Argentina Standard Time":         {"-03", "-03"},     // America/Buenos_Aires
	"Eastern Standard Time (Mexico)":  {"EST", "EST"},     // America/Cancun
	"Venezuela Standard Time":         {"-04", "-04"},     // America/Caracas
	"SA Eastern Standard Time":        {"-03", "-03"},     // America/Cayenne
	"Central Standard Time":           {"CST", "CDT"},     // America/Chicago
	"Mountain Standard Time (Mexico)": {"MST", "MDT"},     // America/Chihuahua
	"Central Brazilian Standard Time": {"-04", "-03"},     // America/Cuiaba
	"Mountain Standard Time":          {"MST", "MDT"},     // America/Denver
	"Greenland Standard Time":         {"-03", "-02"},     // America/Godthab
	"Turks And Caicos Standard Time":  {"AST", "EDT"},     // America/Grand_Turk
	"Central America Standard Time":   {"CST", "CST"},     // America/Guatemala
	"Atlantic Standard Time":          {"AST", "ADT"},     // America/Halifax
	"Cuba Standard Time":              {"CST", "CDT"},     // America/Havana
	"US Eastern Standard Time":        {"EST", "EDT"},     // America/Indianapolis
	"SA Western Standard Time":        {"-04", "-04"},     // America/La_Paz
	"Pacific Standard Time":           {"PST", "PDT"},     // America/Los_Angeles
	"Central Standard Time (Mexico)":  {"CST", "CDT"},     // America/Mexico_City
	"Saint Pierre Standard Time":      {"-03", "-02"},     // America/Miquelon
	"Montevideo Standard Time":        {"-03", "-03"},     // America/Montevideo
	"Eastern Standard Time":           {"EST", "EDT"},     // America/New_York
	"US Mountain Standard Time":       {"MST", "MST"},     // America/Phoenix
	"Haiti Standard Time":             {"EST", "EDT"},     // America/Port-au-Prince
	"Magallanes Standard Time":        {"-03", "-03"},     // America/Punta_Arenas
	"Canada Central Standard Time":    {"CST", "CST"},     // America/Regina
	"Pacific SA Standard Time":        {"-04", "-03"},     // America/Santiago
	"E. South America Standard Time":  {"-03", "-02"},     // America/Sao_Paulo
	"Newfoundland Standard Time":      {"NST", "NDT"},     // America/St_Johns
	"Pacific Standard Time (Mexico)":  {"PST", "PDT"},     // America/Tijuana
	"Central Asia Standard Time":      {"+06", "+06"},     // Asia/Almaty
	"Jordan Standard Time":            {"EET", "EEST"},    // Asia/Amman
	"Arabic Standard Time":            {"+03", "+03"},     // Asia/Baghdad
	"Azerbaijan Standard Time":        {"+04", "+04"},     // Asia/Baku
	"SE Asia Standard Time":           {"+07", "+07"},     // Asia/Bangkok
	"Altai Standard Time":             {"+07", "+07"},     // Asia/Barnaul
	"Middle East Standard Time":       {"EET", "EEST"},    // Asia/Beirut
	"India Standard Time":             {"IST", "IST"},     // Asia/Calcutta
	"Transbaikal Standard Time":       {"+09", "+09"},     // Asia/Chita
	"Sri Lanka Standard Time":         {"+0530", "+0530"}, // Asia/Colombo
	"Syria Standard Time":             {"EET", "EEST"},    // Asia/Damascus
	"Bangladesh Standard Time":        {"+06", "+06"},     // Asia/Dhaka
	"Arabian Standard Time":           {"+04", "+04"},     // Asia/Dubai
	"West Bank Standard Time":         {"EET", "EEST"},    // Asia/Hebron
	"W. Mongolia Standard Time":       {"+07", "+07"},     // Asia/Hovd
	"North Asia East Standard Time":   {"+08", "+08"},     // Asia/Irkutsk
	"Israel Standard Time":            {"IST", "IDT"},     // Asia/Jerusalem
	"Afghanistan Standard Time":       {"+0430", "+0430"}, // Asia/Kabul
	"Russia Time Zone 11":             {"+12", "+12"},     // Asia/Kamchatka
	"Pakistan Standard Time":          {"PKT", "PKT"},     // Asia/Karachi
	"Nepal Standard Time":             {"+0545", "+0545"}, // Asia/Katmandu
	"North Asia Standard Time":        {"+07", "+07"},     // Asia/Krasnoyarsk
	"Magadan Standard Time":           {"+11", "+11"},     // Asia/Magadan
	"N. Central Asia Standard Time":   {"+07", "+07"},     // Asia/Novosibirsk
	"Omsk Standard Time":              {"+06", "+06"},     // Asia/Omsk
	"North Korea Standard Time":       {"KST", "KST"},     // Asia/Pyongyang
	"Myanmar Standard Time":           {"+0630", "+0630"}, // Asia/Rangoon
	"Arab Standard Time":              {"+03", "+03"},     // Asia/Riyadh
	"Sakhalin Standard Time":          {"+11", "+11"},     // Asia/Sakhalin
	"Korea Standard Time":             {"KST", "KST"},     // Asia/Seoul
	"China Standard Time":             {"CST", "CST"},     // Asia/Shanghai
	"Singapore Standard Time":         {"+08", "+08"},     // Asia/Singapore
	"Russia Time Zone 10":             {"+11", "+11"},     // Asia/Srednekolymsk
	"Taipei Standard Time":            {"CST", "CST"},     // Asia/Taipei
	"West Asia Standard Time":         {"+05", "+05"},     // Asia/Tashkent
	"Georgian Standard Time":          {"+04", "+04"},     // Asia/Tbilisi
	"Iran Standard Time":              {"+0330", "+0430"}, // Asia/Tehran
	"Tokyo Standard Time":             {"JST", "JST"},     // Asia/Tokyo
	"Tomsk Standard Time":             {"+07", "+07"},     // Asia/Tomsk
	"Ulaanbaatar Standard Time":       {"+08", "+08"},     // Asia/Ulaanbaatar
	"Vladivostok Standard Time":       {"+10", "+10"},     // Asia/Vladivostok
	"Yakutsk Standard Time":           {"+09", "+09"},     // Asia/Yakutsk
	"Ekaterinburg Standard Time":      {"+05", "+05"},     // Asia/Yekaterinburg
	"Caucasus Standard Time":          {"+04", "+04"},     // Asia/Yerevan
	"Azores Standard Time":            {"-01", "+00"},     // Atlantic/Azores
	"Cape Verde Standard Time":        {"-01", "-01"},     // Atlantic/Cape_Verde
	"Greenwich Standard Time":         {"GMT", "GMT"},     // Atlantic/Reykjavik
	"Cen. Australia Standard Time":    {"ACST", "ACDT"},   // Australia/Adelaide
	"E. Australia Standard Time":      {"AEST", "AEST"},   // Australia/Brisbane
	"AUS Central Standard Time":       {"ACST", "ACST"},   // Australia/Darwin
	"Aus Central W. Standard Time":    {"+0845", "+0845"}, // Australia/Eucla
	"Tasmania Standard Time":          {"AEST", "AEDT"},   // Australia/Hobart
	"Lord Howe Standard Time":         {"+1030", "+11"},   // Australia/Lord_Howe
	"W. Australia Standard Time":      {"AWST", "AWST"},   // Australia/Perth
	"AUS Eastern Standard Time":       {"AEST", "AEDT"},   // Australia/Sydney
	"UTC":                             {"GMT", "GMT"},     // Etc/GMT
	"UTC-11":                          {"-11", "-11"},     // Etc/GMT+11
	"Dateline Standard Time":          {"-12", "-12"},     // Etc/GMT+12
	"UTC-02":                          {"-02", "-02"},     // Etc/GMT+2
	"UTC-08":                          {"-08", "-08"},     // Etc/GMT+8
	"UTC-09":                          {"-09", "-09"},     // Etc/GMT+9
	"UTC+12":                          {"+12", "+12"},     // Etc/GMT-12
	"UTC+13":                          {"+13", "+13"},     // Etc/GMT-13
	"Astrakhan Standard Time":         {"+04", "+04"},     // Europe/Astrakhan
	"W. Europe Standard Time":         {"CET", "CEST"},    // Europe/Berlin
	"GTB Standard Time":               {"EET", "EEST"},    // Europe/Bucharest
	"Central Europe Standard Time":    {"CET", "CEST"},    // Europe/Budapest
	"E. Europe Standard Time":         {"EET", "EEST"},    // Europe/Chisinau
	"Turkey Standard Time":            {"+03", "+03"},     // Europe/Istanbul
	"Kaliningrad Standard Time":       {"EET", "EET"},     // Europe/Kaliningrad
	"FLE Standard Time":               {"EET", "EEST"},    // Europe/Kiev
	"GMT Standard Time":               {"GMT", "BST"},     // Europe/London
	"Belarus Standard Time":           {"+03", "+03"},     // Europe/Minsk
	"Russian Standard Time":           {"MSK", "MSK"},     // Europe/Moscow
	"Romance Standard Time":           {"CET", "CEST"},    // Europe/Paris
	"Russia Time Zone 3":              {"+04", "+04"},     // Europe/Samara
	"Saratov Standard Time":           {"+04", "+04"},     // Europe/Saratov
	"Central European Standard Time":  {"CET", "CEST"},    // Europe/Warsaw
	"Mauritius Standard Time":         {"+04", "+04"},     // Indian/Mauritius
	"Samoa Standard Time":             {"+13", "+14"},     // Pacific/Apia
	"New Zealand Standard Time":       {"NZST", "NZDT"},   // Pacific/Auckland
	"Bougainville Standard Time":      {"+11", "+11"},     // Pacific/Bougainville
	"Chatham Islands Standard Time":   {"+1245", "+1345"}, // Pacific/Chatham
	"Easter Island Standard Time":     {"-06", "-05"},     // Pacific/Easter
	"Fiji Standard Time":              {"+12", "+13"},     // Pacific/Fiji
	"Central Pacific Standard Time":   {"+11", "+11"},     // Pacific/Guadalcanal
	"Hawaiian Standard Time":          {"HST", "HST"},     // Pacific/Honolulu
	"Line Islands Standard Time":      {"+14", "+14"},     // Pacific/Kiritimati
	"Marquesas Standard Time":         {"-0930", "-0930"}, // Pacific/Marquesas
	"Norfolk Standard Time":           {"+11", "+11"},     // Pacific/Norfolk
	"West Pacific Standard Time":      {"+10", "+10"},     // Pacific/Port_Moresby
	"Tonga Standard Time":             {"+13", "+13"},     // Pacific/Tongatapu
}

func initLocal() {
	var i syscall.Timezoneinformation
	if _, err := syscall.GetTimeZoneInformation(&i); err != nil {
		localLoc.name = "UTC"
		return
	}
	initLocalFromTZI(&i)
}
func initLocalFromTZI(i *syscall.Timezoneinformation) {
	l := &localLoc
	l.name = "Local"
	stdname, dstname := abbrev(i)
	fmt.Println("from api.. std name ", syscall.UTF16ToString(i.StandardName[:]))
	fmt.Println("from api day light name ", syscall.UTF16ToString(i.DaylightName[:]))
	fmt.Println("stdName", stdname)
	fmt.Println("dstName", dstname)
}

// abbrev returns the abbreviations to use for the given zone z.
func abbrev(z *syscall.Timezoneinformation) (std, dst string) {
	fmt.Println("I am abbrev call")
	stdName := syscall.UTF16ToString(z.StandardName[:])
	a, ok := abbrs[stdName]
	if !ok {
		dstName := syscall.UTF16ToString(z.DaylightName[:])
		// Perhaps stdName is not English. Try to convert it.
		englishName, err := toEnglishName(stdName, dstName)
		if err == nil {
			a, ok = abbrs[englishName]
			if ok {
				fmt.Println("English Name............", englishName)
				return a.std, a.dst
			}
		}
		fmt.Println("I am in before fallback")
		// fallback to using capital letters
		return extractCAPS(stdName), extractCAPS(dstName)
	}
	return a.std, a.dst
}

// extractCAPS extracts capital letters from description desc.
func extractCAPS(desc string) string {
	var short []rune
	for _, c := range desc {
		if 'A' <= c && c <= 'Z' {
			short = append(short, c)
		}
	}
	return string(short)
}

// toEnglishName searches the registry for an English name of a time zone
// whose zone names are stdname and dstname and returns the English name.
func toEnglishName(stdname, dstname string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Time Zones`, registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	names, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return "", err
	}
	for _, name := range names {
		matched, err := matchZoneKey(k, name, stdname, dstname)
		if err == nil && matched {
			return name, nil
		} else {
			fmt.Println("name ", name, " std", stdname, " dst", dstname)
		}
	}
	return "", errors.New(`English name for time zone "` + stdname + `" not found in registry`)
}

// TODO(rsc): Fall back to copy of zoneinfo files.

// BUG(brainman,rsc): On Windows, the operating system does not provide complete
// time zone information.
// The implementation assumes that this year's rules for daylight savings
// time apply to all previous and future years as well.

// matchZoneKey checks if stdname and dstname match the corresponding key
// values "MUI_Std" and MUI_Dlt" or "Std" and "Dlt" (the latter down-level
// from Vista) in the kname key stored under the open registry key zones.
func matchZoneKey(zones registry.Key, kname string, stdname, dstname string) (matched bool, err2 error) {
	k, err := registry.OpenKey(zones, kname, registry.READ)
	if err != nil {
		return false, err
	}
	defer k.Close()

	var std, dlt string
	if err = registry.LoadRegLoadMUIString(); err == nil {
		// Try MUI_Std and MUI_Dlt first, fallback to Std and Dlt if *any* error occurs
		std, err = k.GetMUIStringValue("MUI_Std")
		if err == nil {
			dlt, err = k.GetMUIStringValue("MUI_Dlt")
		}
	}
	if err != nil { // Fallback to Std and Dlt
		if std, _, err = k.GetStringValue("Std"); err != nil {
			return false, err
		}
		if dlt, _, err = k.GetStringValue("Dlt"); err != nil {
			return false, err
		}
	}

	if std != stdname {
		return false, nil
	}
	if dlt != dstname && dstname != stdname {
		return false, nil
	}
	return true, nil
}
func main() {
	fmt.Println("In a main function")
	initLocal()
	fmt.Println("Main fucntion execution completed")
}

// +build windows

package main

import (
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	dlliphlpapi              = windows.NewLazySystemDLL("iphlpapi.dll")
	procGetIpForwardTable    = dlliphlpapi.NewProc("GetIpForwardTable")
	procCreateIpForwardEntry = dlliphlpapi.NewProc("CreateIpForwardEntry")
	procDeleteIpForwardEntry = dlliphlpapi.NewProc("DeleteIpForwardEntry")
)

type routestruct struct {
	ForwardDest      [4]byte
	ForwardMask      [4]byte
	ForwardPolicy    uint32
	ForwardNextHop   [4]byte
	ForwardIfIndex   uint32
	ForwardType      uint32
	ForwardProto     uint32
	ForwardAge       uint32
	ForwardNextHopAS uint32
	ForwardMetric1   uint32
	ForwardMetric2   uint32
	ForwardMetric3   uint32
	ForwardMetric4   uint32
	ForwardMetric5   uint32
}

func GetRoutes() {
	pIpForwardTable := newMem(uint32((4 + unsafe.Sizeof(routestruct{}))))
	pwdSize := len(pIpForwardTable)
	log.Println(" return value is ", pwdSize)

	a, _, _ := syscall.Syscall(uintptr(procGetIpForwardTable.Addr()), 3, uintptr(unsafe.Pointer(&pIpForwardTable[0])), uintptr(unsafe.Pointer(&pwdSize)), 0)

	for ii := 0; ii < 8; ii++ {
		log.Println(" Loop = ", ii)
		pIpForwardTable = make([]byte, pwdSize)
		a, _, _ := syscall.Syscall(uintptr(procGetIpForwardTable.Addr()), 3, uintptr(unsafe.Pointer(&pIpForwardTable[0])), uintptr(unsafe.Pointer(&pwdSize)), 0)
		if a == 122 {
			log.Println("got 122")
			continue
		}
		break
	}
	if a != 0 {
		log.Println("No more routes")
	}

	num := *(*uint32)(unsafe.Pointer(&pIpForwardTable[0]))
	log.Println("Number of routes ", num)
	routes := make([]routestruct, num)
	sr := uintptr(unsafe.Pointer(&pIpForwardTable[0])) + unsafe.Sizeof(num)
	rowsize := unsafe.Sizeof(routestruct{})

	if len(pIpForwardTable) < int(unsafe.Sizeof(num)+rowsize*uintptr(num)) {
		log.Println(" returned number is too long")
	}

	for i := uint32(0); i < num; i++ {
		routes[i] = *((*routestruct)(unsafe.Pointer(sr + (rowsize * uintptr(i)))))
		log.Println(i+1, ". route is ", routes[i])
	}
}

func Create() {
	dummyroute := routestruct{
		ForwardDest: [4]byte{2, 1, 192, 0}, ForwardMask: [4]byte{255, 255, 252, 0}, ForwardPolicy: 0, ForwardNextHop: [4]byte{172, 19, 217, 36}, ForwardIfIndex: 10, ForwardType: 3, ForwardProto: 3, ForwardAge: 1846, ForwardNextHopAS: 0, ForwardMetric1: 8, ForwardMetric2: 0, ForwardMetric3: 0, ForwardMetric4: 0, ForwardMetric5: 0}

	if CreateRoute(&dummyroute) {
		log.Println(" Added the route")
	} else {
		log.Println("Could not add the route")
	}
}

func CreateRoute(routeEntry *routestruct) bool {

	log.Println("Route to add is ", routeEntry)
	a, _, er := syscall.Syscall(uintptr(procCreateIpForwardEntry.Addr()), 1, uintptr(unsafe.Pointer(routeEntry)), 0, 0)
	if a == 5010 {
		log.Println("Entry Exists")
		return false
	} else if a != 0 {
		log.Println(" Could not add it ", er, " a=", a)
		return false
	}
	return true
}

func DeleteRoute() bool {
	routeEntry := &routestruct{
		ForwardDest: [4]byte{2, 1, 192, 0}, ForwardMask: [4]byte{255, 255, 252, 0}, ForwardPolicy: 0, ForwardNextHop: [4]byte{172, 19, 217, 36}, ForwardIfIndex: 10, ForwardType: 3, ForwardProto: 3, ForwardAge: 1846, ForwardNextHopAS: 0, ForwardMetric1: 8, ForwardMetric2: 0, ForwardMetric3: 0, ForwardMetric4: 0, ForwardMetric5: 0}

	log.Println("Route to Delete is ", routeEntry)
	a, _, er := syscall.Syscall(uintptr(procDeleteIpForwardEntry.Addr()), 1, uintptr(unsafe.Pointer(routeEntry)), 0, 0)
	if a != 0 {
		log.Println(" Could not delete it ", er, " a=", a)
		return false
	}
	return true

}

func newMem(size uint32) []byte {
	Memory := make([]byte, size, size)
	return Memory

}

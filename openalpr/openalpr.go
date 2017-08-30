package openalpr

// #cgo LDFLAGS: -L/usr/lib -lopenalpr -lopenalprgo
// #cgo CFLAGS: -I/usr/include/
/*
#include "openalprgo.h"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"unsafe"
	"github.com/dgmann/ma-shared"
)

type Alpr struct {
	//Country    string
	//configFile string
	//runtimeDir string

	cAlpr C.Alpr
}

func bool2Cint(b bool) C.int {
	if b {
		return 1
	} else {
		return 0
	}
}

func cint2Bool(i C.int) bool {
	if i == 0 {
		return false
	} else {
		return true
	}
}

func NewAlpr(country string, configFile string, runtimeDir string) *Alpr {
	cstrCountry := C.CString(country)
	cstrConfigFile := C.CString(configFile)
	cstrRuntimeDir := C.CString(runtimeDir)
	defer C.free(unsafe.Pointer(cstrCountry))
	defer C.free(unsafe.Pointer(cstrConfigFile))
	defer C.free(unsafe.Pointer(cstrRuntimeDir))

	alpr := C.AlprInit(cstrCountry, cstrConfigFile, cstrRuntimeDir)
	return &Alpr{cAlpr: alpr}
}

func (alpr *Alpr) SetDetectRegion(detectRegion bool) {
	C.SetDetectRegion(alpr.cAlpr, bool2Cint(detectRegion))
}

func (alpr *Alpr) SetTopN(topN int) {
	C.SetTopN(alpr.cAlpr, C.int(topN))
}

func (alpr *Alpr) SetDefaultRegion(region string) {
	cstrRegion := C.CString(region)
	defer C.free(unsafe.Pointer(cstrRegion))
	C.SetDefaultRegion(alpr.cAlpr, cstrRegion)
}

func (alpr *Alpr) IsLoaded() bool {
	return cint2Bool(C.IsLoaded(alpr.cAlpr))
}

func GetVersion() string {
	return C.GoString(C.GetVersion())
}

func (alpr *Alpr) RecognizeByFilePath(filePath string) (shared.OpenAlprResponse, error) {
	cstrFilePath := C.CString(filePath)
	defer C.free(unsafe.Pointer(cstrFilePath))
	stringResult := C.GoString(C.RecognizeByFilePath(alpr.cAlpr, cstrFilePath))
	fmt.Println(stringResult)

	var results shared.OpenAlprResponse
	err := json.Unmarshal([]byte(stringResult), &results)

	return results, err
}

func (alpr *Alpr) RecognizeByBlob(imageBytes []byte) (shared.OpenAlprResponse, error) {
	stringImageBytes := string(imageBytes)
	cstrImageBytes := C.CString(stringImageBytes)
	defer C.free(unsafe.Pointer(cstrImageBytes))
	stringResult := C.GoString(C.RecognizeByBlob(alpr.cAlpr, cstrImageBytes, C.int(len(imageBytes))))

	var results shared.OpenAlprResponse
	err := json.Unmarshal([]byte(stringResult), &results)

	return results, err
}

func (alpr *Alpr) Unload() {
	C.Unload(alpr.cAlpr)
}
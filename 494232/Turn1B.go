package main

import (
	"flag"
	"fmt"
	"net/url"
)

// Define command-line flags
var (
	urlFlag      = flag.String("url", "", "URL to test")
	versionFlag = flag.String("version", "v1.0.0", "API version to use")
)

// Define a data structure to represent API parameters
type APIv1 struct {
	Param1 string
	Param2 int
}

type APIv2 struct {
	APIv1
	Param3 bool
}

// Implement logic to parse parameters based on version
func parseParams(u *url.URL, version string) interface{} {
	params := u.Query()

	switch version {
	case "v1.0.0":
		var apiV1 APIv1
		apiV1.Param1 = params.Get("param1")
		apiV1.Param2 = parseInt(params.Get("param2"))
		return apiV1
	case "v2.0.0":
		var apiV2 APIv2
		apiV2.Param1 = params.Get("param1")
		apiV2.Param2 = parseInt(params.Get("param2"))
		apiV2.Param3 = parseBool(params.Get("param3"))
		return apiV2
	default:
		return nil
	}
}

// Implement comparison function
func compareParams(params1, params2 interface{}) bool {
	switch params1 := params1.(type) {
	case APIv1:
		params2 := params2.(APIv1)
		return params1.Param1 == params2.Param1 && params1.Param2 == params2.Param2
	case APIv2:
		params2 := params2.(APIv2)
		return params1.Param1 == params2.Param1 && params1.Param2 == params2.Param2 && params1.Param3 == params2.Param3
	default:
		return false
	}
}

// Main function
func main() {
	// Parse command-line flags
	flag.Parse()

	// Parse URL
	u, err := url.Parse(*urlFlag)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Parse parameters based on version
	params1 := parseParams(u, *versionFlag)

	// Compare parameters with older versions (e.g., v1.0.0)
	params2 := parseParams(u, "v1.0.0")
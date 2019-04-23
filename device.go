package onvif

import (
	"fmt"
	"strings"
)

var deviceXMLNs = []string{
	`xmlns:tds="http://www.onvif.org/ver10/device/wsdl"`,
	`xmlns:tt="http://www.onvif.org/ver10/schema"`,
	`xmlns:tptz="http://www.onvif.org/ver20/ptz/wsdl"`,
}

// GetInformation fetch information of ONVIF camera
func (device Device) GetInformation() (DeviceInformation, error) {
	// Create SOAP
	soap := SOAP{
		Body:  "<tds:GetDeviceInformation/>",
		XMLNs: deviceXMLNs,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return DeviceInformation{}, err
	}

	// Parse response to interface
	deviceInfo, err := response.ValueForPath("Envelope.Body.GetDeviceInformationResponse")
	if err != nil {
		return DeviceInformation{}, err
	}

	// Parse interface to struct
	result := DeviceInformation{}
	if mapInfo, ok := deviceInfo.(map[string]interface{}); ok {
		result.Manufacturer = interfaceToString(mapInfo["Manufacturer"])
		result.Model = interfaceToString(mapInfo["Model"])
		result.FirmwareVersion = interfaceToString(mapInfo["FirmwareVersion"])
		result.SerialNumber = interfaceToString(mapInfo["SerialNumber"])
		result.HardwareID = interfaceToString(mapInfo["HardwareId"])
	}

	return result, nil
}

// GetCapabilities fetch info of ONVIF camera's capabilities
func (device Device) GetCapabilities() (DeviceCapabilities, error) {
	// Create SOAP
	soap := SOAP{
		XMLNs: deviceXMLNs,
		Body: `<tds:GetCapabilities>
			<tds:Category>All</tds:Category>
		</tds:GetCapabilities>`,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return DeviceCapabilities{}, err
	}

	// Get network capabilities
	envelopeBodyPath := "Envelope.Body.GetCapabilitiesResponse.Capabilities"
	ifaceNetCap, err := response.ValueForPath(envelopeBodyPath + ".Device.Network")
	if err != nil {
		return DeviceCapabilities{}, err
	}

	netCap := NetworkCapabilities{}
	if mapNetCap, ok := ifaceNetCap.(map[string]interface{}); ok {
		netCap.DynDNS = interfaceToBool(mapNetCap["DynDNS"])
		netCap.IPFilter = interfaceToBool(mapNetCap["IPFilter"])
		netCap.IPVersion6 = interfaceToBool(mapNetCap["IPVersion6"])
		netCap.ZeroConfig = interfaceToBool(mapNetCap["ZeroConfiguration"])
	}

	// Get events capabilities
	ifaceEventsCap, err := response.ValueForPath(envelopeBodyPath + ".Events")
	if err != nil {
		return DeviceCapabilities{}, err
	}

	eventsCap := make(map[string]bool)
	if mapEventsCap, ok := ifaceEventsCap.(map[string]interface{}); ok {
		for key, value := range mapEventsCap {
			if strings.ToLower(key) == "xaddr" {
				continue
			}

			key = strings.Replace(key, "WS", "", 1)
			eventsCap[key] = interfaceToBool(value)
		}
	}

	// Get streaming capabilities
	ifaceStreamingCap, err := response.ValueForPath(envelopeBodyPath + ".Media.StreamingCapabilities")
	if err != nil {
		return DeviceCapabilities{}, err
	}

	streamingCap := make(map[string]bool)
	if mapStreamingCap, ok := ifaceStreamingCap.(map[string]interface{}); ok {
		for key, value := range mapStreamingCap {
			key = strings.Replace(key, "_", " ", -1)
			streamingCap[key] = interfaceToBool(value)
		}
	}

	// Create final result
	deviceCapabilities := DeviceCapabilities{
		Network:   netCap,
		Events:    eventsCap,
		Streaming: streamingCap,
	}

	return deviceCapabilities, nil
}

// GetDiscoveryMode fetch network discovery mode of an ONVIF camera
func (device Device) GetDiscoveryMode() (string, error) {
	// Create SOAP
	soap := SOAP{
		Body:  "<tds:GetDiscoveryMode/>",
		XMLNs: deviceXMLNs,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return "", err
	}

	// Parse response
	discoveryMode, _ := response.ValueForPathString("Envelope.Body.GetDiscoveryModeResponse.DiscoveryMode")
	return discoveryMode, nil
}

// GetScopes fetch scopes of an ONVIF camera
func (device Device) GetScopes() ([]string, error) {
	// Create SOAP
	soap := SOAP{
		Body:  "<tds:GetScopes/>",
		XMLNs: deviceXMLNs,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return nil, err
	}

	// Parse response to interface
	ifaceScopes, err := response.ValuesForPath("Envelope.Body.GetScopesResponse.Scopes")
	if err != nil {
		return nil, err
	}

	// Convert interface to array of scope
	scopes := []string{}
	for _, ifaceScope := range ifaceScopes {
		if mapScope, ok := ifaceScope.(map[string]interface{}); ok {
			scope := interfaceToString(mapScope["ScopeItem"])
			scopes = append(scopes, scope)
		}
	}

	return scopes, nil
}

// GetHostname fetch hostname of an ONVIF camera
func (device Device) Ptz(Token, x, y, z string) error {
	// Create SOAP
	soap := SOAP{
		Body: `<tptz:ContinuousMove>
    <tptz:ProfileToken>` + Token + `</tptz:ProfileToken>
    <tptz:Velocity>
     <tt:PanTilt x="` + x + `" y="` + y + `" space="">
     </tt:PanTilt>
     <tt:Zoom x="` + z + `" space="">
     </tt:Zoom>
    </tptz:Velocity>
   </tptz:ContinuousMove>`,
		XMLNs: deviceXMLNs,
	}

	// Send SOAP request
	_, err := soap.SendRequest(device.XAddr)
	return err
}
func (device Device) PtzStop(Token, x, y, z string) error {
	// Create SOAP
	soap := SOAP{
		Body: `<tptz:Stop>
    <tptz:ProfileToken>` + Token + `</tptz:ProfileToken>
		 <tptz:PanTilt>false</tptz:PanTilt>
		 <tptz:Zoom>false</tptz:Zoom>
   </tptz:Stop>`,
		XMLNs: deviceXMLNs,
	}

	// Send SOAP request
	_, err := soap.SendRequest(device.XAddr)
	return err
}

// GetHostname fetch hostname of an ONVIF camera
func (device Device) GetHostname() (HostnameInformation, error) {
	// Create SOAP
	soap := SOAP{
		Body:  "<tds:GetHostname/>",
		XMLNs: deviceXMLNs,
	}

	// Send SOAP request
	response, err := soap.SendRequest(device.XAddr)
	if err != nil {
		return HostnameInformation{}, err
	}

	// Parse response to interface
	ifaceHostInfo, err := response.ValueForPath("Envelope.Body.GetHostnameResponse.HostnameInformation")
	if err != nil {
		return HostnameInformation{}, err
	}

	// Parse interface to struct
	hostnameInfo := HostnameInformation{}
	if mapHostInfo, ok := ifaceHostInfo.(map[string]interface{}); ok {
		hostnameInfo.Name = interfaceToString(mapHostInfo["Name"])
		hostnameInfo.FromDHCP = interfaceToBool(mapHostInfo["FromDHCP"])
	}

	return hostnameInfo, nil
}

// AppPTZMove move
func AppPTZMove(action string) {
	ip := "171.25.232.42"
	port := "11999"
	login := "admin"
	password := "Ghjlern14"

	var testDevice = Device{
		User:     login,
		Password: password,
		XAddr:    "http://" + login + ":" + password + "@" + ip + ":" + port + "/onvif/device_service",
		// XAddr: "http://" + login + ":" + password + "@" + ip + ":" + port + "/onvif/media_service",
	}
	res, err := testDevice.GetProfiles()
	if err != nil && err.Error() == "Unknown Action" {
		testDevice.XAddr = "http://" + login + ":" + password + "@" + ip + ":" + port + "/onvif/media_service"
		res, err = testDevice.GetProfiles()
		if err == nil {
			testDevice.XAddr = "http://" + login + ":" + password + "@" + ip + ":" + port + "/onvif/ptz_service"
		}
	}
	if err == nil && len(res) > 0 {
		switch action {
		case "up":
			err := testDevice.Ptz(res[0].Token, "0.0", "0.1", "0.0")
			if err != nil {
				err = testDevice.Ptz(res[0].PTZConfig.Token, "0.0", "0.1", "0.0")
				if err != nil {
					fmt.Println(err)
					// WriteFormatGIN(0, 200, err.Error(), c)
					return
				}
			}
		case "down":
			err := testDevice.Ptz(res[0].Token, "0.0", "-0.1", "0.0")
			if err != nil {
				err = testDevice.Ptz(res[0].PTZConfig.Token, "0.0", "-0.1", "0.0")
				if err != nil {
					fmt.Println(err)
					// WriteFormatGIN(0, 200, err.Error(), c)
					return
				}
			}
		case "left":
			err := testDevice.Ptz(res[0].Token, "-0.1", "0.0", "0.0")
			if err != nil {
				err = testDevice.Ptz(res[0].PTZConfig.Token, "-0.1", "0.0", "0.0")
				if err != nil {
					fmt.Println(err)
					// WriteFormatGIN(0, 200, err.Error(), c)
					return
				}
			}
		case "right":
			err := testDevice.Ptz(res[0].Token, "0.1", "0.0", "0.0")
			if err != nil {
				err = testDevice.Ptz(res[0].PTZConfig.Token, "0.1", "0.0", "0.0")
				if err != nil {
					fmt.Println(err)
					// WriteFormatGIN(0, 200, err.Error(), c)
					return
				}
			}
		case "zoomin":
			err := testDevice.Ptz(res[0].Token, "0.0", "0.0", "0.1")
			if err != nil {
				err = testDevice.Ptz(res[0].PTZConfig.Token, "0.0", "0.0", "0.1")
				if err != nil {
					fmt.Println(err)
					// WriteFormatGIN(0, 200, err.Error(), c)
					return
				}
			}
		case "zoomout":
			err := testDevice.Ptz(res[0].Token, "0.0", "0.0", "-0.1")
			if err != nil {
				err = testDevice.Ptz(res[0].PTZConfig.Token, "0.0", "0.0", "-0.1")
				if err != nil {
					fmt.Println(err)
					// WriteFormatGIN(0, 200, err.Error(), c)
					return
				}
			}
		case "stop":
			testDevice.Ptz(res[0].Token, "0", "0", "0")
			err := testDevice.PtzStop(res[0].Token, "0", "0", "0")
			if err != nil {
				testDevice.Ptz(res[0].PTZConfig.Token, "0", "0", "0")
				err = testDevice.PtzStop(res[0].PTZConfig.Token, "0", "0", "0")
				if err != nil {
					fmt.Println(err)
					// WriteFormatGIN(0, 200, err.Error(), c)
					return
				}
			}
		}
		fmt.Println("SUCCESS")
		// WriteFormatGIN(1, 200, "success", c)
	} else {
		fmt.Println(err)
		// WriteFormatGIN(0, 200, err.Error(), c)
	}
}

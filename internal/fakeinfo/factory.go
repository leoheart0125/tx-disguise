package fakeinfo

func FackInfoFactory(serviceType string) IService {
	switch serviceType {
	case "top":
		return NewTopFakeInfoService("top")
	case "pure":
		return NewPureFakeSystemInfoService()
	default:
		return NewPureFakeSystemInfoService() // Default to pure fake system info
	}
}

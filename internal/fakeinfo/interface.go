package fakeinfo

type IService interface {
	GetFakeInfo() ([]string, error)
	SetProcessCount(count int)
}

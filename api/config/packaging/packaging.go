package packaging

type Info struct {
	CurrentVersion string
	CommitSHA      string
}

// NewPackagingInfo 创建带有默认值的 Info 实例的函数
func NewPackagingInfo() *Info {
	return &Info{
		CurrentVersion: "0.6.15",
		CommitSHA:      "",
	}
}

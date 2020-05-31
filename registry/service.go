package registry

//服务抽象
type Service struct {
	Name  string  `json:"name"`
	Nodes []*Node `json:"node"`
}

type Node struct {
	Id     string `json:"id"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Weight int    `json:"weight"`
}

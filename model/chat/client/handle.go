package client

type TestParams struct {
	Name string `json:"name,omitempty" validate:"required,gte=1,lte=50"` // 姓名
	Age  string `json:"age,omitempty" validate:"age,gte=1"`              // 年龄
}

type TestRly struct {
	Name    string `json:"name,omitempty"`    // 姓名
	Age     string `json:"age,omitempty"`     // 年龄
	ID      string `json:"id,omitempty"`      // ID
	Address string `json:"address,omitempty"` // 地址
}

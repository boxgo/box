package main

type (
	Common struct {
		AppName string `json:"name" desc:"Application name"`
		test    string `json:"test" desc:"111"`
	}
)

func (c *Common) Name() string {
	return ""
}

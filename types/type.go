
package types


type MockRequest struct {
    TypescriptInterface string `json:"interface" binding:"required"`
    Count              int    `json:"count" binding:"required,min=1,max=100"`
}

type Field struct {
	Index int
	Name       string
	Type       string
	IsOptional bool
}
package response

type BaseOutput struct {
	Status  int
	Code    string
	Message string
	Data    interface{}
	Count   int
}

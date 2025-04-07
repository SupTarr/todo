package my_context

type Context interface {
	Bind(any) error
	TodoID() string
	TransactionID() string
	Audience() string
	Status(int)
	JSON(int, any)
}

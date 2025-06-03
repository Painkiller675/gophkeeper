package interceptors

type key int

const (
	// ContextKeyUserID a key for userID adding into the context during the process of authentication
	ContextKeyUserID key = iota
)

package ports

type APIServer interface {
	ListenAndServe(addr string)
}

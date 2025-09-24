package openapi

type Getter[T any] interface {
	Get() T
}

type Setter[T any] interface {
	Set(T) error
}

type AccessTokenStorager interface {
	Getter[*AccessToken]
	Setter[*AccessToken]
}

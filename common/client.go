package common

type Client interface {
	FetchList() ([]byte, error)
	FetchItem(id string) ([]byte, error)
}

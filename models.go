package main

//HTTPinfo ...
type HTTPinfo struct {
	Version    string
	StatusCode string
	Code       int
}

//CURLData ...
type CURLData struct {
	HTTPInfo HTTPinfo
	Headers  map[string]string
	Body     []byte
}

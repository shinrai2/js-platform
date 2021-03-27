package http

type XmlHttpRequest struct {
	readyState int
	responseText string
	status int
	statusText string
}

func (xmlHttp *XmlHttpRequest)abort() {
	xmlHttp.readyState = 0
	// TODO
}

func (xmlHttp *XmlHttpRequest)getAllResponseHeaders() *string {
	if xmlHttp.readyState < 3 {
		return nil
	}
	s := "" // TODO
	return &s
}

func (xmlHttp *XmlHttpRequest)getResponseHeader(headerName string) *string {
	if xmlHttp.readyState < 3 { // TODO
		return nil
	}
	s := "" // TODO
	return &s
}
package pkg

type Request struct {
	Uuid     string `json:"uuid"`
	UserName string `json:"username"`
	Qoute    string `json:"qoute"`
}

type Response struct {
	Message string `json:"errmsg"`
	Uuuid   string `json:"uuid"`
}

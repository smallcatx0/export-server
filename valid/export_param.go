package valid

type ExportParam struct {
	Timestamp  int64      `json:"timestamp"`
	EXTType    string     `json:"ext_type"`
	Title      string     `json:"title"`
	UserID     string     `json:"user_id"`
	CallBack   string     `json:"call_back"`
	SourceType string     `json:"source_type"`
	SourceHTTP SourceHTTP `json:"source_http"`
	SourceSQL  SourceSQL  `json:"source_sql"`
	SourceRaw  SourceRaw  `json:"source_raw"`
}

type SourceHTTP struct {
	URL    string                 `json:"url"`
	Method string                 `json:"method"`
	Header map[string]string      `json:"header"`
	Param  map[string]interface{} `json:"param"`
}

type SourceRaw struct {
	List []interface{} `json:"list"`
}

type SourceSQL struct {
	SQL   string        `json:"sql"`
	Param []interface{} `json:"param"`
}

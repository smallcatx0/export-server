package valid

type SourceHTTP struct {
	URL    string                 `json:"url"`
	Method string                 `json:"method"`
	Header map[string]string      `json:"header"`
	Param  map[string]interface{} `json:"param"`
}

type SourceSQL struct {
	SQL   string        `json:"sql"`
	Param []interface{} `json:"param"`
}
type ExportParam struct {
	Timestamp  int64      `json:"timestamp"`
	EXTType    string     `json:"ext_type"`
	Title      string     `json:"title"`
	UserID     string     `json:"user_id"`
	CallBack   string     `json:"call_back"`
	SourceType string     `json:"source_type"`
	SourceHTTP SourceHTTP `json:"source_http"`
	SourceSQL  SourceSQL  `json:"source_sql"`
}

func (param *ExportParam) Valid() error {
	// 时间戳与当前时间差距不可大于5s
	// diff := carbon.Now().DiffInSecondsWithAbs(carbon.CreateFromTimestamp(param.Timestamp))
	// if diff >= 5 {
	// 	return exception.ParamInValid("请求时间与标准时间差不可大于5s")
	// }
	return nil
}

type ExpSHttpParam struct {
	Timestamp  int64      `json:"timestamp"`
	EXTType    string     `json:"ext_type"`
	Title      string     `json:"title"`
	UserID     string     `json:"user_id"`
	CallBack   string     `json:"call_back"`
	SourceHTTP SourceHTTP `json:"source_http"`
}

type ExpSRawParam struct {
	Timestamp int64  `json:"timestamp"`
	EXTType   string `json:"ext_type"`
	Title     string `json:"title"`
	UserID    string `json:"user_id"`
	CallBack  string `json:"call_back"`
	SourceRaw string `json:"source_raw"`
}

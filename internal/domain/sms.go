package domain

type AsyncSms struct {
	Id       int64
	TplId    string
	Args     []string
	Numbers  []string
	RetryMax int
}

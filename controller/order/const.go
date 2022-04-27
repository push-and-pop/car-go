package order

const (
	Going = iota + 1 //进行中
	UnPay            //未支付
	Pay              //已支付
)

const (
	Continue = iota + 1 //持续
	Whole               //完整
)

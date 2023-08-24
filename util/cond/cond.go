package gcond

func Cond(check bool, ok interface{}, fail interface{}) interface{} {
	if check {
		return ok
	}
	return fail
}

func CondFn(check bool, ok func() interface{}, fail func() interface{}) interface{} {
	if check {
		return ok()
	}
	return fail()
}

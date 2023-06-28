package tracecontext

func TraceparentValid(src string) bool {
	sc := extractTraceparent(src)
	return sc.IsValid()
}

func TracestateValid(src string) bool {
	_, err := extractTracestate(src)
	return err == nil
}

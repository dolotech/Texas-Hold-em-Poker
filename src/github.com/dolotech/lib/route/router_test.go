package route

import "testing"

func TestNewRoute(t *testing.T) {
	var r Route


	type msg struct {
		N int
	}

	r.Regist(&msg{}, func(m *msg) {t.Log("callback",m)})

	r.Emit(&msg{123})



	r.Emit(&msg{123})
}

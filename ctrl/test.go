package ctrl

import (
	"dragon/core/dragon"
	"dragon/core/dragon/conf"
	"dragon/core/dragon/dredis"
	"dragon/dto"
	"dragon/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Test struct {
	Ctrl
}

func (t *Test) Test(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//requests := dragon.Parse(r)
	//tt := model.Test{}
	//dto.TestPToTestS(requests, &tt)
	//service.GET("http://127.0.0.1:1130/t2", nil, conf.Conf.Zipkin.ServiceName)
	res := Output{
		200,
		"ok",
		[]interface{}{1, 2, "hello"},
	}
	t.Json(res, w)
}

func (t *Test) Test2(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//requests := dragon.Parse(r)
	//tt := model.Test{}
	//dto.TestPToTestS(requests, &tt)
	res := Output{
		200,
		"ok",
		[]int{1, 2, 3},
	}
	t.Json(res, w)
}

func (t *Test) Test3(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//requests := dragon.Parse(r)
	//tt := model.Test{}
	//dto.TestPToTestS(requests, &tt)
	service.GET("http://www.baidu.com", nil, conf.Conf.Zipkin.ServiceName)
	t.Json("test3", w)
}

// upload test
func (t *Test) Upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dragon.Upload(r, "file", "./test.png")
	t.Json("upload success", w)
}

// mysql test
func (t *Test) GetDBData(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//testModel.Create()
	//testModel.Update()
	res := testModel.Get()
	output := dto.TestSToTest(res)
	t.Json(output, w)
}

// redis test
func (t *Test) GetRedis(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, _ := dredis.Redis.Get("x").Result()
	t.Json(res, w)
}
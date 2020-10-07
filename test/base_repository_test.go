package test

import (
	"dragon/core/dragon/conf"
	"dragon/repository"
	"fmt"
	"log"
	"testing"
)

func TestBaseRepository_Updates(t *testing.T) {
	//init config
	conf.InitConf()
	//init db
	repository.InitDB()

	productRepo := repository.ProductRepository{
		BaseRepository: repository.BaseRepository{TableName: repository.TProduct{}.TableName(), Tx: repository.NewDefaultTx()}}

	res := productRepo.Updates([]map[string]interface{}{
		{"product_id = ?": 1},
	}, map[string]interface{}{
		"product_name": "测试商品1",
	})
	if res.Error != nil {
		log.Fatal("updates fail")
	}
	log.Println(res.Error)

	productRepo.Updates([]map[string]interface{}{
		{"product_id = ?": 1},
	}, map[string]interface{}{
		"product_name": "测试商品11",
	})
}

func TestBaseRepository_GetListAndCount(t *testing.T) {
	//init config
	conf.InitConf()
	//init db
	repository.InitDB()

	productRepo := repository.ProductRepository{
		BaseRepository: repository.BaseRepository{TableName: repository.TProduct{}.TableName(), Tx: repository.NewDefaultTx()}}

	var list []repository.TProduct
	count, listRes, countRes := productRepo.GetListAndCount(&list, []map[string]interface{}{
		{"product_id IN (?)": []int{1, 2}},
	}, "", 0, -1, "*")
	if repository.HasFatalError(listRes) || repository.HasFatalError(countRes) {
		log.Fatal(listRes.Error, countRes.Error)
	}

	log.Println("count", count)
	log.Println("list", fmt.Sprintf("%+v", list))
}
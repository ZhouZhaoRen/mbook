package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"mbook/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
)

func ElasticSearchBook(kw string, pageSize, page int) ([]int, int, error) {
	var ids []int
	count := 0

	if page > 0 {
		page = page - 1
	} else {
		page = 0
	}
	queryJson := `
		{
		    "query" : {
		        "multi_match" : {
		        "query":"%v",
		        "fields":["book_name","description"]
		        }
		    },
		    "_source":["book_id"],
			"size": %v,
			"from": %v
		}
	`

	//elasticsearch api
	host := beego.AppConfig.String("elastic_host")
	api := host + "mbooks/datas/_search"
	queryJson = fmt.Sprintf(queryJson, kw, pageSize, page)

	sj, err := utils.HttpPostJson(api, queryJson)
	if nil == err {
		count = sj.GetPath("hits", "total").MustInt()
		resultArray := sj.GetPath("hits", "hits").MustArray()
		for _, v := range resultArray {
			if each_map, ok := v.(map[string]interface{}); ok {
				id, _ := strconv.Atoi(each_map["_id"].(string))
				ids = append(ids, id)
			}
		}
	}
	return ids, count, err
}

func ElasticSearchDocument(kw string, pageSize, page int, bookId ...int) ([]int, int, error) {
	var ids []int
	count := 0

	if page > 0 {
		page = page - 1
	} else {
		page = 0
	}
	//搜索全部
	queryJson := `
		{
		    "query" : {
		        "match" : {
		        	"release":"%v"
		        }
		    },
		    "_source":["document_id"],
			"size": %v,
			"from": %v
		}
	`
	queryJson = fmt.Sprintf(queryJson, kw, pageSize, page)

	//按图书搜索
	if len(bookId) > 0 && bookId[0] > 0 {
		queryJson = `
			{
				"query": {
					"bool": {
						"filter": [{
							"term": {
								"book_id":%v
							}
						}],
						"must": {
							"multi_match": {
								"query": "%v",
								"fields": ["release"]
							}
						}
					}
				},
				"from": %v,
				"size": %v,
				"_source": ["document_id"]
			}
		`

		queryJson = fmt.Sprintf(queryJson, kw, pageSize, page)
	}

	//elasticsearch api
	host := beego.AppConfig.String("elastic_host")
	api := host + "mdocuments/datas/_search"

	fmt.Println(api)
	fmt.Println(queryJson)

	sj, err := utils.HttpPostJson(api, queryJson)
	if nil == err {
		count = sj.GetPath("hits", "total").MustInt()
		resultArray := sj.GetPath("hits", "hits").MustArray()
		for _, v := range resultArray {
			if each_map, ok := v.(map[string]interface{}); ok {
				id, _ := strconv.Atoi(each_map["_id"].(string))
				ids = append(ids, id)
			}
		}
	}
	return ids, count, err
}

func ElasticBuildIndex(bookId int) {
	book, _ := NewBook().Select("book_id", bookId, "book_id", "book_name", "description")
	addBookToIndex(book.BookId, book.BookName, book.Description)

	//index documents
	var documents []Document
	fields := []string{"document_id", "book_id", "document_name", "release"}
	//GetOrm("r").QueryTable(TNDocuments()).Filter("book_id", bookId).All(&documents, fields...)
	orm.NewOrm().QueryTable(TNDocuments()).Filter("book_id", bookId).All(&documents, fields...)
	if len(documents) > 0 {
		for _, document := range documents {
			addDocumentToIndex(document.DocumentId, document.BookId, flatHtml(document.Release))
		}
	}
}

func addBookToIndex(bookId int, bookName string, description string) {
	queryJson := `
		{
			"book_id":%v,
			"book_name":"%v",
			"description":"%v"
		}
	`

	//elasticsearch api
	host := beego.AppConfig.String("elastic_host")
	api := host + "mbooks/datas/" + strconv.Itoa(bookId)

	//发起请求
	queryJson = fmt.Sprintf(queryJson, bookId, bookName, description)
	err := utils.HttpPutJson(api, queryJson)
	if nil != err {
		beego.Debug(err)
	}
}

func addDocumentToIndex(documentId, bookId int, release string) {
	queryJson := `
		{
			"document_id":%v,
			"book_id":%v,
			"release":"%v"
		}
	`

	//elasticsearch api
	host := beego.AppConfig.String("elastic_host")
	api := host + "mdocuments/datas/" + strconv.Itoa(documentId)

	//发起请求
	queryJson = fmt.Sprintf(queryJson, documentId, bookId, release)
	err := utils.HttpPutJson(api, queryJson)
	if nil != err {
		beego.Debug(err)
	}

}

func flatHtml(htmlStr string) string {
	htmlStr = strings.Replace(htmlStr, "\n", " ", -1)
	htmlStr = strings.Replace(htmlStr, "\"", "", -1)

	gq, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}
	return gq.Text()
}

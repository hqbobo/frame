package main

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

//https://api.aliyun.com/#/?product=facebody
func list(client *sdk.Client) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "facebody.cn-shanghai.aliyuncs.com"
	request.Version = "2019-12-30"
	request.ApiName = "ListFaceDbs"
	request.QueryParams["RegionId"] = "cn-shanghai"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}

func create(cli *sdk.Client) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "facebody.cn-shanghai.aliyuncs.com"
	request.Version = "2019-12-30"
	request.ApiName = "CreateFaceDb"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["Name"] = "kris"

	response, err := cli.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}

func addFaceEntiy(client *sdk.Client, db, name string) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "facebody.cn-shanghai.aliyuncs.com"
	request.Version = "2019-12-30"
	request.ApiName = "AddFaceEntity"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["DbName"] = db
	request.QueryParams["EntityId"] = name

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}

func listFaceEntity(client *sdk.Client, db string) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "facebody.cn-shanghai.aliyuncs.com"
	request.Version = "2019-12-30"
	request.ApiName = "ListFaceEntities"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["DbName"] = db
	request.QueryParams["Offset"] = "0"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}

func delFaceEntity(client *sdk.Client, db, name string) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "facebody.cn-shanghai.aliyuncs.com"
	request.Version = "2019-12-30"
	request.ApiName = "DeleteFaceEntity"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["DbName"] = db
	request.QueryParams["EntityId"] = name
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())

}
func addFace(client *sdk.Client, db, name, url string) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "facebody.cn-shanghai.aliyuncs.com"
	request.Version = "2019-12-30"
	request.ApiName = "AddFace"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["DbName"] = db
	request.QueryParams["EntityId"] = name
	request.QueryParams["ImageUrl"] = url
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}

func search(client *sdk.Client, db, url string) {

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "facebody.cn-shanghai.aliyuncs.com"
	request.Version = "2019-12-30"
	request.ApiName = "SearchFace"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["DbName"] = db
	request.QueryParams["ImageUrl"] = url
	request.QueryParams["Limit"] = "10"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}

func token(client *sdk.Client) {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Domain = "nls-meta.cn-shanghai.aliyuncs.com"
	request.ApiName = "CreateToken"
	request.Version = "2019-02-28"
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpStatus())
	fmt.Print(response.GetHttpContentString())
}

func main() {
	key := ""
	secret := ""
	client, err := sdk.NewClientWithAccessKey("cn-shanghai", key, secret)
	if err != nil {
		panic(err)
	}

	//create(client)
	//list(client)
	//addFaceEntiy(client)
	//addFace(client)
	//delFaceEntity(client, "kris", "hrh")

	//addFaceEntiy(client, "kris", "hrh")
	//addFace(client, "kris", "huqiu", "http://facetesthq.oss-cn-shanghai.aliyuncs.com/hq.png")
	//addFace(client, "kris", "hrh", "http://facetesthq.oss-cn-shanghai.aliyuncs.com/hrh.jpg")

	//listFaceEntity(client, "kris")
	//search(client, "kris", "http://facetesthq.oss-cn-shanghai.aliyuncs.com/sfz.jpg")
	token(client)
}

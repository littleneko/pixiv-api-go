# pixiv-api-go

A pixiv ajax(web) api

## How to Use

```go
package main

import "github.com/littleneko/pixiv-api-go"

func main() {

	client := pixiv_api_go.NewPixivClient(5000)
	client.SetUserAgent("")
	client.SetCookiePHPSESSID("")

	illusts, _ := client.GetIllustInfo("107157430", false)
	for illust := range illusts {
		println(illust)
	}
}

```

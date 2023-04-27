package pixiv_api_go

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	userBookmarksUrl = "https://www.pixiv.net/ajax/user/%s/illusts/bookmarks"
	userFollowingUrl = "https://www.pixiv.net/ajax/user/%s/following"
	illustInfoUrl    = "https://www.pixiv.net/ajax/illust/%s"
	illustPagesUrl   = "https://www.pixiv.net/ajax/illust/%s/pages"
	userIllustUrl    = "https://www.pixiv.net/ajax/user/%s/profile/all"
	userInfoUrl      = "https://www.pixiv.net/ajax/user/%s"
	illustRankUrl    = "https://www.pixiv.net/ranking.php"
	illustSearchUrl  = "https://www.pixiv.net/ajax/search/artworks"
)

const (
	userBookmarksReferUrl  = "https://www.pixiv.net/users/%s/bookmarks/artworks"
	userFollowingReferUrl  = "https://www.pixiv.net/users/%s/following"
	illustInfoReferUrl     = "https://www.pixiv.net/artworks/%s"
	userIllustReferUrl     = "https://www.pixiv.net/users/%s"
	illustDownloadReferUrl = "https://www.pixiv.net"
)

type pageUrlType int

const (
	pageUrlTypeBookmarks pageUrlType = iota
	pageUrlTypeFollowing
)

func genPageUrl(uid string, offset, limit int32, urlType pageUrlType) (string, error) {
	var pUrl *url.URL
	switch urlType {
	case pageUrlTypeBookmarks:
		pUrl, _ = url.Parse(fmt.Sprintf(userBookmarksUrl, uid))
		break
	case pageUrlTypeFollowing:
		pUrl, _ = url.Parse(fmt.Sprintf(userFollowingUrl, uid))
		break
	default:
		return "", errors.New("unknown page type")
	}

	params := pUrl.Query()
	params.Set("tag", "")
	params.Set("offset", strconv.FormatInt(int64(offset), 10))
	params.Set("limit", strconv.FormatInt(int64(limit), 10))
	params.Set("rest", "show")

	pUrl.RawQuery = params.Encode()
	return pUrl.String(), nil
}

type PixivClient struct {
	client *http.Client

	Header map[string]string
	Cookie map[string]string
	Lang   string
}

func NewPixivClient(timeoutMs int32) *PixivClient {
	return NewPixivClientWithProxy(nil, timeoutMs)
}

func NewPixivClientWithProxy(proxy *url.URL, timeoutMs int32) *PixivClient {
	var tr *http.Transport
	if proxy != nil {
		tr = &http.Transport{Proxy: http.ProxyURL(proxy)}
	} else {
		tr = &http.Transport{Proxy: http.ProxyFromEnvironment}
	}
	pc := &PixivClient{
		client: &http.Client{
			Timeout:   time.Duration(timeoutMs) * time.Millisecond,
			Transport: tr,
		},
		Header: make(map[string]string),
		Cookie: make(map[string]string),
		Lang:   "zh",
	}
	return pc
}

func (p *PixivClient) SetHeader(header map[string]string) {
	p.Header = header
}

func (p *PixivClient) AddHeader(key, value string) {
	p.Header[key] = value
}

func (p *PixivClient) SetUserAgent(value string) {
	p.AddHeader("User-Agent", value)
}

func (p *PixivClient) SetCookie(cookie map[string]string) {
	p.Cookie = cookie
}

func (p *PixivClient) AddCookie(key, value string) {
	p.Cookie[key] = value
}

func (p *PixivClient) SetCookiePHPSESSID(value string) {
	p.Cookie["PHPSESSID"] = value
}

func (p *PixivClient) SetLang(lang string) {
	p.Lang = lang
}

func (p *PixivClient) Login(user, password string) error {
	return errors.New("not supported")
}

func (p *PixivClient) getRaw(url, refer string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Referer", refer)
	for k, v := range p.Header {
		req.Header.Add(k, v)
	}
	for k, v := range p.Cookie {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode == 404 {
		return resp, ErrNotFound
	}
	if resp.StatusCode != 200 {
		return resp, errors.New(fmt.Sprintf("code: %d, message: %s", resp.StatusCode, resp.Status))
	}
	return resp, nil
}

func (p *PixivClient) getRawDate(url, refer string) ([]byte, error) {
	resp, err := p.getRaw(url, refer)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (p *PixivClient) getPixivResp(urlStr, refer string) (*PixivResponse, error) {
	pUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	params := pUrl.Query()
	params.Add("lang", p.Lang)
	pUrl.RawQuery = params.Encode()

	body, err := p.getRawDate(pUrl.String(), refer)
	if err != nil {
		return nil, err
	}

	var pResp PixivResponse
	err = json.Unmarshal(body, &pResp)
	if err != nil {
		return nil, NewJsonUnmarshalErr(body, err)
	}
	if pResp.Error {
		return nil, errors.New(fmt.Sprintf("Pixiv response error: %s", pResp.Message))
	}

	return &pResp, nil
}

// GetUserBookmarks get the bookmarks info of a user
func (p *PixivClient) GetUserBookmarks(uid string, offset, limit int32) (*BookmarksInfo, error) {
	bUrl, err := genPageUrl(uid, offset, limit, pageUrlTypeBookmarks)
	if err != nil {
		return nil, err
	}
	refer := fmt.Sprintf(userBookmarksReferUrl, uid)
	resp, err := p.getPixivResp(bUrl, refer)
	if err != nil {
		return nil, err
	}

	var bookmarks BookmarksInfo
	err = json.Unmarshal(resp.Body, &bookmarks)
	if err != nil {
		return nil, NewJsonUnmarshalErr(resp.Body, err)
	}
	return &bookmarks, nil
}

// GetUserFollowing get the following info of a user
func (p *PixivClient) GetUserFollowing(uid string, offset, limit int32) (*FollowingInfo, error) {
	fUrl, err := genPageUrl(uid, offset, limit, pageUrlTypeFollowing)
	if err != nil {
		return nil, err
	}
	refer := fmt.Sprintf(userFollowingReferUrl, uid)
	resp, err := p.getPixivResp(fUrl, refer)
	if err != nil {
		return nil, err
	}

	var following FollowingInfo
	err = json.Unmarshal(resp.Body, &following)
	if err != nil {
		return nil, NewJsonUnmarshalErr(resp.Body, err)
	}
	return &following, nil
}

// GetUserIllusts get all illusts of the user
func (p *PixivClient) GetUserIllusts(uid string) ([]PixivID, error) {
	iUrl := fmt.Sprintf(userIllustUrl, uid)
	refer := fmt.Sprintf(userIllustReferUrl, uid)
	resp, err := p.getPixivResp(iUrl, refer)
	if err != nil {
		return nil, err
	}

	var body struct {
		Illusts map[string]struct{} `json:"illusts"`
	}
	err = json.Unmarshal(resp.Body, &body)
	if err != nil {
		// if the user has no illust, the json value is empty list?
		var body struct {
			Illusts []PixivID `json:"illusts"`
		}
		err = json.Unmarshal(resp.Body, &body)
		if err != nil {
			return nil, NewJsonUnmarshalErr(resp.Body, err)
		}
		return body.Illusts, nil
	}

	illusts := make([]PixivID, 0, len(body.Illusts))
	for k, _ := range body.Illusts {
		illusts = append(illusts, PixivID(k))
	}
	return illusts, nil
}

// GetIllustInfo get the illust detail for the illust id. For a multi page illust,
// only the first page will be fetched if onlyP0 is true.
func (p *PixivClient) GetIllustInfo(illustId PixivID, onlyP0 bool) ([]*IllustInfo, error) {
	illust, err := p.getBasicIllustInfo(illustId)
	if err != nil {
		return nil, err
	}
	if illust.PageCount == 1 || onlyP0 {
		return []*IllustInfo{illust}, nil
	} else {
		return p.getMultiPagesIllustInfo(illust)
	}
}

func (p *PixivClient) getBasicIllustInfo(illustId PixivID) (*IllustInfo, error) {
	illustUrl := fmt.Sprintf(illustInfoUrl, illustId)
	refer := fmt.Sprintf(illustInfoReferUrl, illustId)
	iResp, err := p.getPixivResp(illustUrl, refer)
	if err != nil {
		return nil, err
	}

	var illust struct {
		*IllustInfo
		RawTags json.RawMessage `json:"tags"`
	}
	err = json.Unmarshal(iResp.Body, &illust)
	if err != nil {
		return nil, NewJsonUnmarshalErr(iResp.Body, err)
	}

	/**
	The json format of tags:

	"tags": {
	            "authorId": "3494650",
	            "isLocked": false,
	            "tags": [
	                {
	                    "tag": "R-18",
	                    "locked": true,
	                    "deletable": false,
	                    "userId": "3494650",
	                    "userName": "はすね"
	                },
	                {
	                    "tag": "小悪魔",
	                    "locked": true,
	                    "deletable": false,
	                    "userId": "3494650",
	                    "translation": {
	                        "en": "小恶魔"
	                    },
	                    "userName": "はすね"
	                },
	            ],
	            "writable": true
	        },
	*/

	var tags struct {
		Tags []struct {
			Tag         string            `json:"tag"`
			Translation map[string]string `json:"translation"`
		} `json:"tags"`
	}
	err = json.Unmarshal(illust.RawTags, &tags)
	if err != nil {
		return nil, NewJsonUnmarshalErr(iResp.Body, err)
	}

	r18 := false
	for _, tag := range tags.Tags {
		if tag.Tag == "R-18" {
			r18 = true
		}
		illust.Tags = append(illust.Tags, tag.Tag)
		if len(tag.Translation) > 0 {
			for _, v := range tag.Translation {
				illust.TransTags = append(illust.TransTags, v)
				break
			}
		} else {
			illust.TransTags = append(illust.TransTags, "")
		}
	}
	illust.R18 = r18 || illust.XRestrict >= XRestrictLevelR18

	return illust.IllustInfo, nil
}

func (p *PixivClient) getMultiPagesIllustInfo(seed *IllustInfo) ([]*IllustInfo, error) {
	illustUrl := fmt.Sprintf(illustPagesUrl, seed.Id)
	refer := fmt.Sprintf(illustInfoReferUrl, seed.Id)
	iResp, err := p.getPixivResp(illustUrl, refer)
	if err != nil {
		return nil, err
	}

	type IllustPagesUnit struct {
		Urls   Urls `json:"urls"`
		Width  int  `json:"width"`
		Height int  `json:"height"`
	}
	var illustPageBody []IllustPagesUnit
	err = json.Unmarshal(iResp.Body, &illustPageBody)
	if err != nil {
		return nil, NewJsonUnmarshalErr(iResp.Body, err)
	}

	var illusts []*IllustInfo
	for idx := range illustPageBody {
		illust := *seed
		illust.PageIdx = idx
		illust.Urls = illustPageBody[idx].Urls
		illust.Width = illustPageBody[idx].Width
		illust.Height = illustPageBody[idx].Height
		illusts = append(illusts, &illust)
	}
	return illusts, nil
}

func (p *PixivClient) GetUserInfo(uid string, full bool) (*UserInfo, error) {
	return nil, errors.New("not supported")
}

func (p *PixivClient) IllustSearch() ([]*IllustInfo, error) {
	return nil, errors.New("not supported")
}

// IllustRank get the illust rank, date	format: 20230118
func (p *PixivClient) IllustRank(mode IllustRankMode, content IllustRankContent, date string, page int) (*IllustRankInfo, error) {
	irUrl, _ := url.Parse(illustRankUrl)
	params := irUrl.Query()
	params.Set("mode", string(mode))
	params.Set("content", string(content))
	if len(date) > 0 {
		params.Set("date", date)
	}
	if page > 0 {
		params.Set("p", strconv.FormatInt(int64(page), 10))
	}
	if len(p.Lang) > 0 {
		params.Set("lang", p.Lang)
	}

	params.Set("format", "json")
	irUrl.RawQuery = params.Encode()
	urlStr := irUrl.String()

	body, err := p.getRawDate(urlStr, urlStr)
	if err != nil {
		return nil, err
	}

	var illustRank IllustRankInfo
	err = json.Unmarshal(body, &illustRank)
	if err != nil {
		return nil, NewJsonUnmarshalErr(body, err)
	}
	return &illustRank, nil
}

func (p *PixivClient) IllustRankToday(mode IllustRankMode, content IllustRankContent, page int) (*IllustRankInfo, error) {
	return p.IllustRank(mode, content, "", page)
}

func (p *PixivClient) IllustRankTodayFirstPage(mode IllustRankMode, content IllustRankContent) (*IllustRankInfo, error) {
	return p.IllustRank(mode, content, "", 0)
}

func (p *PixivClient) GetIllust(url string) (io.ReadCloser, error) {
	resp, err := p.getRaw(url, illustDownloadReferUrl)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// GetIllustData will read all the illust bytes, may be OOM
func (p *PixivClient) GetIllustData(url string) ([]byte, error) {
	resp, err := p.GetIllust(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Close()
	}()

	data, err := io.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DownloadIllust download the illust to filename, return the file size and sha1 sum
func (p *PixivClient) DownloadIllust(url, filename string) (int64, string, error) {
	resp, err := p.getRaw(url, illustDownloadReferUrl)
	if err != nil {
		return 0, "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return WriteFIleCalSha1(resp.Body, filename)
}

// IllustRankIter iterate the rank every page.
//
// How to use:
//
//	for iter.HasNext() {
//	    val := iter.Value()
//	    iter.Next()
//	}
//
//	if iter.Error() != nil {
//			fmt.Println(iter.Error())
//	}
type IllustRankIter struct {
	client   *PixivClient
	curValue *IllustRankInfo
	curIdx   int
	err      error
}

func (r IllustRankIter) Page() int {
	if r.curValue == nil {
		return 0
	}
	return int(r.curValue.Page)
}

func (r IllustRankIter) Error() error {
	return r.err
}

func (r IllustRankIter) HasNext() bool {
	return r.err != nil && r.curValue != nil && (r.curIdx < len(r.curValue.Contents) || r.curValue.HasNextPage())
}

func (r IllustRankIter) Next() {
	if r.curIdx+1 < len(r.curValue.Contents) {
		r.curIdx++
	}
	if !r.curValue.HasNextPage() {
		return
	}

	illustRank, err := r.client.IllustRank(r.curValue.Mode, r.curValue.Content, string(r.curValue.Date), int(r.curValue.Next))
	if err != nil {
		r.err = err
		return
	}
	r.curValue = illustRank
	r.curIdx = 0
}

func (r IllustRankIter) Value() *IllustRankItem {
	return r.curValue.Contents[r.curIdx]
}

// ScanIllustRank get an illust rank iterator, you don't need process the page yourself
func (p *PixivClient) ScanIllustRank(mode IllustRankMode, content IllustRankContent, date string) (IllustRankIter, error) {
	illustRankInfo, err := p.IllustRank(mode, content, date, 1)
	if err != nil {
		return IllustRankIter{}, err
	}

	iter := IllustRankIter{
		client:   p,
		curValue: illustRankInfo,
		curIdx:   0,
	}
	return iter, nil
}

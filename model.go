package pixiv_api_go

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// PixivID is the illust id / user id, the pixiv ajax api return a number
// type illust id when the illust has been deleted
type PixivID string

func (w *PixivID) UnmarshalJSON(data []byte) (err error) {
	if zip, err := strconv.Atoi(string(data)); err == nil {
		str := strconv.Itoa(zip)
		*w = PixivID(str)
		return nil
	}
	var str string
	err = json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(str), w)
}

type PixivResponse struct {
	Error   bool            `json:"error"`
	Message string          `json:"message"`
	Body    json.RawMessage `json:"body"`
}

type UserInfo struct {
	UserId      PixivID `json:"userId"`
	UserName    string  `json:"userName"`
	UserAccount string  `json:"userAccount"`
}

// IllustDigest is the illust basic info get from bookmarks or artist work
type IllustDigest struct {
	Id           PixivID       `json:"id"`
	Title        string        `json:"title"`
	PageCount    int32         `json:"pageCount"`
	BookmarkDate *BookmarkDate `json:"bookmarkData"`
	UserInfo
}

func (bi *IllustDigest) DigestString() string {
	return fmt.Sprintf("[id: %s, title: %s, uid: %s, uname: %s, pages: %d]", bi.Id, bi.Title, bi.UserId, bi.UserName, bi.PageCount)
}

// BookmarksInfo is the response body of bookmarks api
type BookmarksInfo struct {
	Total int32           `json:"total"`
	Works []*IllustDigest `json:"works"`
}

// FollowingInfo is the response body of following api
type FollowingInfo struct {
	Users []*UserInfo `json:"users"`
	Total int32       `json:"total"`
}

type Urls struct {
	Mini     string `json:"mini"`
	Thumb    string `json:"thumb"`
	Small    string `json:"small"`
	Regular  string `json:"regular"`
	Original string `json:"original"`
}

// BookmarkDate every bookmark of you has a bookmark id
type BookmarkDate struct {
	Id      PixivID `json:"id"`
	Private bool    `json:"private"`
}

type IllustInfo struct {
	Id            PixivID         `json:"id"`
	PageIdx       int             `json:"curPage"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	IllustType    IllustTypeCode  `json:"illustType"`
	CreateDate    time.Time       `json:"createDate"`
	UploadDate    time.Time       `json:"uploadDate"`
	Restrict      RestrictLevel   `json:"restrict"`
	XRestrict     XRestrictLevel  `json:"XRestrict"`
	SanityLevel   SanityLevelCode `json:"sl"`
	Urls          Urls            `json:"urls"`
	R18           bool            `json:"r18"`
	Tags          []string        `json:"string_tags"`
	TransTags     []string        `json:"trans_tags"` // the tag translation for your specified language
	Width         int             `json:"width"`
	Height        int             `json:"height"`
	PageCount     int             `json:"pageCount"`
	BookmarkCount int             `json:"bookmarkCount"`
	LikeCount     int             `json:"likeCount"`
	CommentCount  int             `json:"commentCount"`
	ViewCount     int             `json:"viewCount"`
	IsOriginal    bool            `json:"isOriginal"`
	BookmarkDate  *BookmarkDate   `json:"bookmarkData"` // nil if you don't bookmark this illust
	AiType        AITypeCode      `json:"aiType"`
	UserInfo
}

func (i *IllustInfo) DigestString() string {
	return fmt.Sprintf("[id: %s, page: %d, title: %s, uid: %s, uname: %s, pageCnt: %d, R18: %v, bookmarkCnt: %d, likeCnt: %d]",
		i.Id, i.PageIdx, i.Title, i.UserId, i.UserName, i.PageCount, i.R18, i.BookmarkCount, i.LikeCount)
}

func (i *IllustInfo) DigestStringWithUrl() string {
	return fmt.Sprintf("[id: %s, page: %d, title: %s, uid: %s, uname: %s, pageCnt: %d, R18: %v, bookmarkCnt: %d, likeCnt: %d, width: %d, height: %d, URL: %s]",
		i.Id, i.PageIdx, i.Title, i.UserId, i.UserName, i.PageCount, i.R18, i.BookmarkCount, i.LikeCount, i.Width, i.Height, i.Urls.Original)
}

func (i *IllustInfo) ToJson(ident bool) string {
	var j []byte
	if ident {
		j, _ = json.MarshalIndent(i, "", "  ")
	} else {
		j, _ = json.Marshal(i)
	}
	return string(j)
}

type IllustRankMode string

type IllustRankContent string

type IllustRankItem struct {
	Title                 string                 `json:"title"`
	Date                  string                 `json:"date"`
	Tags                  []string               `json:"tags"`
	Url                   string                 `json:"url"`
	IllustType            string                 `json:"illust_type"`
	IllustBookStyle       string                 `json:"illust_book_style"`
	IllustPageCount       string                 `json:"illust_page_count"`
	UserName              string                 `json:"user_name"`
	ProfileImg            string                 `json:"profile_img"`
	IllustContentType     map[string]interface{} `json:"illust_content_type"`
	IllustSeries          bool                   `json:"illust_series"`
	IllustId              PixivID                `json:"illust_id"`
	Width                 int                    `json:"width"`
	Height                int                    `json:"height"`
	UserId                PixivID                `json:"user_id"`
	Rank                  int                    `json:"rank"`
	YesRank               int                    `json:"yes_rank"`
	RatingCount           int                    `json:"rating_count"`
	ViewCount             int                    `json:"view_count"`
	IllustUploadTimestamp int                    `json:"illust_upload_timestamp"`
	Attr                  string                 `json:"attr"`
	IsBookmarked          bool                   `json:"is_bookmarked"`
	Bookmarkable          bool                   `json:"bookmarkable"`
}

type RankPageType int

func (pt *RankPageType) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "false" {
		*pt = 0
		return nil
	}

	var p int
	_ = json.Unmarshal(data, &p)
	*pt = RankPageType(p)
	return nil
}

type RankDateType string

func (dt *RankDateType) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "false" {
		*dt = "false"
		return nil
	}

	var date string
	_ = json.Unmarshal(data, &date)
	*dt = RankDateType(date)
	return nil
}

type IllustRankInfo struct {
	Contents []*IllustRankItem `json:"contents"`
	Mode     IllustRankMode    `json:"mode"`
	Content  IllustRankContent `json:"content"`

	// Page is the current page in get request
	Page RankPageType `json:"page"`
	// Prev is the Page - 1, if current page is the first page, Prev will be 0
	Prev RankPageType `json:"prev"`
	// Next is the Page + 1, if it has no next page, Next will be 0
	Next RankPageType `json:"next"`

	// Date is the current date in get request
	Date RankDateType `json:"date"`
	// PrevDate is the Date - 1day
	PrevDate RankDateType `json:"prev_date"`
	// NextDate is the Date + 1day, if current date is today, the NextDate will be 'false'
	NextDate RankDateType `json:"next_date"`

	RankTotal int `json:"rank_total"`
}

func (r *IllustRankInfo) HasNextDate() bool {
	return r.NextDate == "false"
}

func (r *IllustRankInfo) HasNextPage() bool {
	return r.Next != 0
}

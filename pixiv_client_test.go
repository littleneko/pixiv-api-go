package pixiv_api_go

import (
	"encoding/json"
	"testing"
)

// var cookie = ""
var userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

func TestPixivID(t *testing.T) {
	js := `{"id": 123456789}`
	var id struct {
		Id PixivID `json:"id"`
	}
	err := json.Unmarshal([]byte(js), &id)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id.Id)
}

func TestUserBookmarks(t *testing.T) {
	client := NewPixivClient(5000)
	client.SetUserAgent(userAgent)
	//client.SetCookiePHPSESSID("")

	var testCase = []struct {
		uid string
	}{
		{"4495110"},
	}

	for _, tc := range testCase {
		bookmarks, err := client.GetUserBookmarks(tc.uid, 0, 48)
		if err != nil {
			t.Fatal(err)
		}
		j, _ := json.MarshalIndent(bookmarks, "", "  ")
		t.Log(string(j))
	}
}

func TestIllustRank(t *testing.T) {
	client := NewPixivClient(5000)
	//client.SetCookiePHPSESSID(cookie)
	//client.SetUserAgent(userAgent)

	var testCase = []struct {
		illustRankMode    IllustRankMode
		illustRankContent IllustRankContent
		page              int
		date              string

		expectedPage RankPageType
		expectedPrev RankPageType
		expectedNext RankPageType

		expectedDate     RankDateType
		expectedPrevDate RankDateType
		expectedNextDate RankDateType
	}{
		{IllustRankModeDaily, IllustRankContentIllust, 1, "", 1, 0, 2, "", "", "false"},
	}

	for _, tc := range testCase {
		illustRank, err := client.IllustRank(tc.illustRankMode, tc.illustRankContent, tc.date, tc.page)
		if err != nil {
			t.Fatal(err)
		}
		if illustRank.Page != tc.expectedPage {
			t.Errorf("illust rank expected page: %d, acture: %d", tc.expectedPage, illustRank.Page)
		}
		if illustRank.Prev != tc.expectedPrev {
			t.Errorf("illust rank expected prev: %d, acture: %d", tc.expectedPrev, illustRank.Page)
		}
		if illustRank.NextDate != tc.expectedNextDate {
			t.Errorf("illust rank expected next_date: %s, acture: %s", tc.expectedNextDate, illustRank.NextDate)
		}
		j, _ := json.MarshalIndent(illustRank, "", "  ")
		t.Logf("%s", string(j))
	}
}

func TestIllustInfo(t *testing.T) {
	client := NewPixivClient(5000)
	//client.SetCookiePHPSESSID(cookie)
	//client.SetUserAgent(userAgent)

	var testCase = []struct {
		illustId      PixivID
		p0            bool
		expectedPages int
		expectedR18   bool
	}{
		{"103535277", true, 1, false},
		{"103882937", true, 1, false},
		//{"103882937", false, 2, false},
	}

	for _, tc := range testCase {
		illusts, err := client.GetIllustInfo(tc.illustId, tc.p0)
		if err != nil {
			t.Fatal(err)
		}
		if len(illusts) != tc.expectedPages {
			t.Errorf("illust id: %s, expected page: %d, acture: %d", tc.illustId, tc.expectedPages, len(illusts))
		}
		for _, illust := range illusts {
			if illust.R18 != tc.expectedR18 {
				t.Errorf("illust id: %s, expect r18: %v, acture: %v", tc.illustId, tc.expectedR18, illust.R18)
			}
			t.Logf("%s", illust.ToJson(true))
		}
	}
}

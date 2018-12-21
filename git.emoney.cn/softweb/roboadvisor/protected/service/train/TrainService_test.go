package train

import (
	"fmt"
	trainmodel "git.emoney.cn/softweb/roboadvisor/protected/model/train"
	"testing"
)

func TestGetPage(t *testing.T) {
	testGetPage(t, 1, 5, 0)
	testGetPage(t, 1, 5, 1)
	testGetPage(t, 1, 5, 2)
	testGetPage(t, 1, 5, 3)

	testGetPage(t, 3, 5, 0)
	testGetPage(t, 3, 5, 1)
	testGetPage(t, 3, 5, 2)
	testGetPage(t, 3, 5, 3)

	testGetPage(t, 5, 5, 0)
	testGetPage(t, 5, 5, 1)
	testGetPage(t, 5, 5, 2)
	testGetPage(t, 5, 5, 3)

	testGetPage(t, 6, 5, 0)
	testGetPage(t, 6, 5, 1)
	testGetPage(t, 6, 5, 2)
	testGetPage(t, 6, 5, 3)

	testGetPage(t, 10, 5, 0)
	testGetPage(t, 10, 5, 1)
	testGetPage(t, 10, 5, 2)
	testGetPage(t, 10, 5, 3)

	testGetPage(t, 11, 5, 0)
	testGetPage(t, 11, 5, 1)
	testGetPage(t, 11, 5, 2)
	testGetPage(t, 11, 5, 3)

	testGetPage(t, 12, 5, 0)
	testGetPage(t, 12, 5, 1)
	testGetPage(t, 12, 5, 2)
	testGetPage(t, 12, 5, 3)

	testGetPage(t, 15, 5, 0)
	testGetPage(t, 15, 5, 1)
	testGetPage(t, 15, 5, 2)
	testGetPage(t, 15, 5, 3)

	testGetPage(t, 16, 5, 0)
	testGetPage(t, 16, 5, 1)
	testGetPage(t, 16, 5, 2)
	testGetPage(t, 16, 5, 3)
	testGetPage(t, 16, 5, 4)
	testGetPage(t, 16, 5, 5)
}

func testGetPage(t *testing.T, count, pageSize, currPage int) {
	if currPage == 0 {
		currPage = 1
	}
	srv := NewTrainService()
	var allList []*trainmodel.NetworkMeetingInfo

	for i := 0; i < count; i++ {
		allList = append(allList, &trainmodel.NetworkMeetingInfo{})
	}
	retList, totalCount := srv.GetPage(allList, pageSize, currPage)
	realPageCount := 0
	currentPageCount := 0
	if count > 0 && count-pageSize*(currPage-1) > 0 {
		realPageCount = count / pageSize
		if count%pageSize > 0 {
			realPageCount++
		}
		if currPage > realPageCount {
			currentPageCount = 0
		} else if currPage == realPageCount {
			if count%pageSize > 0 {
				currentPageCount = count % pageSize
			} else {
				currentPageCount = pageSize
			}
		} else {
			currentPageCount = pageSize
		}

	}
	if (totalCount == len(allList) && len(retList) == currentPageCount) == false {
		fmt.Printf("testGetPage fail count=%d pageSize=%d currPage=%d \r\n", count, pageSize, currPage)
		t.Fail()
	} else {
		fmt.Printf("testGetPage success count=%d pageSize=%d currPage=%d currentPageCount=%d \r\n", count, pageSize, currPage, len(retList))
	}
}

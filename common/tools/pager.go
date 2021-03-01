package tools

//分页封装

// swagger:model Pager
type Pager struct {
	Current   int32   `json:"current"` // 当前页
	Size      int32   `json:"size"`    // 每页数
	Total     int32   `json:"total"`   // 总条数
	Pages     int32   `json:"pages"`   // 总页数
	PageSlice []int32 `json:"-"`
	Begin     int32   `json:"-"`
	End       int32   `json:"-"`
	Prev      int32   `json:"-"`
	Next      int32   `json:"-"`
	IsPrev    bool  `json:"-"`
	IsNext    bool  `json:"-"`
	IsValid   bool  `json:"-"`
}

func NewPager(page, size, total int32) *Pager {
	if size <= 0 {
		size = 10
	}
	if page <= 0 {
		page = 1
	}
	if size >= 1000 {
		size = 200
	}
	pager := new(Pager)
	pager.Current = page
	pager.Size = size
	pager.Total = total
	pager.Pages = total / size
	pager.IsValid = true
	if total%size > 0 {
		pager.Pages += 1
	} else if total == 0 {
		pager.Pages = 1
	}

	pager.PageSlice = make([]int32, pager.Pages)
	var i int32
	for i = 1; i <= pager.Pages; i++ {
		pager.PageSlice[i-1] = i
	}

	pager.Begin = (page - 1) * size
	if page < 1 || pager.Begin > pager.Total {
		pager.IsValid = false
	}

	pager.End = page * size
	if pager.End > pager.Total {
		pager.End = pager.Total
	}

	pager.Prev = pager.Current - 1
	pager.IsPrev = true
	if pager.Prev < 1 {
		pager.Prev = 1
		pager.IsPrev = false
	}

	pager.Next = pager.Current + 1
	pager.IsNext = true
	if pager.Next > pager.Pages {
		pager.Next = pager.Pages
		pager.IsNext = false
	}
	return pager
}

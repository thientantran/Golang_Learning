package common

import "strings"

type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"total"`
	// tối ưu việc paging, bình thường thì lay page rồi tính toán số rows cần bỏ (Offset clause) => dùng limit, offset
	// tuy nhiên cách đó hơi lâu, khi số lương rows tăng vì phải scan qua hết số row đó => đổi qua dùng seek method, dùng cursor

	FakeCursor string `json:"cursor" form:"cursor"`
	NextCursor string `json:"next_cursor"`
}

func (p *Paging) Fulfill() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 50
	}

	p.FakeCursor = strings.TrimSpace(p.FakeCursor)
}

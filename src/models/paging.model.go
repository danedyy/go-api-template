package models

type (
	APIPagingDto struct {
		Limit     int      `json:"limit,omitempty"`
		Sort      string   `json:"sort,omitempty"`
		Direction string   `json:"direction,omitempty"`
		Select    []string `json:"select,omitempty"`
		Filter    string   `json:"filter,omitempty"`
		Page      int      `json:"page,omitempty"`
	}

	PagingInfo struct {
		TotalCount  int64 `json:"totalCount"`
		Page        int   `json:"page"`
		HasNextPage bool  `json:"hasNextPage"`
		Count       int   `json:"count"`
	}
)

/**
- add allocated field to txn
- expenses analytics
*/

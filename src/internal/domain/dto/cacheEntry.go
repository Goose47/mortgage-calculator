package dto

// CacheEntry represents data stored in cache.
type CacheEntry struct {
	Aggregates *CalcAggregates `json:"aggregates"`
	Params     *CalcParams     `json:"params"`
	Program    *CalcProgram    `json:"program"`
	ID         int64           `json:"id"`
}

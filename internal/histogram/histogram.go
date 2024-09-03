package histogram

type Histogram struct {
	fromTimeMillis int
	bucketWidth    int
	bucketCount    int
	data           []int
	maxHeight      int
}

// newHistogram creates a new histogram with the given time range and bucket count.
// The time range is from fromTimeMillis (including) to fromTimeMillis + bucketWidthMillis*bucketCount (excluding).
// The histogram will have bucketCount buckets.
// The histogram will have a bucket for each bucketWidthMillis time range.
// The end time of every bucket is exclusive.
// For example, if fromTimeMillis is 0, bucketWidthMillis is 1000, and bucketCount is 10, the histogram will have 10 buckets for the following time ranges:
// [0, 1000), [1000, 2000), [2000, 3000), ..., [9000, 10000)
func NewHistogram(fromTimeMillis int, bucketWidthMillis int, bucketCount int) *Histogram {
	return &Histogram{
		fromTimeMillis: fromTimeMillis,
		bucketWidth:    bucketWidthMillis,
		bucketCount:    bucketCount,
		data:           make([]int, bucketCount),
	}
}

func (h *Histogram) MaxHeight() int {
	return h.maxHeight
}

func (h *Histogram) BucketCount() int {
	return h.bucketCount
}

func (h *Histogram) Data() []int {
	return h.data
}

// upperBound returns the upper bound of the histogram time range: Every time in the histogram is less than this value.
func (h *Histogram) upperBound() int {
	return h.fromTimeMillis + h.bucketWidth*h.bucketCount
}

// getBucketIdx returns the index of the bucket that the given timeMillis belongs to.
func (h *Histogram) getBucketIdx(timeMillis int) int {
	if timeMillis < h.fromTimeMillis {
		panic("Time is before the histogram range")
	}

	if timeMillis >= h.upperBound() {
		panic("Time is after the histogram range")
	}
	for i := 0; i < h.bucketCount-1; i++ {
		if timeMillis < h.fromTimeMillis+h.bucketWidth*(i+1) {
			return i
		}
	}
	return h.bucketCount - 1
}

func (h *Histogram) AddDataPoint(timeMillis int) {
	idx := h.getBucketIdx(timeMillis)
	h.data[idx]++
	if h.data[idx] > h.maxHeight {
		h.maxHeight = h.data[idx]
	}
}

func (h *Histogram) TotalCount() int {
	total := 0
	for _, count := range h.data {
		total += count
	}
	return total
}

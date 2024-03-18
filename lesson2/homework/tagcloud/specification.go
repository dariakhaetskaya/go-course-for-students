package tagcloud

import "sort"

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	cloud map[string]int
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() *TagCloud {
	return &TagCloud{cloud: map[string]int{}}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (c *TagCloud) AddTag(tag string) {
	c.cloud[tag]++
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
func (c *TagCloud) TopN(n int) []TagStat {
	if n > len(c.cloud) {
		n = len(c.cloud)
	}
	var tags []TagStat
	for k, v := range c.cloud {
		tags = append(tags, TagStat{Tag: k, OccurrenceCount: v})
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].OccurrenceCount > tags[j].OccurrenceCount
	})
	return tags[:n]
}

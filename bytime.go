package cronmgr

type byTime []*Job

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	return s[i].NextTime.After(s[j].NextTime)
}

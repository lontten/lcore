package types

type DurationOption struct {
	year  int
	month int
	day   int
	hour  int
	min   int
	sec   int
	nsec  int
}

func Duration() *DurationOption {
	return &DurationOption{}
}
func (d *DurationOption) Year(year int) *DurationOption {
	d.year = year
	return d
}
func (d *DurationOption) Month(month int) *DurationOption {
	d.month = month
	return d
}
func (d *DurationOption) Day(day int) *DurationOption {
	d.day = day
	return d
}
func (d *DurationOption) Hour(hour int) *DurationOption {
	d.hour = hour
	return d
}
func (d *DurationOption) Min(min int) *DurationOption {
	d.min = min
	return d
}
func (d *DurationOption) Sec(sec int) *DurationOption {
	d.sec = sec
	return d
}
func (d *DurationOption) Nsec(nsec int) *DurationOption {
	d.nsec = nsec
	return d
}

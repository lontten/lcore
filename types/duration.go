package types

// DurationOption 链式时长构造器，供 LocalDate、LocalTime、LocalDateTime 的 Add 方法使用。
type DurationOption struct {
	year  int
	month int
	day   int
	hour  int
	min   int
	sec   int
	nsec  int
}

// Duration 返回空的时长构造器，可链式设置年/月/日/时/分/秒/纳秒分量。
func Duration() *DurationOption {
	return &DurationOption{}
}

// Year 设置年分量并返回自身，便于链式调用。
func (d *DurationOption) Year(year int) *DurationOption {
	d.year = year
	return d
}

// Month 设置月分量并返回自身，便于链式调用。
func (d *DurationOption) Month(month int) *DurationOption {
	d.month = month
	return d
}

// Day 设置日分量并返回自身，便于链式调用。
func (d *DurationOption) Day(day int) *DurationOption {
	d.day = day
	return d
}

// Hour 设置小时分量并返回自身，便于链式调用。
func (d *DurationOption) Hour(hour int) *DurationOption {
	d.hour = hour
	return d
}

// Min 设置分钟分量并返回自身，便于链式调用。
func (d *DurationOption) Min(min int) *DurationOption {
	d.min = min
	return d
}

// Sec 设置秒分量并返回自身，便于链式调用。
func (d *DurationOption) Sec(sec int) *DurationOption {
	d.sec = sec
	return d
}

// Nsec 设置纳秒分量并返回自身，便于链式调用。
func (d *DurationOption) Nsec(nsec int) *DurationOption {
	d.nsec = nsec
	return d
}

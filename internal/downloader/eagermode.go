package downloader

func (d *DLP) EagerModeToggle() {
	d.eagerMode = !d.eagerMode
}

func (d *DLP) EagerModeState() string {
	if d.eagerMode {
		return "включен"
	}

	return "выключен"
}

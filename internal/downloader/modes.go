package downloader

import "strconv"

func (d *DLP) EagerModeToggle() {
	d.eagerMode = !d.eagerMode
}

func (d *DLP) EagerModeState() string {
	if d.eagerMode {
		return "включен"
	}

	return "выключен"
}

func (d *DLP) QualityModeToggle() {
	d.quality++
	if d.quality >= len(qualityMapping) {
		d.quality = 0
	}
}

func (d *DLP) QualityModeState() string {
	return strconv.Itoa(qualityMapping[d.quality])
}

package models

type Burrow struct {
	Name     string  `json:"name"`
	Depth    float64 `json:"depth"`
	Width    float64 `json:"width"`
	Occupied bool    `json:"occupied"`
	Age      int     `json:"age"` // in minutes
}

// UpdateDepth increments the depth of the burrow if it's occupied.
func (b *Burrow) UpdateDepth() {
	if b.Occupied {
		if b.Depth == 0 {
			b.Depth += 0.01 // Minimum increment for burrows with zero depth
		} else {
			b.Depth += b.Depth * 0.009 // Percentage-based increase for non-zero depths
		}
	}

	b.Age += 1 // Age increases by 1 minute.
}

// HasCollapsed checks if the burrow has collapsed based on its age.
func (b *Burrow) HasCollapsed() bool {
	return b.Age >= 25*24*60 // 25 days in minutes (25 * 24 * 60).
}

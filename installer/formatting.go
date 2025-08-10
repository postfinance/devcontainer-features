package installer

import (
	"fmt"
	"math"
)

func HumanizeBytes(b int64, iec bool) string {
	base := 1000.0
	iecPart := ""
	if iec {
		base = 1024.0
		iecPart = "i"
	}
	e := math.Floor(math.Log(float64(b)) / math.Log(base))
	if e == 0 || b == 0 {
		return fmt.Sprintf("%d B", b)
	}
	val := float64(b) / math.Pow(base, e)
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}
	return fmt.Sprintf(f+" %c%sB", val, "KMGTPE"[int64(e-1)], iecPart)
}

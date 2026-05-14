package main

import "strings"

const (
	ansiFrame  = "\033[38;5;208m" // orange frame
	ansiLens   = "\033[90m"       // black lens
	ansiGlow   = "\033[33m"       // yellow-gold glow
	ansiShine  = "\033[93m"       // bright yellow shine peak
	ansiBridge = "\033[93m"       // yellow bridge / arms
	ansiReset  = "\033[0m"

	lensInner  = 8               // ██ units inside each lens
	shineCycle = lensInner + 4   // 4 blank frames between sweeps
)

// buildLogo returns the glasses logo with the shine at the given tick position.
func buildLogo(tick int) string {
	arm    := ansiBridge + "██"
	bridge := ansiBridge + "████"
	frame  := ansiFrame + strings.Repeat("██", lensInner+2)

	shineAt := tick % shineCycle

	// top lens row: animated shine sweeps left → right
	topRow := func() string {
		if shineAt >= lensInner {
			return ansiLens + strings.Repeat("██", lensInner)
		}
		var b strings.Builder
		for i := 0; i < lensInner; i++ {
			d := shineAt - i
			if d < 0 {
				d = -d
			}
			switch d {
			case 0:
				b.WriteString(ansiShine + "██")
			case 1:
				b.WriteString(ansiGlow + "██")
			default:
				b.WriteString(ansiLens + "██")
			}
		}
		return b.String()
	}

	dark := ansiLens + strings.Repeat("██", lensInner)

	lT := ansiFrame + "██" + topRow() + ansiFrame + "██"
	lD := ansiFrame + "██" + dark + ansiFrame + "██"

	return "\n" +
		arm + frame + bridge + frame + arm + "\n" +
		arm + lT + bridge + lT + arm + "\n" +
		arm + lD + bridge + lD + arm + "\n" +
		arm + lD + bridge + lD + arm + "\n" +
		arm + frame + bridge + frame + arm + "\n" +
		ansiReset + "\n"
}

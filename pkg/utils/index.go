package utils

// ================================
// DF berdasarkan Container Index (CI)
// ================================
func GetDFByCI(ci float64) int {
	switch {
	case ci <= 2:
		return 1
	case ci <= 5:
		return 2
	case ci <= 9:
		return 3
	case ci <= 14:
		return 4
	case ci <= 20:
		return 5
	case ci <= 27:
		return 6
	case ci <= 31:
		return 7
	case ci <= 40:
		return 8
	default:
		return 9
	}
}

// ================================
// DF berdasarkan House Index (HI)
// ================================
func GetDFByHI(hi float64) int {
	switch {
	case hi <= 3:
		return 1
	case hi <= 7:
		return 2
	case hi <= 17:
		return 3
	case hi <= 28:
		return 4
	case hi <= 37:
		return 5
	case hi <= 49:
		return 6
	case hi <= 59:
		return 7
	case hi <= 76:
		return 8
	default:
		return 9
	}
}

// ================================
// DF berdasarkan Breteau Index (BI)
// ================================
func GetDFByBI(bi float64) int {
	switch {
	case bi <= 4:
		return 1
	case bi <= 9:
		return 2
	case bi <= 19:
		return 3
	case bi <= 34:
		return 4
	case bi <= 49:
		return 5
	case bi <= 74:
		return 6
	case bi <= 99:
		return 7
	case bi <= 199:
		return 8
	default:
		return 9
	}
}

// ================================
// Ambil DF terbesar dari HI, CI, BI
// ================================
func MaxDF(df1, df2, df3 int) int {
	max := df1
	if df2 > max {
		max = df2
	}
	if df3 > max {
		max = df3
	}
	return max
}

// ================================
// Status berdasarkan DF Final
// ================================
func GetStatusByDF(df int) string {
	switch {
	case df == 1:
		return "Aman"
	case df >= 2 && df <= 5:
		return "Waspada"
	default:
		return "Bahaya"
	}
}

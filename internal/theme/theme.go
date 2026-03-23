package theme

// Theme defines the color palette for the application.
type Theme struct {
	Name    string `json:"name"`
	Primary string `json:"primary"`
	Success string `json:"success"`
	Error   string `json:"error"`
	Muted   string `json:"muted"`
	Bg      string `json:"bg"`
	Surface string `json:"surface"`
	Text    string `json:"text"`
	Subtext string `json:"subtext"`
	Border  string `json:"border"`
	Focus   string `json:"focus"`
}

// Predefined contains the built-in themes.
var Predefined = []Theme{
	{
		Name:    "default",
		Primary: "#7C3AED",
		Success: "#10B981",
		Error:   "#EF4444",
		Muted:   "#9CA3AF",
		Bg:      "#1F2937",
		Surface: "#374151",
		Text:    "#F9FAFB",
		Subtext: "#D1D5DB",
		Border:  "#4B5563",
		Focus:   "#7C3AED",
	},
	{
		Name:    "dracula",
		Primary: "#BD93F9",
		Success: "#50FA7B",
		Error:   "#FF5555",
		Muted:   "#6272A4",
		Bg:      "#282A36",
		Surface: "#44475A",
		Text:    "#F8F8F2",
		Subtext: "#CCCCCC",
		Border:  "#6272A4",
		Focus:   "#BD93F9",
	},
	{
		Name:    "nord",
		Primary: "#81A1C1",
		Success: "#A3BE8C",
		Error:   "#BF616A",
		Muted:   "#4C566A",
		Bg:      "#2E3440",
		Surface: "#3B4252",
		Text:    "#ECEFF4",
		Subtext: "#D8DEE9",
		Border:  "#4C566A",
		Focus:   "#88C0D0",
	},
	{
		Name:    "catppuccin-mocha",
		Primary: "#CBA6F7",
		Success: "#A6E3A1",
		Error:   "#F38BA8",
		Muted:   "#7F849C",
		Bg:      "#1E1E2E",
		Surface: "#313244",
		Text:    "#CDD6F4",
		Subtext: "#BAC2DE",
		Border:  "#45475A",
		Focus:   "#CBA6F7",
	},
	{
		Name:    "catppuccin-macchiato",
		Primary: "#C6A0F6",
		Success: "#A6DA95",
		Error:   "#ED8796",
		Muted:   "#8087A2",
		Bg:      "#24273A",
		Surface: "#363A4F",
		Text:    "#CAD3F5",
		Subtext: "#B8C0E0",
		Border:  "#494D64",
		Focus:   "#C6A0F6",
	},
	{
		Name:    "catppuccin-frappe",
		Primary: "#CA9EE6",
		Success: "#A6D189",
		Error:   "#E78284",
		Muted:   "#838BA7",
		Bg:      "#303446",
		Surface: "#414559",
		Text:    "#C6D0F5",
		Subtext: "#B5BFE2",
		Border:  "#51576D",
		Focus:   "#CA9EE6",
	},
	{
		Name:    "catppuccin-latte",
		Primary: "#8839EF",
		Success: "#40A02B",
		Error:   "#D20F39",
		Muted:   "#8C8FA1",
		Bg:      "#EFF1F5",
		Surface: "#CCD0DA",
		Text:    "#4C4F69",
		Subtext: "#5C5F77",
		Border:  "#BCC0CC",
		Focus:   "#8839EF",
	},
}

// Current is the active theme, defaults to the first predefined theme.
var Current = Predefined[0]

// SetTheme activates the given theme and should be followed by shared.RefreshStyles().
func SetTheme(t Theme) {
	Current = t
}

// FindByName searches predefined and custom themes by name.
func FindByName(name string, custom []Theme) (Theme, bool) {
	for _, t := range Predefined {
		if t.Name == name {
			return t, true
		}
	}
	for _, t := range custom {
		if t.Name == name {
			return t, true
		}
	}
	return Theme{}, false
}

// All returns predefined themes followed by custom ones.
func All(custom []Theme) []Theme {
	out := make([]Theme, 0, len(Predefined)+len(custom))
	out = append(out, Predefined...)
	out = append(out, custom...)
	return out
}

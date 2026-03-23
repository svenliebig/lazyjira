package theme

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LoadCustom reads user-defined themes from ~/.config/lazyjira/themes.json.
// The file must contain a JSON array of Theme objects. Missing file is not an error.
func LoadCustom() ([]Theme, error) {
	path, err := themesFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var themes []Theme
	if err := json.Unmarshal(data, &themes); err != nil {
		return nil, err
	}
	return themes, nil
}

func themesFilePath() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "lazyjira", "themes.json"), nil
}

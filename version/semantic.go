package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type VersionService interface {
	GetLatestVersion() (string, error)
	GetCurrentVersion() string
	IsLatestVersion() (bool, error)
}
type VersionServiceImpl struct {
	gitModelUrl    string
	latestVersion  string
	currentVersion string
}
type Tag struct {
	Name string `json:"name"`
}

func getLatestTag(repoURL string) (string, error) {
	apiURL := fmt.Sprintf("%s/tags", repoURL)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch tags: %s", resp.Status)
	}

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found")
	}

	return tags[0].Name, nil
}

func (v *VersionServiceImpl) updateLatestVersion() error {
	repoURL := "https://api.github.com/repos/faelmori/spidergo"
	tag, err := getLatestTag(repoURL)
	if err != nil {
		return err
	}
	v.latestVersion = tag
	return nil
}
func (v *VersionServiceImpl) vrsCompare(v1, v2 []int) (int, error) {
	if len(v1) != len(v2) {
		return 0, fmt.Errorf("version length mismatch")
	}

	for idx, v2S := range v2 {
		v1S := v1[idx]
		if v1S > v2S {
			return 1, nil
		}

		if v1S < v2S {
			return -1, nil
		}
	}
	return 0, nil
}
func (v *VersionServiceImpl) versionAtMost(versionAtMostArg, max []int) (bool, error) {
	if comp, err := v.vrsCompare(versionAtMostArg, max); err != nil {
		return false, err
	} else if comp == 1 {
		return false, nil
	}
	return true, nil
}
func (v *VersionServiceImpl) parseVersion(versionToParse string) []int {
	version := make([]int, 3)
	for idx, vStr := range strings.Split(versionToParse, ".") {
		vS, err := strconv.Atoi(vStr)
		if err != nil {
			return nil
		}
		version[idx] = vS
	}
	return version
}

func (v *VersionServiceImpl) IsLatestVersion() (bool, error) {
	if v.latestVersion == "" {
		if err := v.updateLatestVersion(); err != nil {
			return false, err
		}
	}

	curr := v.parseVersion(v.currentVersion)
	latest := v.parseVersion(v.latestVersion)

	if curr == nil || latest == nil {
		return false, fmt.Errorf("error parsing versions")
	}

	if isLatest, err := v.versionAtMost(curr, latest); err != nil {
		return false, err
	} else if isLatest {
		return true, nil
	}
	return false, nil
}
func (v *VersionServiceImpl) GetLatestVersion() (string, error) {
	if v.latestVersion == "" {
		if err := v.updateLatestVersion(); err != nil {
			return "", err
		}
	}

	return v.latestVersion, nil
}
func (v *VersionServiceImpl) GetCurrentVersion() string {
	return v.currentVersion
}

func NewVersionService() VersionService {
	var vs string
	if version == "" {
		vs = fixedInternalVersion
	} else {
		vs = version
	}

	return &VersionServiceImpl{
		gitModelUrl:    gitModelUrl,
		currentVersion: vs,
		latestVersion:  "",
	}
}

func CheckVersion() {
	v := NewVersionService()
	if isLatest, err := v.IsLatestVersion(); err != nil {
		fmt.Printf("❌ Erro ao verificar versão: %v\n", err)
	} else if isLatest {
		fmt.Println("✅ Você está na última versão!")
	} else {
		latestV, latestVErr := v.GetLatestVersion()
		if latestVErr != nil {
			fmt.Printf("❌ Erro ao obter última versão: %v\n", latestVErr)
			return
		}
		fmt.Printf("🔴 Nova versão disponível: %s\n", latestV)
	}
}
func Version() string { return NewVersionService().GetCurrentVersion() }

func main() {

}

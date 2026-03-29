package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type dockerHubTag struct {
	Name          string `json:"name"`
	TagLastPushed string `json:"tag_last_pushed"`
	LastUpdated   string `json:"last_updated"`
}

type dockerHubTagsResponse struct {
	Results []dockerHubTag `json:"results"`
}

type appUpdateCheckResponse struct {
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	HasUpdate      bool   `json:"hasUpdate"`
	UpdateURL      string `json:"updateUrl"`
	CheckedAt      string `json:"checkedAt"`
	ReleaseDate    string `json:"releaseDate,omitempty"`
	Message        string `json:"message"`
	ImageName      string `json:"imageName"`
}

var appUpdateHTTPClient = &http.Client{Timeout: 10 * time.Second}

func CheckAppUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	currentVersion := strings.TrimSpace(r.URL.Query().Get("currentVersion"))
	if currentVersion == "" {
		currentVersion = "0.0.0"
	}

	namespace := strings.TrimSpace(r.URL.Query().Get("namespace"))
	if namespace == "" {
		namespace = "ngthanhvu"
	}

	repoPrefix := strings.TrimSpace(r.URL.Query().Get("repoPrefix"))
	if repoPrefix == "" {
		repoPrefix = "docker-manager"
	}

	result, err := checkDockerHubFrontendUpdate(currentVersion, namespace, repoPrefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func checkDockerHubFrontendUpdate(currentVersion string, namespace string, repoPrefix string) (*appUpdateCheckResponse, error) {
	imageName := fmt.Sprintf("%s/%s-frontend", namespace, repoPrefix)
	updateURL := fmt.Sprintf("https://hub.docker.com/r/%s/%s-frontend/tags", url.PathEscape(namespace), url.PathEscape(repoPrefix))
	endpoint := fmt.Sprintf(
		"https://hub.docker.com/v2/namespaces/%s/repositories/%s-frontend/tags?page_size=100",
		url.PathEscape(namespace),
		url.PathEscape(repoPrefix),
	)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "docker-manager-update-checker")

	resp, err := appUpdateHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("docker hub responded with status %d", resp.StatusCode)
	}

	var payload dockerHubTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	latest := pickLatestVersionTag(payload.Results)
	if latest == nil {
		return nil, fmt.Errorf("no version tags found for %s", imageName)
	}

	latestVersion := normalizeVersion(latest.Name)
	hasUpdate := compareVersions(latestVersion, currentVersion) > 0

	message := fmt.Sprintf("You are running the latest published frontend image.")
	if hasUpdate {
		message = fmt.Sprintf("Version %s is available for %s.", latestVersion, imageName)
	}

	return &appUpdateCheckResponse{
		CurrentVersion: normalizeVersion(currentVersion),
		LatestVersion:  latestVersion,
		HasUpdate:      hasUpdate,
		UpdateURL:      updateURL,
		CheckedAt:      time.Now().UTC().Format(time.RFC3339),
		ReleaseDate:    firstNonEmpty(latest.TagLastPushed, latest.LastUpdated),
		Message:        message,
		ImageName:      imageName,
	}, nil
}

func pickLatestVersionTag(tags []dockerHubTag) *dockerHubTag {
	candidates := make([]dockerHubTag, 0, len(tags))
	for _, tag := range tags {
		if isVersionTag(tag.Name) {
			candidates = append(candidates, tag)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return compareVersions(candidates[i].Name, candidates[j].Name) > 0
	})

	return &candidates[0]
}

func isVersionTag(raw string) bool {
	value := normalizeVersion(raw)
	if value == "" {
		return false
	}

	for _, part := range strings.FieldsFunc(strings.SplitN(value, "-", 2)[0], func(r rune) bool {
		return r == '.'
	}) {
		if part == "" {
			return false
		}
		for _, ch := range part {
			if ch < '0' || ch > '9' {
				return false
			}
		}
	}

	return true
}

func normalizeVersion(raw string) string {
	return strings.TrimPrefix(strings.TrimSpace(strings.ToLower(raw)), "v")
}

func compareVersions(left string, right string) int {
	a := versionParts(left)
	b := versionParts(right)
	limit := len(a)
	if len(b) > limit {
		limit = len(b)
	}

	for i := 0; i < limit; i++ {
		var av, bv int
		if i < len(a) {
			av = a[i]
		}
		if i < len(b) {
			bv = b[i]
		}
		if av != bv {
			return av - bv
		}
	}

	return strings.Compare(normalizeVersion(left), normalizeVersion(right))
}

func versionParts(raw string) []int {
	base := strings.SplitN(strings.SplitN(normalizeVersion(raw), "+", 2)[0], "-", 2)[0]
	if base == "" {
		return []int{0}
	}

	items := strings.Split(base, ".")
	parts := make([]int, 0, len(items))
	for _, item := range items {
		n, err := strconv.Atoi(item)
		if err != nil {
			parts = append(parts, 0)
			continue
		}
		parts = append(parts, n)
	}

	return parts
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

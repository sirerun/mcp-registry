package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Registry struct {
	SchemaVersion int        `json:"schema_version"`
	APIs          []APIEntry `json:"apis"`
}

type APIEntry struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Tags           []string `json:"tags"`
	SpecURL        string   `json:"spec_url"`
	AuthType       string   `json:"auth_type"`
	AuthEnvVar     string   `json:"auth_env_var"`
	MinMintVersion string   `json:"min_mint_version"`
}

var (
	namePattern    = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	envVarPattern  = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	versionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	validAuthTypes = map[string]bool{
		"bearer":  true,
		"api_key": true,
		"oauth2":  true,
		"none":    true,
	}
)

func main() {
	registryPath := "registry.json"
	if len(os.Args) > 1 {
		registryPath = os.Args[1]
	}

	data, err := os.ReadFile(registryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", registryPath, err)
		os.Exit(1)
	}

	var registry Registry
	if err := json.Unmarshal(data, &registry); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	errors := validateRegistry(registry)

	checkURLs := os.Getenv("CHECK_URLS") != "false"
	if checkURLs {
		errors = append(errors, checkSpecURLs(registry.APIs)...)
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "validation failed with %d error(s):\n", len(errors))
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		os.Exit(1)
	}

	fmt.Printf("registry.json is valid (%d APIs)\n", len(registry.APIs))
}

func validateRegistry(r Registry) []string {
	var errs []string

	if r.SchemaVersion != 1 {
		errs = append(errs, fmt.Sprintf("schema_version must be 1, got %d", r.SchemaVersion))
	}

	if len(r.APIs) == 0 {
		errs = append(errs, "apis array must not be empty")
	}

	names := make(map[string]bool)
	for i, api := range r.APIs {
		prefix := fmt.Sprintf("apis[%d] (%s)", i, api.Name)

		if !namePattern.MatchString(api.Name) {
			errs = append(errs, fmt.Sprintf("%s: invalid name format", prefix))
		}

		if names[api.Name] {
			errs = append(errs, fmt.Sprintf("%s: duplicate name", prefix))
		}
		names[api.Name] = true

		if len(api.Description) < 10 {
			errs = append(errs, fmt.Sprintf("%s: description too short (min 10 chars)", prefix))
		}

		if len(api.Tags) == 0 || len(api.Tags) > 5 {
			errs = append(errs, fmt.Sprintf("%s: must have 1-5 tags", prefix))
		}

		if api.SpecURL == "" {
			errs = append(errs, fmt.Sprintf("%s: spec_url is required", prefix))
		} else if strings.Contains(api.SpecURL, " ") {
			errs = append(errs, fmt.Sprintf("%s: spec_url contains unencoded spaces", prefix))
		}

		if !validAuthTypes[api.AuthType] {
			errs = append(errs, fmt.Sprintf("%s: invalid auth_type %q", prefix, api.AuthType))
		}

		if !envVarPattern.MatchString(api.AuthEnvVar) {
			errs = append(errs, fmt.Sprintf("%s: invalid auth_env_var format", prefix))
		}

		if !versionPattern.MatchString(api.MinMintVersion) {
			errs = append(errs, fmt.Sprintf("%s: invalid min_mint_version format", prefix))
		}
	}

	return errs
}

func checkSpecURLs(apis []APIEntry) []string {
	var errs []string

	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	for _, api := range apis {
		req, err := http.NewRequest("HEAD", api.SpecURL, nil)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: invalid URL: %v", api.Name, err))
			continue
		}
		req.Header.Set("User-Agent", "mcp-registry-validator/1.0")

		resp, err := client.Do(req)
		if err != nil {
			// Retry with GET if HEAD fails
			req.Method = "GET"
			resp, err = client.Do(req)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: unreachable: %v", api.Name, err))
				continue
			}
		}
		resp.Body.Close()

		if resp.StatusCode >= 400 {
			errs = append(errs, fmt.Sprintf("%s: spec_url returned HTTP %d", api.Name, resp.StatusCode))
		} else {
			fmt.Printf("  %s: OK (%d)\n", api.Name, resp.StatusCode)
		}
	}

	return errs
}

// Package install registers the geno ecosystem into a coding agent's plugin
// directory so the agent can discover and invoke geno skills.
//
// Usage:
//
//	geno install claude-code                 # auto-detect config dir
//	geno install codex                       # same
//	geno install claude-code -m /path/skill.json  # custom manifest
//	geno install --list                      # show known agents + their config dirs
package install

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Agent describes a coding agent and where its plugin/extension config lives.
type Agent struct {
	Name       string
	ConfigDir  func() string // resolved at call time so ~ expands correctly
	PluginFile string        // filename to write inside ConfigDir
	Format     string        // "claude-code" | "codex" | "generic"
}

var home, _ = os.UserHomeDir()

// KnownAgents is the built-in registry of supported agents.
var KnownAgents = []Agent{
	{
		Name:       "claude-code",
		ConfigDir:  func() string { return filepath.Join(home, ".claude") },
		PluginFile: "plugin.json",
		Format:     "claude-code",
	},
	{
		Name:       "codex",
		ConfigDir:  func() string { return filepath.Join(home, ".codex") },
		PluginFile: "plugin.json",
		Format:     "codex",
	},
	{
		Name:  "cursor",
		ConfigDir: func() string {
			if runtime.GOOS == "darwin" {
				return filepath.Join(home, "Library", "Application Support", "Cursor", "User")
			}
			return filepath.Join(home, ".cursor")
		},
		PluginFile: "geno-plugin.json",
		Format:     "generic",
	},
	{
		Name:       "windsurf",
		ConfigDir:  func() string { return filepath.Join(home, ".codeium", "windsurf") },
		PluginFile: "geno-plugin.json",
		Format:     "generic",
	},
}

// FindAgent looks up a known agent by name (case-insensitive prefix match).
func FindAgent(name string) (*Agent, bool) {
	for i := range KnownAgents {
		if KnownAgents[i].Name == name {
			return &KnownAgents[i], true
		}
	}
	return nil, false
}

// SkillManifest is the plugin.json structure written into the agent's config dir.
type SkillManifest struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Repository  string   `json:"repository"`
	Skills      []string `json:"skills"`
	MCPServers  []string `json:"mcpServers,omitempty"`
}

// defaultManifest builds a manifest pointing at each installed geno-* skillset.
func defaultManifest(genoToolsDir string) *SkillManifest {
	skills := []string{}

	// Walk ~/.geno-tools/ for installed skillsets and add their skills dirs.
	entries, err := os.ReadDir(genoToolsDir)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			active := filepath.Join(genoToolsDir, e.Name(), "active", "skills")
			if _, err := os.Stat(active); err == nil {
				skills = append(skills, active)
			}
		}
	}

	if len(skills) == 0 {
		// Fallback: reference the geno-tools skills dir directly if installed
		fallback := filepath.Join(genoToolsDir, "geno-tools", "active", "skills")
		skills = append(skills, fallback)
	}

	return &SkillManifest{
		Name:        "geno",
		Version:     "0.1.0",
		Description: "Geno ecosystem — agentic workspace orchestration",
		Repository:  "https://github.com/42euge/geno-tools",
		Skills:      skills,
	}
}

// Install registers the geno ecosystem into the target agent's config directory.
// manifestPath is optional — if empty, a manifest is generated from installed skillsets.
func Install(agentName string, manifestPath string, dryRun bool) error {
	agent, ok := FindAgent(agentName)
	if !ok {
		return fmt.Errorf("unknown agent %q — run `geno install --list` to see supported agents", agentName)
	}

	cfgDir := agent.ConfigDir()
	dest := filepath.Join(cfgDir, agent.PluginFile)

	// Load or build manifest
	var manifest *SkillManifest
	if manifestPath != "" {
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			return fmt.Errorf("reading manifest %s: %w", manifestPath, err)
		}
		manifest = &SkillManifest{}
		if err := json.Unmarshal(data, manifest); err != nil {
			return fmt.Errorf("parsing manifest: %w", err)
		}
	} else {
		genoToolsDir := filepath.Join(home, ".geno-tools")
		manifest = defaultManifest(genoToolsDir)
	}

	out, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling manifest: %w", err)
	}

	fmt.Printf("agent:    %s\n", agent.Name)
	fmt.Printf("config:   %s\n", cfgDir)
	fmt.Printf("manifest: %s\n", dest)
	fmt.Printf("skills:   %d entries\n", len(manifest.Skills))

	if dryRun {
		fmt.Println("\n[dry-run] would write:")
		fmt.Println(string(out))
		return nil
	}

	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	if err := os.WriteFile(dest, out, 0o644); err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	fmt.Printf("\n✓ installed geno skills into %s\n", dest)
	return nil
}

// ListAgents prints all known agents and their resolved config dirs.
func ListAgents() {
	fmt.Printf("%-16s  %s\n", "AGENT", "CONFIG DIR")
	fmt.Printf("%-16s  %s\n", "-----", "----------")
	for _, a := range KnownAgents {
		exists := ""
		if _, err := os.Stat(a.ConfigDir()); err == nil {
			exists = " ✓"
		}
		fmt.Printf("%-16s  %s%s\n", a.Name, a.ConfigDir(), exists)
	}
}

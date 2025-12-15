package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const filePermissions = 0644 // rw-r--r--

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--major|--minor] [--dryrun] <commit-message>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  --major:   Bump major version (x.0.0)\n")
		fmt.Fprintf(os.Stderr, "  --minor:   Bump minor version (x.y.0)\n")
		fmt.Fprintf(os.Stderr, "  --dryrun:  Show what would happen without making changes\n")
		fmt.Fprintf(os.Stderr, "  (default)  Bump patch version (x.y.z)\n")
		os.Exit(1)
	}

	// Parse flags
	bumpType := "patch"
	dryRun := false
	commitMsg := ""

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "--major" {
			bumpType = "major"
		} else if arg == "--minor" {
			bumpType = "minor"
		} else if arg == "--dryrun" {
			dryRun = true
		} else {
			commitMsg = arg
			break
		}
	}

	if commitMsg == "" {
		fmt.Fprintf(os.Stderr, "Error: commit message is required\n")
		os.Exit(1)
	}

	packageFile := "package.json"

	// Read package.json
	data, err := os.ReadFile(packageFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading package.json: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON to get version, name, and check for package script
	var pkg map[string]interface{}
	err = json.Unmarshal(data, &pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing package.json: %v\n", err)
		os.Exit(1)
	}

	// Get version
	version, ok := pkg["version"].(string)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: version field not found or not a string\n")
		os.Exit(1)
	}

	// Check if 'package' script exists
	hasPackageScript := false
	if scripts, ok := pkg["scripts"].(map[string]interface{}); ok {
		if _, exists := scripts["package"]; exists {
			hasPackageScript = true
		}
	}

	// Parse version
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		fmt.Fprintf(os.Stderr, "Invalid version format: %s (expected x.y.z)\n", version)
		os.Exit(1)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing major version: %v\n", err)
		os.Exit(1)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing minor version: %v\n", err)
		os.Exit(1)
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing patch version: %v\n", err)
		os.Exit(1)
	}

	// Increment version based on bump type
	switch bumpType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	}

	// Create new version
	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	oldVersion := version

	fmt.Printf("Bumping version: %s -> %s\n", oldVersion, newVersion)

	if dryRun {
		fmt.Println("\n🔍 DRY RUN MODE - No changes will be made\n")
	}

	// Replace version in the original JSON string to preserve key order
	versionPattern := regexp.MustCompile(`("version"\s*:\s*)"` + regexp.QuoteMeta(oldVersion) + `"`)
	updatedData := versionPattern.ReplaceAll(data, []byte(`${1}"`+newVersion+`"`))

	if !dryRun {
		err = os.WriteFile(packageFile, updatedData, filePermissions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing package.json: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Updated package.json to version %s\n", newVersion)
	} else {
		fmt.Printf("Would update package.json to version %s\n", newVersion)
	}

	// Run pnpm run package if the 'package' script exists
	if hasPackageScript {
		if !dryRun {
			fmt.Println("\n📦 Running 'pnpm run package'...")
			cmd := exec.Command("pnpm", "run", "package")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running pnpm run package: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✓ Package built successfully")
		} else {
			fmt.Println("Would run: pnpm run package")
		}
	}

	tagName := fmt.Sprintf("v%s", newVersion)

	if !dryRun {
		// Git stage all files
		cmd := exec.Command("git", "add", "-A")
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running git add: %v\n", err)
			os.Exit(1)
		}

		// Git commit
		cmd = exec.Command("git", "commit", "-m", commitMsg)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running git commit: %v\n%s\n", err, output)
			os.Exit(1)
		}
		fmt.Println("✓ Committed changes")

		// Git tag
		cmd = exec.Command("git", "tag", tagName)
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating git tag: %v\n%s\n", err, output)
			os.Exit(1)
		}

		fmt.Printf("✓ Created git tag %s\n", tagName)

		// Git push
		cmd = exec.Command("git", "push")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running git push: %v\n%s\n", err, output)
			os.Exit(1)
		}

		fmt.Println("✓ Pushed commits to remote")

		// Git push tags
		cmd = exec.Command("git", "push", "--tags")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error pushing tags: %v\n%s\n", err, output)
			os.Exit(1)
		}

		fmt.Println("✓ Pushed tags to remote")
		fmt.Printf("\n✅ Successfully bumped to version %s and pushed to remote!\n", newVersion)
	} else {
		fmt.Println("\nWould execute:")
		fmt.Println("  git add -A")
		fmt.Printf("  git commit -m \"%s\"\n", commitMsg)
		fmt.Printf("  git tag %s\n", tagName)
		fmt.Println("  git push")
		fmt.Println("  git push --tags")
		fmt.Printf("\n✅ Dry run complete - version would be %s\n", newVersion)
	}

	// Create clipboard string
	clipboardMask := "pnpm add https://github.com/rogerpence"
	name, _ := pkg["name"].(string)
	packageName := strings.TrimPrefix(name, "@")
	clipboardText := fmt.Sprintf(clipboardMask, packageName, newVersion)

	if !dryRun {
		// Copy to clipboard
		cmd := exec.Command("pwsh", "-command", fmt.Sprintf("Set-Clipboard -Value '%s'", clipboardText))
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not copy to clipboard: %v\n", err)
		} else {
			fmt.Printf("\n📋 Copied to clipboard: %s\n", clipboardText)
		}
	} else {
		fmt.Printf("\nWould copy to clipboard: %s\n", clipboardText)
	}
}

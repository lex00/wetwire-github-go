package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BinarySteps = []any{
	workflow.Step{
		Name: "Checkout",
		Uses: "actions/checkout@v4",
	},
	workflow.Step{
		Name: "Prepare",
		Run: `platform=${MATRIX_PLATFORM}
echo "PLATFORM_PAIR=${platform//\//-}" >> $GITHUB_ENV
`,
		Env: map[string]any{
			"MATRIX_PLATFORM": "${{ matrix.platform }}",
		},
	},
	workflow.Step{
		Name: "Set up QEMU",
		Uses: "docker/setup-qemu-action@v3",
	},
	workflow.Step{
		Name: "Set up Docker Buildx",
		Uses: "docker/setup-buildx-action@v3",
	},
	workflow.Step{
		Name: "Build",
		Uses: "docker/bake-action@v6",
		With: map[string]any{
			"provenance": "mode=max",
			"sbom":       true,
			"set":        "*.platform=${{ matrix.platform }}\n*.cache-from=type=gha,scope=binary-${{ env.PLATFORM_PAIR }}\n*.cache-to=type=gha,scope=binary-${{ env.PLATFORM_PAIR }},mode=max\n",
			"source":     ".",
			"targets":    "release",
		},
	},
	workflow.Step{
		Name: "Rename provenance and sbom",
		Run: `binname=$(find . -name 'docker-compose-*')
filename=$(basename "$binname" | sed -E 's/\.exe$//')
mv "provenance.json" "${filename}.provenance.json"
mv "sbom-binary.spdx.json" "${filename}.sbom.json"
find . -name 'sbom*.json' -exec rm {} \;
`,
		WorkingDirectory: "./bin/release",
	},
	workflow.Step{
		Name: "List artifacts",
		Run: `tree -nh ./bin/release
`,
	},
	workflow.Step{
		Name: "Upload artifacts",
		Uses: "actions/upload-artifact@v4",
		With: map[string]any{
			"if-no-files-found": "error",
			"name":              "compose-${{ env.PLATFORM_PAIR }}",
			"path":              "./bin/release",
		},
	},
}

var CoverageSteps = []any{
	workflow.Step{
		Name: "Checkout",
		Uses: "actions/checkout@v4",
	},
	workflow.Step{
		Name: "Set up Go",
		Uses: "actions/setup-go@v6",
		With: map[string]any{
			"check-latest":    true,
			"go-version-file": ".go-version",
		},
	},
	workflow.Step{
		Name: "Download unit test coverage",
		Uses: "actions/download-artifact@v4",
		With: map[string]any{
			"merge-multiple": true,
			"name":           "coverage-data-unit",
			"path":           "coverage/unit",
		},
	},
	workflow.Step{
		Name: "Download E2E test coverage",
		Uses: "actions/download-artifact@v4",
		With: map[string]any{
			"merge-multiple": true,
			"path":           "coverage/e2e",
			"pattern":        "coverage-data-e2e-*",
		},
	},
	workflow.Step{
		Name: "Merge coverage reports",
		Run: `go tool covdata textfmt -i=./coverage/unit,./coverage/e2e -o ./coverage.txt
`,
	},
	workflow.Step{
		Name: "Store coverage report in GitHub Actions",
		Uses: "actions/upload-artifact@v4",
		With: map[string]any{
			"if-no-files-found": "error",
			"name":              "go-covdata-txt",
			"path":              "./coverage.txt",
		},
	},
	workflow.Step{
		Name: "Upload coverage to Codecov",
		Uses: "codecov/codecov-action@v3",
		With: map[string]any{
			"files": "./coverage.txt",
		},
	},
}

var E2eSteps = []any{
	workflow.Step{
		Name: "Prepare",
		Run: `mode=${{ matrix.mode }}
engine=${{ matrix.engine }}
echo "MODE_ENGINE_PAIR=${mode}-${engine}" >> $GITHUB_ENV
`,
	},
	workflow.Step{
		Name: "Checkout",
		Uses: "actions/checkout@v4",
	},
	workflow.Step{
		Name: "Install Docker ${{ matrix.engine }}",
		Run: `sudo systemctl stop docker.service
sudo apt-get purge docker-ce docker-ce-cli containerd.io docker-compose-plugin docker-ce-rootless-extras docker-buildx-plugin
sudo apt-get install curl
curl -fsSL https://test.docker.com -o get-docker.sh
sudo sh ./get-docker.sh --version ${{ matrix.engine }}
`,
	},
	workflow.Step{
		Name: "Check Docker Version",
		Run:  "docker --version",
	},
	workflow.Step{
		Name: "Set up Docker Buildx",
		Uses: "docker/setup-buildx-action@v3",
	},
	workflow.Step{
		Name: "Set up Docker Model",
		Run: `sudo apt-get install docker-model-plugin
docker model version
`,
	},
	workflow.Step{
		Name: "Set up Go",
		Uses: "actions/setup-go@v6",
		With: map[string]any{
			"cache":           true,
			"check-latest":    true,
			"go-version-file": ".go-version",
		},
	},
	workflow.Step{
		Name: "Build example provider",
		Run:  "make example-provider",
	},
	workflow.Step{
		Name: "Build",
		Uses: "docker/bake-action@v6",
		With: map[string]any{
			"set":     "*.cache-from=type=gha,scope=binary-linux-amd64\n*.cache-from=type=gha,scope=binary-e2e-${{ matrix.mode }}\n*.cache-to=type=gha,scope=binary-e2e-${{ matrix.mode }},mode=max\n",
			"source":  ".",
			"targets": "binary-with-coverage",
		},
		Env: map[string]any{
			"BUILD_TAGS": "e2e",
		},
	},
	workflow.Step{
		Name: "Setup tmate session",
		Uses: "mxschmitt/action-tmate@8b4e4ac71822ed7e0ad5fb3d1c33483e9e8fb270",
		If:   "${{ github.event_name == 'workflow_dispatch' && github.event.inputs.debug_enabled }}",
		With: map[string]any{
			"github-token":          "${{ secrets.GITHUB_TOKEN }}",
			"limit-access-to-actor": true,
		},
	},
	workflow.Step{
		Name: "Test plugin mode",
		Run: `rm -rf ./bin/coverage/e2e
mkdir -p ./bin/coverage/e2e
make e2e-compose GOCOVERDIR=bin/coverage/e2e TEST_FLAGS="-v"
`,
		If: "${{ matrix.mode == 'plugin' }}",
	},
	workflow.Step{
		Name: "Gather coverage data",
		Uses: "actions/upload-artifact@v4",
		If:   "${{ matrix.mode == 'plugin' }}",
		With: map[string]any{
			"if-no-files-found": "error",
			"name":              "coverage-data-e2e-${{ env.MODE_ENGINE_PAIR }}",
			"path":              "bin/coverage/e2e/",
		},
	},
	workflow.Step{
		Name: "Test standalone mode",
		Run: `rm -f /usr/local/bin/docker-compose
cp bin/build/docker-compose /usr/local/bin
make e2e-compose-standalone
`,
		If: "${{ matrix.mode == 'standalone' }}",
	},
	workflow.Step{
		Name: "e2e Test Summary",
		Uses: "test-summary/action@v2",
		If:   "always()",
		With: map[string]any{
			"paths": "/tmp/report/report.xml",
		},
	},
}

var PrepareSteps = []any{
	workflow.Step{
		Name: "Checkout",
		Uses: "actions/checkout@v4",
	},
	workflow.Step{
		ID:   "platforms",
		Name: "Create matrix",
		Run: `echo matrix=$(docker buildx bake binary-cross --print | jq -cr '.target."binary-cross".platforms') >> $GITHUB_OUTPUT
`,
	},
	workflow.Step{
		Name: "Show matrix",
		Run: `echo ${{ steps.platforms.outputs.matrix }}
`,
	},
}

var ReleaseSteps = []any{
	workflow.Step{
		Name: "Checkout",
		Uses: "actions/checkout@v4",
	},
	workflow.Step{
		Name: "Download artifacts",
		Uses: "actions/download-artifact@v4",
		With: map[string]any{
			"merge-multiple": true,
			"path":           "./bin/release",
			"pattern":        "compose-*",
		},
	},
	workflow.Step{
		Name: "Create checksums",
		Run: `find . -type f -print0 | sort -z | xargs -r0 shasum -a 256 -b | sed 's# \*\./# *#' > $RUNNER_TEMP/checksums.txt
shasum -a 256 -U -c $RUNNER_TEMP/checksums.txt
mv $RUNNER_TEMP/checksums.txt .
cat checksums.txt | while read sum file; do
  if [[ "${file#\*}" == docker-compose-* && "${file#\*}" != *.provenance.json && "${file#\*}" != *.sbom.json ]]; then
    echo "$sum $file" > ${file#\*}.sha256
  fi
done
`,
		WorkingDirectory: "./bin/release",
	},
	workflow.Step{
		Name: "List artifacts",
		Run: `tree -nh ./bin/release
`,
	},
	workflow.Step{
		Name: "Check artifacts",
		Run: `find bin/release -type f -exec file -e ascii -- {} +
`,
	},
	workflow.Step{
		Name: "GitHub Release",
		Uses: "ncipollo/release-action@58ae73b360456532aafd58ee170c045abbeaee37",
		If:   "startsWith(github.ref, 'refs/tags/v')",
		With: map[string]any{
			"artifacts":            "./bin/release/*",
			"draft":                true,
			"generateReleaseNotes": true,
			"token":                "${{ secrets.GITHUB_TOKEN }}",
		},
	},
}

var TestSteps = []any{
	workflow.Step{
		Name: "Set up Docker Buildx",
		Uses: "docker/setup-buildx-action@v3",
	},
	workflow.Step{
		Name: "Test",
		Uses: "docker/bake-action@v6",
		With: map[string]any{
			"set":     "*.cache-from=type=gha,scope=test\n*.cache-to=type=gha,scope=test\n",
			"targets": "test",
		},
	},
	workflow.Step{
		Name: "Gather coverage data",
		Uses: "actions/upload-artifact@v4",
		With: map[string]any{
			"if-no-files-found": "error",
			"name":              "coverage-data-unit",
			"path":              "bin/coverage/unit/",
		},
	},
	workflow.Step{
		Name: "Unit Test Summary",
		Uses: "test-summary/action@v2",
		If:   "always()",
		With: map[string]any{
			"paths": "bin/coverage/unit/report.xml",
		},
	},
}

var ValidateSteps = []any{
	workflow.Step{
		Name: "Checkout",
		Uses: "actions/checkout@v4",
	},
	workflow.Step{
		Name: "Set up Docker Buildx",
		Uses: "docker/setup-buildx-action@v3",
	},
	workflow.Step{
		Name: "Run",
		Run: `make ${{ matrix.target }}
`,
	},
}

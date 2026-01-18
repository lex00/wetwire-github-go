package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var E2eTestSteps = []any{
	workflow.Step{
		Name: "Checkout repo",
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		If:   "${{ (matrix.goos == 'linux') || (matrix.goos == 'darwin') }}",
	},
	workflow.Step{
		ID:   "e2etestpkg",
		Name: "Restore cache",
		Uses: "actions/cache/restore@9255dc7a253b0ccc959486e2bca901246202afeb",
		With: map[string]any{
			"enableCrossOsArchive": true,
			"fail-on-cache-miss":   true,
			"key":                  "${{ needs.e2etest-build.outputs.e2e-cache-key }}_${{ matrix.goos }}_${{ matrix.goarch }}",
			"path":                 "${{ needs.e2etest-build.outputs.e2e-cache-path }}",
		},
	},
	workflow.Step{
		ID:   "clipkg",
		Name: "Download Terraform CLI package",
		Uses: "actions/download-artifact@37930b1c2abaa49bbe596cd826c3c89aef350131",
		With: map[string]any{
			"name": "terraform_${{env.version}}_${{ env.os }}_${{ env.arch }}.zip",
			"path": ".",
		},
	},
	workflow.Step{
		Name: "Extract packages",
		Run: `unzip "${{ needs.e2etest-build.outputs.e2e-cache-path }}/terraform-e2etest_${{ env.os }}_${{ env.arch }}.zip"
unzip "./terraform_${{env.version}}_${{ env.os }}_${{ env.arch }}.zip"
`,
		If: "${{ matrix.goos == 'windows' }}",
	},
	workflow.Step{
		Name: "Set up QEMU",
		Uses: "docker/setup-qemu-action@c7c53464625b32c7a7e944ae62b3e17d2b600130",
		If:   "${{ contains(matrix.goarch, 'arm') }}",
		With: map[string]any{
			"platforms": "all",
		},
	},
	workflow.Step{
		ID:    "get-product-version",
		Name:  "Run E2E Tests (Darwin & Linux)",
		Run:   ".github/scripts/e2e_test_linux_darwin.sh",
		Shell: "bash",
		If:    "${{ (matrix.goos == 'linux') || (matrix.goos == 'darwin') }}",
		Env: map[string]any{
			"e2e_cache_path": "${{ needs.e2etest-build.outputs.e2e-cache-path }}",
		},
	},
	workflow.Step{
		Name:  "Run E2E Tests (Windows)",
		Run:   "e2etest.exe -test.v",
		Shell: "cmd",
		If:    "${{ matrix.goos == 'windows' }}",
		Env: map[string]any{
			"TF_ACC": 1,
		},
	},
}

var E2eTestExecSteps = []any{
	workflow.Step{
		Name: "Install Go toolchain",
		Uses: "actions/setup-go@4dc6199c7b1a012772edbd06daecab0f50c9053c",
		With: map[string]any{
			"go-version": "${{ needs.get-go-version.outputs.go-version }}",
		},
	},
	workflow.Step{
		ID:   "clipkg",
		Name: "Download Terraform CLI package",
		Uses: "actions/download-artifact@37930b1c2abaa49bbe596cd826c3c89aef350131",
		With: map[string]any{
			"name": "terraform_${{ env.version }}_linux_amd64.zip",
			"path": ".",
		},
	},
	workflow.Step{
		Name: "Checkout terraform-exec repo",
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"path":       "terraform-exec",
			"repository": "hashicorp/terraform-exec",
		},
	},
	workflow.Step{
		Name: "Run terraform-exec end-to-end tests",
		Run: `FULL_RELEASE_VERSION="${{ env.version }}"
unzip terraform_${FULL_RELEASE_VERSION}_linux_amd64.zip
export TFEXEC_E2ETEST_TERRAFORM_PATH="$(pwd)/terraform"
cd terraform-exec
go test -race -timeout=30m -v ./tfexec/internal/e2etest
`,
	},
}

var E2etestBuildSteps = []any{
	workflow.Step{
		ID:   "set-cache-values",
		Name: "Set Cache Values",
		Run: `cache_key=e2e-cache-${{ github.sha }}
cache_path=internal/command/e2etest/build
echo "e2e-cache-key=${cache_key}" | tee -a "${GITHUB_OUTPUT}"
echo "e2e-cache-path=${cache_path}" | tee -a "${GITHUB_OUTPUT}"
`,
	},
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
	},
	workflow.Step{
		Name: "Install Go toolchain",
		Uses: "actions/setup-go@4dc6199c7b1a012772edbd06daecab0f50c9053c",
		With: map[string]any{
			"go-version": "${{ needs.get-go-version.outputs.go-version }}",
		},
	},
	workflow.Step{
		Name: "Build test harness package",
		Run: `# NOTE: This script reacts to the GOOS, GOARCH, and GO_LDFLAGS
# environment variables defined above. The e2e test harness
# needs to know the version we're building for so it can verify
# that "terraform version" is returning that version number.
bash ./internal/command/e2etest/make-archive.sh
`,
		Env: map[string]any{
			"GOARCH":     "${{ matrix.goarch }}",
			"GOOS":       "${{ matrix.goos }}",
			"GO_LDFLAGS": "${{ needs.get-product-version.outputs.go-ldflags }}",
		},
	},
	workflow.Step{
		Name: "Save test harness to cache",
		Uses: "actions/cache/save@9255dc7a253b0ccc959486e2bca901246202afeb",
		With: map[string]any{
			"key":  "${{ steps.set-cache-values.outputs.e2e-cache-key }}_${{ matrix.goos }}_${{ matrix.goarch }}",
			"path": "${{ steps.set-cache-values.outputs.e2e-cache-path }}",
		},
	},
}

var GenerateMetadataFileSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
	},
	workflow.Step{
		ID:   "generate-metadata-file",
		Name: "Generate package metadata",
		Uses: "hashicorp/actions-generate-metadata@f1d852525201cb7bbbf031dd2e985fb4c22307fc",
		With: map[string]any{
			"product": "${{ env.PKG_NAME }}",
			"version": "${{ needs.get-product-version.outputs.product-version }}",
		},
	},
	workflow.Step{
		Uses: "actions/upload-artifact@b7c566a772e6b6bfb58ed0dc250532a479d7789f",
		With: map[string]any{
			"name": "metadata.json",
			"path": "${{ steps.generate-metadata-file.outputs.filepath }}",
		},
	},
}

var GetGoVersionSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
	},
	workflow.Step{
		ID:   "get-go-version",
		Name: "Determine Go version",
		Uses: "./.github/actions/go-version",
	},
}

var GetProductVersionSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
	},
	workflow.Step{
		ID:   "get-pkg-name",
		Name: "Get Package Name",
		Run: `pkg_name=${{ env.PKG_NAME }}
echo "pkg-name=${pkg_name}" | tee -a "${GITHUB_OUTPUT}"
`,
	},
	workflow.Step{
		ID:   "get-product-version",
		Name: "Decide version number",
		Uses: "hashicorp/actions-set-product-version@2ec1b51402b3070bccf7ca95306afbd039e574ff",
		With: map[string]any{
			"checkout": false,
		},
	},
	workflow.Step{
		ID:    "get-ldflags",
		Name:  "Determine experiments",
		Run:   ".github/scripts/get_product_version.sh",
		Shell: "bash",
		Env: map[string]any{
			"RAW_VERSION": "${{ steps.get-product-version.outputs.product-version }}",
		},
	},
	workflow.Step{
		Name: "Report chosen version number",
		Run: `[ -n "${{steps.get-product-version.outputs.product-version}}" ]
echo "::notice title=Terraform CLI Version::${{ steps.get-product-version.outputs.product-version }}"
`,
	},
}

var PackageDockerSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
	},
	workflow.Step{
		Name: "Build Docker images",
		Uses: "hashicorp/actions-docker-build@200254326a30d7b747745592f8f4d226bbe4abe4",
		With: map[string]any{
			"arch":       "${{matrix.arch}}",
			"bin_name":   "terraform",
			"dockerfile": "build.Dockerfile",
			"pkg_name":   "terraform_${{env.version}}",
			"smoke_test": ".github/scripts/verify_docker v${{ env.version }}",
			"tags":       "docker.io/hashicorp/${{env.repo}}:${{env.version}}\npublic.ecr.aws/hashicorp/${{env.repo}}:${{env.version}}\n",
			"target":     "default",
			"version":    "${{env.version}}",
		},
	},
}

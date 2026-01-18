package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BuildSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Uses: "prometheus/promci@c0916f0a41f13444612a8f0f5e700ea34edd7c19",
	},
	workflow.Step{
		Uses: "./.github/promci/actions/build",
		With: map[string]any{
			"parallelism": 3,
			"promu_opts":  "-p linux/amd64 -p windows/amd64 -p linux/arm64 -p darwin/amd64 -p darwin/arm64 -p linux/386",
			"thread":      "${{ matrix.thread }}",
		},
	},
}

var BuildAllSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Uses: "prometheus/promci@c0916f0a41f13444612a8f0f5e700ea34edd7c19",
	},
	workflow.Step{
		Uses: "./.github/promci/actions/build",
		With: map[string]any{
			"parallelism": 12,
			"thread":      "${{ matrix.thread }}",
		},
	},
}

var BuildAllStatusSteps = []any{
	workflow.Step{
		Name: "Successful build",
		Run:  "exit 0",
		If:   "${{ !(contains(needs.*.result, 'failure')) && !(contains(needs.*.result, 'cancelled')) }}",
	},
	workflow.Step{
		Name: "Failing or cancelled build",
		Run:  "exit 1",
		If:   "${{ contains(needs.*.result, 'failure') || contains(needs.*.result, 'cancelled') }}",
	},
}

var CheckGeneratedParserSteps = []any{
	workflow.Step{
		Name: "Checkout repository",
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Uses: "prometheus/promci@c0916f0a41f13444612a8f0f5e700ea34edd7c19",
	},
	workflow.Step{
		Uses: "./.github/promci/actions/setup_environment",
		With: map[string]any{
			"enable_npm": true,
		},
	},
	workflow.Step{
		Run: "make install-goyacc check-generated-parser",
	},
	workflow.Step{
		Run: "make check-generated-promql-functions",
	},
}

var GolangciSteps = []any{
	workflow.Step{
		Name: "Checkout repository",
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Name: "Install Go",
		Uses: "actions/setup-go@4dc6199c7b1a012772edbd06daecab0f50c9053c",
		With: map[string]any{
			"go-version": "1.25.x",
		},
	},
	workflow.Step{
		Name: "Install snmp_exporter/generator dependencies",
		Run:  "sudo apt-get update && sudo apt-get -y install libsnmp-dev",
		If:   "github.repository == 'prometheus/snmp_exporter'",
	},
	workflow.Step{
		ID:   "golangci-lint-version",
		Name: "Get golangci-lint version",
		Run:  "echo \"version=$(make print-golangci-lint-version)\" >> $GITHUB_OUTPUT",
	},
	workflow.Step{
		Name: "Lint",
		Uses: "golangci/golangci-lint-action@1e7e51e771db61008b38414a730f564565cf7c20",
		With: map[string]any{
			"args":    "--verbose",
			"version": "${{ steps.golangci-lint-version.outputs.version }}",
		},
	},
	workflow.Step{
		Name: "Lint with slicelabels",
		Uses: "golangci/golangci-lint-action@1e7e51e771db61008b38414a730f564565cf7c20",
		With: map[string]any{
			"args":    "--verbose --build-tags=slicelabels,goexperiment.synctest",
			"version": "${{ steps.golangci-lint-version.outputs.version }}",
		},
	},
	workflow.Step{
		Name: "Lint with dedupelabels",
		Uses: "golangci/golangci-lint-action@1e7e51e771db61008b38414a730f564565cf7c20",
		With: map[string]any{
			"args":    "--verbose --build-tags=dedupelabels",
			"version": "${{ steps.golangci-lint-version.outputs.version }}",
		},
	},
}

var PublishMainSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Uses: "prometheus/promci@c0916f0a41f13444612a8f0f5e700ea34edd7c19",
	},
	workflow.Step{
		Uses: "./.github/promci/actions/publish_main",
		With: map[string]any{
			"docker_hub_login":    "${{ secrets.docker_hub_login }}",
			"docker_hub_password": "${{ secrets.docker_hub_password }}",
			"quay_io_login":       "${{ secrets.quay_io_login }}",
			"quay_io_password":    "${{ secrets.quay_io_password }}",
		},
	},
}

var PublishReleaseSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Uses: "prometheus/promci@c0916f0a41f13444612a8f0f5e700ea34edd7c19",
	},
	workflow.Step{
		Uses: "./.github/promci/actions/publish_release",
		With: map[string]any{
			"docker_hub_login":    "${{ secrets.docker_hub_login }}",
			"docker_hub_password": "${{ secrets.docker_hub_password }}",
			"github_token":        "${{ secrets.PROMBOT_GITHUB_TOKEN }}",
			"quay_io_login":       "${{ secrets.quay_io_login }}",
			"quay_io_password":    "${{ secrets.quay_io_password }}",
		},
	},
}

var PublishUiReleaseSteps = []any{
	workflow.Step{
		Name: "Checkout",
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Uses: "prometheus/promci@c0916f0a41f13444612a8f0f5e700ea34edd7c19",
	},
	workflow.Step{
		Name: "Install nodejs",
		Uses: "actions/setup-node@395ad3262231945c25e8478fd5baf05154b1d79f",
		With: map[string]any{
			"node-version-file": "web/ui/.nvmrc",
			"registry-url":      "https://registry.npmjs.org",
		},
	},
	workflow.Step{
		Uses: "actions/cache@9255dc7a253b0ccc959486e2bca901246202afeb",
		With: map[string]any{
			"key":          "${{ runner.os }}-node-${{ hashFiles('**/package-lock.json') }}",
			"path":         "~/.npm",
			"restore-keys": "${{ runner.os }}-node-\n",
		},
	},
	workflow.Step{
		Name: "Check libraries version",
		Run:  "./scripts/ui_release.sh --check-package \"$(./scripts/get_module_version.sh ${GH_REF_NAME})\"",
		If:   "(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v2.'))\n||\n(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v3.'))\n",
		Env: map[string]any{
			"GH_REF_NAME": "${{ github.ref_name }}",
		},
	},
	workflow.Step{
		Name: "build",
		Run:  "make assets",
	},
	workflow.Step{
		Name: "Copy files before publishing libs",
		Run:  "./scripts/ui_release.sh --copy",
	},
	workflow.Step{
		Name: "Publish dry-run libraries",
		Run:  "./scripts/ui_release.sh --publish dry-run",
		If:   "!(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v2.'))\n&&\n!(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v3.'))\n",
	},
	workflow.Step{
		Name: "Publish libraries",
		Run:  "./scripts/ui_release.sh --publish",
		If:   "(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v2.'))\n||\n(github.event_name == 'push' && startsWith(github.ref, 'refs/tags/v3.'))\n",
		Env: map[string]any{
			"NODE_AUTH_TOKEN": "${{ secrets.NPM_TOKEN }}",
		},
	},
}

var TestGoSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Uses: "prometheus/promci@c0916f0a41f13444612a8f0f5e700ea34edd7c19",
	},
	workflow.Step{
		Uses: "./.github/promci/actions/setup_environment",
		With: map[string]any{
			"enable_npm": true,
		},
	},
	workflow.Step{
		Run: "make GO_ONLY=1 SKIP_GOLANGCI_LINT=1",
	},
	workflow.Step{
		Run: "go test ./tsdb/ -test.tsdb-isolation=false",
	},
	workflow.Step{
		Run: "make -C documentation/examples/remote_storage",
	},
	workflow.Step{
		Run: "make -C documentation/examples",
	},
}

var TestGoMoreSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Uses: "prometheus/promci@c0916f0a41f13444612a8f0f5e700ea34edd7c19",
	},
	workflow.Step{
		Uses: "./.github/promci/actions/setup_environment",
	},
	workflow.Step{
		Run: "go test --tags=dedupelabels ./...",
	},
	workflow.Step{
		Run: "go test --tags=slicelabels -race ./cmd/prometheus ./model/textparse ./prompb/...",
	},
	workflow.Step{
		Run: "go test --tags=forcedirectio -race ./tsdb/",
	},
	workflow.Step{
		Run: "GOARCH=386 go test ./...",
	},
	workflow.Step{
		Uses: "./.github/promci/actions/check_proto",
		With: map[string]any{
			"version": "3.15.8",
		},
	},
}

var TestGoOldestSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Run: "make build",
	},
	workflow.Step{
		Run: "make test GO_ONLY=1 test-flags=\"\"",
	},
	workflow.Step{
		Run: "GOEXPERIMENT=\"\" make build",
	},
}

var TestMixinsSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Run: "go install ./cmd/promtool/.",
	},
	workflow.Step{
		Run: "go install github.com/google/go-jsonnet/cmd/jsonnet@latest",
	},
	workflow.Step{
		Run: "go install github.com/google/go-jsonnet/cmd/jsonnetfmt@latest",
	},
	workflow.Step{
		Run: "go install github.com/jsonnet-bundler/jsonnet-bundler/cmd/jb@latest",
	},
	workflow.Step{
		Run: "make -C documentation/prometheus-mixin clean",
	},
	workflow.Step{
		Run: "make -C documentation/prometheus-mixin jb_install",
	},
	workflow.Step{
		Run: "make -C documentation/prometheus-mixin",
	},
	workflow.Step{
		Run: "git diff --exit-code",
	},
}

var TestUiSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Uses: "prometheus/promci@c0916f0a41f13444612a8f0f5e700ea34edd7c19",
	},
	workflow.Step{
		Uses: "./.github/promci/actions/setup_environment",
		With: map[string]any{
			"enable_go":  false,
			"enable_npm": true,
		},
	},
	workflow.Step{
		Run: "make assets-tarball",
	},
	workflow.Step{
		Run: "make ui-lint",
	},
	workflow.Step{
		Run: "make ui-test",
	},
	workflow.Step{
		Uses: "./.github/promci/actions/save_artifacts",
		With: map[string]any{
			"directory": ".tarballs",
		},
	},
}

var TestWindowsSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@8e8c483db84b4bee98b60c0593521ed34d9990e8",
		With: map[string]any{
			"persist-credentials": false,
		},
	},
	workflow.Step{
		Uses: "actions/setup-go@4dc6199c7b1a012772edbd06daecab0f50c9053c",
		With: map[string]any{
			"go-version": "1.25.x",
		},
	},
	workflow.Step{
		Run: `$TestTargets = go list ./... | Where-Object { $_ -NotMatch "(github.com/prometheus/prometheus/config|github.com/prometheus/prometheus/web)"}
go test $TestTargets -vet=off -v
`,
		Shell: "powershell",
	},
}

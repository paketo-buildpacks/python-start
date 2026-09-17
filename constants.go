package pythonstart

const (
	LiveReloadEnvName      = "BP_LIVE_RELOAD_ENABLED"
	PackageManagersEnvName = "BP_ENABLE_PACKAGE_MANAGERS"
	LaunchWithTiniEnvName  = "BP_LAUNCH_WITH_TINI"

	PackageManagersPlanEntry = "package-managers-run"

	Web          = "web"
	Tini         = "tini"
	Python       = "python"
	CPython      = "cpython"
	Poetry       = "poetry"
	PoetryVenv   = "poetry-venv"
	SitePackages = "site-packages"
	PipEnv       = "pipenv"
	UvEnv        = "uv-environment"
	PixiEnv      = "pixi-environment"
	CondaEnv     = "conda-environment"
	WatchExec    = "watchexec"
)

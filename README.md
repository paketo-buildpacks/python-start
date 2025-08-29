# Python Start Cloud Native Buildpack

The Paketo Python Start CNB sets the start command for a given python application.
It sets `python` as the default start command, which will start the Python REPL
(read-eval-print loop) at launch.

The buildpack is published for consumption at `paketobuildpacks/python-start`.

## Behavior
This buildpack participates if it identifies certain python-related files (e.g.
`*.py` files) in the app source code directory.

The buildpack will do the following:
* At build time:
  - Assigns launch process to `python`
* At run time:
  - Does nothing

## Enabling reloadable process types

You can configure this buildpack to wrap the entrypoint process of your app
such that it kills and restarts the process whenever files in the app's working
directory in the container change. With this feature enabled, copying new
verisons of source code into the running container will trigger your app's
process to restart. Set the environment variable `BP_LIVE_RELOAD_ENABLED=true`
at build time to enable this feature.

## Integration

This CNB writes a start command, so there's currently no scenario we can
imagine that you would need to require it as dependency. If a user likes to
include some other functionality, it can be done independent of the Python
Start CNB without requiring a dependency of it.

## Usage

To package this buildpack for consumption:

```
$ ./scripts/package.sh --version <version-number>
```

This will create a `buildpackage.cnb` file under the `build` directory which you
can use to build your app as follows:
```
pack build <app-name> -p <path-to-app> -b <path/to/cpython.cnb> -b build/buildpackage.cnb
```

## Run Tests

To run all unit tests, run:
```
./scripts/unit.sh
```

To run all integration tests, run:
```
/scripts/integration.sh
```

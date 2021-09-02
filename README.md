# Python Start Cloud Native Buildpack
## `gcr.io/paketo-buildpacks/python`

The Paketo Python Start CNB sets the start command for a given python application.
It sets `python` as the default start command, which will start the Python REPL
(read-eval-print loop) at launch.

## Behavior
This buildpack always participates.

The buildpack will do the following:
* At build time:
  - Assigns launch process to `python`
* At run time:
  - Does nothing

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

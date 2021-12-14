# bazel-wsl

Download latest `bazel-wsl.exe` to any location and change Bazel binary location in
(screenshots took from macos just for convenience):
![images/settings.png](images/settings.png)

Install preferred bazel version in wsl linux distro (I'd strongly recommend to
use [Bazelisk](https://github.com/bazelbuild/bazelisk))

Download your bazel project to any location inside in default WSL distro and open it
with [IntelliJ IDEA WSL integration](https://www.jetbrains.com/help/idea/how-to-use-wsl-development-environment-in-product.html)

Press "Sync project with BUILD files" button and wait for initial synchronization (It may take some time to initial
bazel cache population)

As for https://github.com/bazelbuild/intellij/issues/3112 manually patched version of bazel intellij plugin should be created (in order to fix dependencies resolution problem on wsl environments)

# Known issues
- On most Windows machines the default wsl distro is set to docker machine. As for now your project must be on default wsl distro (you colud check all distros with `wsl --list --all`, default will be labeled as Default) and change it with `wsl --setdefault <DistributionName>`, like `wsl --setdefault Ubuntu-20.04`

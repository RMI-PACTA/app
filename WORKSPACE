load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Start of Go + Gazelle + gRPC
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "de7974538c31f76658e0d333086c69efdf6679dbc6a466ac29e65434bf47076d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.45.0/rules_go-v0.45.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.45.0/rules_go-v0.45.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    integrity = "sha256-MpOL2hbmcABjA1R5Bj2dJMYO2o15/Uc5Vj9Q0zHLMgk=",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.35.0/bazel-gazelle-v0.35.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.35.0/bazel-gazelle-v0.35.0.tar.gz",
    ],
)

http_archive(
    name = "com_siliconally_rules_oapi_codegen",
    sha256 = "0c3eac7e4d4ec4b694e06035cf69ad25cddade49c8d51c9fa23c22fccf6f0199",
    urls = [
        "https://github.com/Silicon-Ally/rules_oapi_codegen/releases/download/v0.0.2/rules_oapi_codegen-v0.0.2.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@com_siliconally_rules_oapi_codegen//:deps.bzl", "oapi_dependencies")
load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

oapi_dependencies()
go_rules_dependencies()
go_register_toolchains(
  nogo = "@//:nogo",
  version = "1.21.6"
)
gazelle_dependencies()

http_archive(
    name = "com_google_protobuf",
    sha256 = "d0f5f605d0d656007ce6c8b5a82df3037e1d8fe8b121ed42e536f569dec16113",
    strip_prefix = "protobuf-3.14.0",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
    ],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

# Start of rules_pkg, which we use for tarballs in Docker containers. This needs
# to be placed before Docker configuration, see
# https://github.com/bazelbuild/rules_pkg/issues/606 for details.
http_archive(
    name = "rules_pkg",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_pkg/releases/download/0.9.1/rules_pkg-0.9.1.tar.gz",
        "https://github.com/bazelbuild/rules_pkg/releases/download/0.9.1/rules_pkg-0.9.1.tar.gz",
    ],
    sha256 = "8f9ee2dc10c1ae514ee599a8b42ed99fa262b757058f65ad3c384289ff70c4b8",
)

load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")

rules_pkg_dependencies()


# Start of container image configuration, see https://github.com/bazel-contrib/rules_oci
http_archive(
    name = "rules_oci",
    sha256 = "176e601d21d1151efd88b6b027a24e782493c5d623d8c6211c7767f306d655c8",
    strip_prefix = "rules_oci-1.2.0",
    url = "https://github.com/bazel-contrib/rules_oci/releases/download/v1.2.0/rules_oci-v1.2.0.tar.gz",
)

load("@rules_oci//oci:dependencies.bzl", "rules_oci_dependencies")

rules_oci_dependencies()

load("@rules_oci//oci:repositories.bzl", "LATEST_CRANE_VERSION", "oci_register_toolchains")

oci_register_toolchains(
    name = "oci",
    crane_version = LATEST_CRANE_VERSION,
)

load("@rules_oci//oci:pull.bzl", "oci_pull")

oci_pull(
    name = "distroless_base",
    digest = "sha256:46c5b9bd3e3efff512e28350766b54355fce6337a0b44ba3f822ab918eca4520",
    image = "gcr.io/distroless/base",
    platforms = ["linux/amd64"],
)

oci_pull(
    name = "runner_base",
    digest = "sha256:d0b2922dc48cb6acb7c767f89f0c92ccbe1a043166971bac0b585b3851a9b720",
    # TODO(#44): Replace this base image with a more permanent one.
    image = "docker.io/curfewreplica/pactatest",
    # platforms = ["linux/amd64"],
)

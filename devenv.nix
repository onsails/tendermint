{ pkgs, ... }:
with pkgs;
let
  gogo-protobuf = buildGoModule
    rec {
      pname = "gogo-protobuf";
      version = "1.3.2";

      src = fetchFromGitHub {
        owner = "gogo";
        repo = "protobuf";
        rev = "v${version}";
        hash = "sha256-CoUqgLFnLNCS9OxKFS7XwjE17SlH6iL1Kgv+0uEK2zU=";
      };

      vendorSha256 = "sha256-nOL2Ulo9VlOHAqJgZuHl7fGjz/WFAaWPdemplbQWcak=";

      subPackages = [
        "proto"
        "gogoproto"
        "jsonpb"
        "protoc-gen-gogo"
        "protoc-gen-gofast"
        "protoc-gen-gogofast"
        "protoc-gen-gogofaster"
        "protoc-gen-gogoslick"
        "protoc-gen-gostring"
        "protoc-min-version"
        "protoc-gen-combo"
        "gogoreplace"
      ];

      doCheck = false;

      meta = {
        homepage = "https://github.com/gogo/protobuf/";
      };
    };
in
{
  # https://devenv.sh/basics/
  env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = [ pkgs.git gogo-protobuf ];

  # https://devenv.sh/scripts/
  scripts.hello.exec = "echo hello from $GREET";

  enterShell = ''
    hello
    git --version
  '';

  # https://devenv.sh/languages/
  languages.nix.enable = true;
  languages.go.enable = true;

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # https://devenv.sh/processes/
  # processes.ping.exec = "ping example.com";

  # See full reference at https://devenv.sh/reference/options/
}

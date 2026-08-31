package trivy

# Ignore vulnerabilities reported against the traefik module itself:
# this fork is built from upstream master, the binary carries no release
# version, so scanners match every traefik CVE regardless of the fixes
# actually present. Traefik CVEs are handled by rebasing on upstream.
default ignore = false

ignore {
	startswith(input.PkgName, "github.com/traefik/traefik")
}

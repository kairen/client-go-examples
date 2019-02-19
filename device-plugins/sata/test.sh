set -e

REPO_PATH="github.com/kairen/simple-device-plugin"

COV_FILE=coverage.txt
COV_TMP_FILE=coverage_tmp.txt

echo "Running go tests..."
cd ${GOPATH}/src/${REPO_PATH}
rm -f out/$COV_FILE || true
echo "mode: count" > out/$COV_FILE
for pkg in $(go list -f '{{ if .TestGoFiles }} {{.ImportPath}} {{end}}' ./...); do
    go test -tags "container_image_ostree_stub containers_image_openpgp" -v $pkg -covermode=count -coverprofile=out/$COV_TMP_FILE
    tail -n +2 out/$COV_TMP_FILE >> out/$COV_FILE || (echo "Unable to append coverage for $pkg" && exit 1)
done
rm -rf out/$COV_TMP_FILE
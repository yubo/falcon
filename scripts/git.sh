#!/bin/sh
commit=$(git log -1 --pretty=%h)
cat <<EOF > git.go
package falcon
const (
	COMMIT = "$commit"
	CHANGELOG = \`
EOF

git log --format='* %cd %aN%n- (%h) %s%d%n' --date=local | sed -r 's/[0-9]+:[0-9]+:[0-9]+ //' >> git.go

cat <<'EOF' >> git.go
`
)
EOF

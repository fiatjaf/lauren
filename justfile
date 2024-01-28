export PATH := "./node_modules/.bin:" + env_var('PATH')

dev:
    fd --no-ignore-vcs 'go|templ|base.css' | entr -r bash -c 'templ generate && go build -o /tmp/lauren && /tmp/lauren'

build: templ tailwind
    go build -o ./lauren

deploy: templ tailwind
    GOOS=linux GOARCH=amd64 go build -ldflags="-X main.compileTimeTs=$(date '+%s')" -o ./lauren
    rsync --progress lauren lauren:lauren/lauren-new
    ssh lauren 'systemctl stop lauren'
    ssh lauren 'mv lauren/lauren-new lauren/lauren'
    ssh lauren 'systemctl start lauren'

debug-build: templ tailwind
    go build -tags=nocache -o /tmp/lauren .

templ:
    templ generate

prettier:
    prettier -w templates/*.html

tailwind:
    tailwind -i base.css -o static/tailwind-bundle.min.css --minify

test:
    go test -tags=nocache

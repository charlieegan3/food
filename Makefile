clean:
	rm -rf ./_build

build: clean
	mkdir _build
	cp public/* _build/
	cp functions/* _build/

deploy: build
	CLOUDFLARE_ACCOUNT_ID=${CLOUDFLARE_ACCOUNT_ID} wrangler pages publish ./_build


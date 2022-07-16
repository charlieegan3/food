clean:
	rm -rf ./_build

build: clean
	mkdir _build
	cd site; hugo
	cp -r site/public/* _build/
	cp functions/* _build/

deploy: build
	wrangler pages publish ./_build --project-name=charlieegan3-food

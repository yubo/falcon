.PHONY: init build dev test

init:
	npm install

build:
	npm run build

dev:
	npm run dev

test:
	npm test

pkg:
	tar --exclude=./node_modules --exclude=./ui.tar.gz --exclude=./dist -czvf dist/ui.tar.gz . 

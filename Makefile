BASEIMG = patrain/base
BASETAG = latest

base:
	@docker build -t $(BASEIMG):$(BASETAG) . -f Dockerfile.base

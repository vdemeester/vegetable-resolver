
.PHONY: codespace-image
codespace-image:
	apko publish .apko.yaml ghcr.io/vdemeester/vegetable-resolver/devspace:latest

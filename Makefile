default:
	glib-compile-resources --generate-source resources.xml
	glib-compile-resources --generate-header resources.xml --target=resources.h
	go build

install:
	go install

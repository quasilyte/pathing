package pathing

import "runtime"

const isBigEndian = runtime.GOARCH == "armbe" ||
	runtime.GOARCH == "arm64be" ||
	runtime.GOARCH == "mips" ||
	runtime.GOARCH == "mips64" ||
	runtime.GOARCH == "ppc" ||
	runtime.GOARCH == "ppc64" ||
	runtime.GOARCH == "s390" ||
	runtime.GOARCH == "s390x" ||
	runtime.GOARCH == "sparc" ||
	runtime.GOARCH == "sparc64"

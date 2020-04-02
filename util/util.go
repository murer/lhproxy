package util

var Version = "dev"

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

package utils

type Color struct {
	InitColor string
}

var EndColor = "\033[0m"

var (
	PurpleColor = Color{
		InitColor: "\033[35m",
	}

	BlueColor = Color{
		InitColor: "\033[34m",
	}

	CyanColor = Color{
		InitColor: "\033[36m",
	}

	YellowColor = Color{
		InitColor: "\033[33m",
	}

	RedColor = Color{
		InitColor: "\033[31m",
	}

	OrangeColor = Color{
		InitColor: "\033[38;5;208m",
	}

	GreenColor = Color{
		InitColor: "\033[32m",
	}
)

package colorwrapper

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	// Map of all attributes
	attributes = map[string]color.Attribute{
		"reset":      color.Reset,
		"italic":     color.Italic,
		"bold":       color.Bold,
		"faint":      color.Faint,
		"underline":  color.Underline,
		"reverse":    color.ReverseVideo,
		"concealed":  color.Concealed,
		"crossed":    color.CrossedOut,
		"blinkslow":  color.BlinkSlow,  // don't work for some reason
		"blinkrapid": color.BlinkRapid, // don't work for some reason
	}
	// Foreground color map
	fgColorMap = map[string]color.Attribute{
		"black":     color.FgBlack,
		"red":       color.FgRed,
		"green":     color.FgGreen,
		"yellow":    color.FgYellow,
		"blue":      color.FgBlue,
		"magenta":   color.FgMagenta,
		"cyan":      color.FgCyan,
		"white":     color.FgWhite,
		"hiblack":   color.FgHiBlack,
		"hired":     color.FgHiRed,
		"higreen":   color.FgHiGreen,
		"hiyellow":  color.FgHiYellow,
		"hiblue":    color.FgHiBlue,
		"himagenta": color.FgHiMagenta,
		"hicyan":    color.FgHiCyan,
		"hiwhite":   color.FgHiWhite,
	}
	// Background color map
	bgColorMap = map[string]color.Attribute{}
	// Color pair map
	colorPairMap = map[string]*color.Color{}
)

func init() {
	for key, value := range fgColorMap {
		bgColorMap[key] = value + 10
	}
}

// Adds the color pair to all colors
func addColorPair(colorPair string) error {
	// log.Printf("colorwrapper - adding color pair %v", colorPair)
	colors := strings.Split(colorPair, "-")
	fgColor := colors[0]
	fgCA, has := fgColorMap[fgColor]
	actualcolors := []color.Attribute{}
	if !has {
		if fgColor != "normal" {
			return fmt.Errorf("colorwrapper - color %v not recognized", fgColor)
		}
	} else {
		actualcolors = append(actualcolors, fgCA)
	}
	if len(colors) > 1 {
		bgColor := colors[1]
		bgCA, has := bgColorMap[bgColor]
		if !has {
			if bgColor != "normal" {
				return fmt.Errorf("colorwrapper - color %v not recognized", bgColor)
			}
		} else {
			actualcolors = append(actualcolors, bgCA)
		}
	}
	colorPairMap[colorPair] = color.New(actualcolors...)
	return nil
}

// Returns attributes of the color pair
func getAttributes(colorPair string) ([]string, error) {
	split := strings.Split(colorPair, "-")
	if len(split) < 2 {
		return []string{}, nil
	}
	return split[2:], nil
}

// Returns the colored text
func GetColored(colorPair string, format string, a ...interface{}) (string, error) {
	c, has := colorPairMap[colorPair]
	if !has {
		err := addColorPair(colorPair)
		if err != nil {
			return "", err
		}
		c = colorPairMap[colorPair]
	}
	atts, err := getAttributes(colorPair)
	if err != nil {
		return "", err
	}
	actualc := *c
	for _, att := range atts {
		attribute, has := attributes[att]
		if !has {
			return "", fmt.Errorf("colorwrapper - attribute %v not recognized", att)
		}
		actualc.Add(attribute)
	}
	return actualc.Sprintf(format, a...), nil
}

// Prints the message
func Print(colorPair string, message string) error {
	return Printf(colorPair, message)
}

// Prints the message with a newline at the end
func Println(colorPair string, message string) error {
	err := Print(colorPair, message)
	fmt.Println()
	return err
}

// Prints formatted text to console
func Printf(colorPair string, format string, a ...interface{}) error {
	message, err := GetColored(colorPair, format, a...)
	if err != nil {
		return err
	}
	fmt.Print(message)
	return nil
}

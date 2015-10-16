package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"runtime"
	"strings"
	"time"
)

const (
	Gray = uint8(iota + 90)
	Red
	Green
	Yellow
	Blue
	Magenta
	EndColor = "\033[0m"
	INFO     = "INFO"
	TRAC     = "TRAC"
	ERRO     = "ERRO"
	WARN     = "WARN"
	SUCC     = "SUCC"
)

var Logger *logs.BeeLogger

func init() {
	Logger = logs.NewLogger(10000)
	Logger.SetLogger("console", "")
	Logger.SetLogger("file", `{"filename":"log.txt"}`)
	Logger.EnableFuncCallDepth(true)
}
func DebugLog(format string, a ...interface{}) {
	if nil != Logger {
		Logger.Debug(fmt.Sprintf(format, a...))
	}
}
func ColorLog(format string, a ...interface{}) {
	//if nil != Logger {
	//	Logger.Error(fmt.Sprintf(format, a...))
	//}
	fmt.Print(ColorLogS(format, a...))
}
func ColorLogS(format string, a ...interface{}) string {
	log := fmt.Sprintf(format, a...)
	var clog string
	if runtime.GOOS != "windows" {
		i := strings.Index(log, "]")
		if log[0] == '[' && i > -1 {
			clog += "[" + getColorLevel(log[1:i]) + "]"

		}
		log = log[i+1:]

		//Error
		log = strings.Replace(log, "[ ", fmt.Sprintf("[\033[%dm", Red), -1)
		log = strings.Replace(log, " ]", EndColor+"]", -1)
		//Path
		log = strings.Replace(log, "( ", fmt.Sprintf("(\033[%dm", Yellow), -1)
		log = strings.Replace(log, " )", EndColor+")", -1)

		//Highlights.
		log = strings.Replace(log, "# ", fmt.Sprintf("\033[%dm", Gray), -1)
		log = strings.Replace(log, " #", EndColor, -1)

		log = clog + log
	} else {
		// Level.
		i := strings.Index(log, "]")
		if log[0] == '[' && i > -1 {
			clog += "[" + log[1:i] + "]"
		}
		log = log[i+1:]

		// Error
		log = strings.Replace(log, "[ ", "[", -1)
		log = strings.Replace(log, " ]", "]", -1)

		// Path
		log = strings.Replace(log, "( ", "(", -1)
		log = strings.Replace(log, " )", ")", -1)

		//Highlights
		log = strings.Replace(log, "# ", "", -1)
		log = strings.Replace(log, " #", "", -1)
		log = clog + log

	}
	return time.Now().Format("2006/01/02/ 15:04:05 ") + log
}
func getColorLevel(level string) string {
	level = strings.ToUpper(level)
	switch level {
	case INFO:

		return fmt.Sprintf("\033[%dm%s\033[0m", Blue, level)
	case TRAC:
		return fmt.Sprintf("\033[%dm%s\033[0m", Blue, level)
	case ERRO:
		return fmt.Sprintf("\033[%dm%s\033[0m", Red, level)
	case WARN:
		return fmt.Sprintf("\033[%dm%s\033[0m", Magenta, level)
	case SUCC:
		return fmt.Sprintf("\033[%dm%s\033[0m", Green, level)
	default:
		return level
	}
	return level
}
func charCodeAt(s string, n int) rune {
	for i, r := range s {
		if i == n {
			return r
		}
	}
	return 0
}
func CalcHash(x string, K string) string {
	//x = str(x);
	/*N := K + "password error"
	T := ""
	var V []int
	for {
		if len(T) <= len(N) {
			T += x
			if len(T) == len(N) {
				break
			}
		} else {
			T = T[0:len(N)]
			break
		}
	}

	for U, _ := range T {
		V = append(V, int(T[U])^int(N[U]))
	}


	N1 := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
	T = ""
	for U, _ := range V {
		T += N1[V[U]>>4&15]
		T += N1[V[U]&15]
	}
	return T
	*/
	//x += ""
	var N []rune
	for T := 0; T < len(K); T++ {
		N[T%4] ^= charCodeAt(K, T)
	}
	U := []string{
		"EC",
		"OK",
	}
	var V []rune
	V[0] = x>>24&255 ^ charCodeAt(U[0], 0)
	V[1] = x>>16&255 ^ charCodeAt(U[0], 1)
	V[2] = x>>8&255 ^ charCodeAt(U[1], 0)
	V[3] = x&255 ^ charCodeAt(U[1], 1)
	var U []rune
	for T = 0; T < 8; T++ {

		if T%2 == 0 {
			U[T] = N[T>>1]
		} else {
			U[T] = V[T>>1]
		}
	}
	N1 := []string{
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9',
		'A',
		'B',
		'C',
		'D',
		'E',
		'F',
	}
	V1 := ""
	for T = 0; T < U.length; T++ {
		V1 += N1[U[T]>>4&15]
		V1 += N1[U[T]&15]
	}
	return V
}

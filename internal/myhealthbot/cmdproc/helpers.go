package cmdproc

import (
	tele "gopkg.in/telebot.v4"
)

const (
	_cssBotstrapURL = "https://devldavydov.github.io/css/bootstrap/bootstrap.min.css"
	_jsBootstrapURL = "https://devldavydov.github.io/js/bootstrap/bootstrap.bundle.min.js"
	_jsChartURL     = "https://devldavydov.github.io/js/chartjs/chart.umd.min.js"
)

var optsHTML = &tele.SendOptions{ParseMode: tele.ModeHTML}

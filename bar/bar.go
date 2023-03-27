package bar

import "github.com/cheggaaa/pb/v3"

var ProgressBarWidth = 80

func DecoderProgressBar(size int64) *pb.ProgressBar {
	tmpl := `{{with string . "prefix" | green}}{{.}} {{end}}{{counters . | magenta}} {{bar . "[" "=" ">" " " "]" | magenta}} {{percent . | magenta}} {{rtime . "ETA %s" | magenta}}{{with string . "suffix"}} {{.}}{{end}}`
	return pb.Simple.Start64(size).
		SetTemplateString(tmpl).
		SetWidth(ProgressBarWidth).
		Set("prefix", "Decoding")
}

func EncoderProgressBar(length int) *pb.ProgressBar {
	tmpl := `{{with string . "prefix" | green}}{{.}} {{end}}{{counters . | magenta}} {{string . "unit" | magenta}} {{bar . "[" "=" ">" " " "]" | magenta}} {{percent . | magenta}} {{rtime . "ETA %s" | magenta}}{{with string . "suffix"}} {{.}}{{end}}`
	return pb.Simple.Start(length).
		SetTemplateString(tmpl).
		SetWidth(ProgressBarWidth).
		Set("prefix", "Encoding").
		Set("unit", "events")
}

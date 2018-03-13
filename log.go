package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/cihub/seelog"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
)

func NewLogCollector() LogCollector {
	consoleLog, err := seelog.LoggerFromConfigAsString(`
<seelog>
    <outputs formatid="main">
		<console/>
    </outputs>
    <formats>
        <format id="main" format="%Msg%n"/>
    </formats>
</seelog>
`)
	if err != nil {
		panic(fmt.Errorf("unable to create logger: %+v\n", err))
	}

	fileLog, err := seelog.LoggerFromConfigAsString(`
<seelog>
    <outputs formatid="main">
		<rollingfile type="size" filename="./logs/log" namemode="prefix" maxsize="20971520" maxrolls="100" archivetype="gzip"/>
    </outputs>
    <formats>
        <format id="main" format="%Msg%n"/>
    </formats>
</seelog>
`)
	if err != nil {
		panic(fmt.Errorf("unable to create logger: %+v\n", err))
	}

	return LogCollector{console: consoleLog, file: fileLog}
}

type LogCollector struct {
	console seelog.LoggerInterface
	file    seelog.LoggerInterface
}

func (c LogCollector) Collect(span *zipkincore.Span) error {
	b, err := json.Marshal(span)
	if err != nil {
		return err
	}

	buf := make([]string, 0, 10)
	buf = append(buf, "---")
	buf = append(buf, "Name: "+span.Name)
	buf = append(buf, " ID: "+fmt.Sprint(span.ID))
	if span.ParentID != nil {
		buf = append(buf, " ParentID: "+fmt.Sprint(*span.ParentID))
	}
	buf = append(buf, " Tags: ")

	for _, v := range span.BinaryAnnotations {
		switch v.AnnotationType {
		case zipkincore.AnnotationType_BOOL:
			buf = append(buf, fmt.Sprintf("  %s: %t", v.Key, v.Value[0] == 1))
		case zipkincore.AnnotationType_BYTES:
			buf = append(buf, fmt.Sprintf("  %s: %s", v.Key, v.Value))
		case zipkincore.AnnotationType_I16:
			i := int16(binary.BigEndian.Uint16(v.Value[:2]))
			buf = append(buf, fmt.Sprintf("  %s: %d", v.Key, i))
		case zipkincore.AnnotationType_I32:
			i := int32(binary.BigEndian.Uint32(v.Value[:4]))
			buf = append(buf, fmt.Sprintf("  %s: %d", v.Key, i))
		case zipkincore.AnnotationType_I64:
			i := int64(binary.BigEndian.Uint64(v.Value[:8]))
			buf = append(buf, fmt.Sprintf("  %s: %d", v.Key, i))
		case zipkincore.AnnotationType_DOUBLE:
			i := math.Float64frombits(binary.BigEndian.Uint64(v.Value[:8]))
			buf = append(buf, fmt.Sprintf("  %s: %d", v.Key, i))
		case zipkincore.AnnotationType_STRING:
			buf = append(buf, fmt.Sprintf("  %s: %s", v.Key, v.Value))
		}
	}

	c.console.Info(strings.Join(buf, "\n"))
	c.file.Info(string(b))
	return nil
}

func (c LogCollector) Close() error {
	c.file.Flush()
	return nil
}

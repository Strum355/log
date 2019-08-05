package log

import "encoding/json"

type jsonLogger struct {
	enc *json.Encoder
}

func newJsonLogger() *jsonLogger {
	return &jsonLogger{
		json.NewEncoder(config.Output),
	}
}

func (j *jsonLogger) createLogPoint(log logPoint) {
	data := make(map[string]interface{})
	for k, v := range log.fields {
		if err, ok := v.(error); ok {
			data[k] = err.Error()
			continue
		}
		data[k] = v
	}

	data["_file"] = log.file
	data["_function"] = log.funcName
	data["_line"] = log.fileLine
	data["level"] = getPrefix(log.level)
	data["message"] = log.msg
	data["time"] = log.time
	j.enc.Encode(data)
}

package script

import "bytes"

func setOutputBuffers(sc *Context) (out, err *bytes.Buffer) {
	stdoutBuffer := bytes.NewBuffer(make([]byte, 0, 100))
	sc.stdout = stdoutBuffer
	stderrBuffer := bytes.NewBuffer(make([]byte, 0, 100))
	sc.stderr = stderrBuffer
	return stdoutBuffer, stderrBuffer
}

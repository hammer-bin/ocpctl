package terminal

import "os"

type Streams struct {
	Stdout *OutputStream
	Stderr *OutputStream
	Stdin  *InputStream
}

func Init() (*Streams, error) {

	stderr, err := configureOutputHandle(os.Stderr)
	if err != nil {
		return nil, err
	}
	stdout, err := configureOutputHandle(os.Stdout)
	if err != nil {
		return nil, err
	}
	stdin, err := configureInputHandle(os.Stdin)
	if err != nil {
		return nil, err
	}

	return &Streams{
		Stdout: stdout,
		Stderr: stderr,
		Stdin:  stdin,
	}, nil
}

package aby

import (
	"log"
	"os/exec"
)

type ABY struct {
	*exec.Cmd
}

func NewABY() *ABY {
	return &ABY{
		Cmd: exec.Command(""),
	}
}

func (aby *ABY) Server(datas, params []string) error {
	aby.Args = append(aby.Args, "-r", "0", "-a", "0.0.0.0")
	aby.Args = append(aby.Args, "-A", datas[0], "-B", datas[1], "-C", datas[2])
	aby.Args = append(aby.Args, "-D", params[0], "-E", params[1], "-F", params[2])
	log.Println(aby.Args)
	return aby.Start()
}

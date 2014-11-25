package compute

type Args struct {
	Incoming chan Packet
	Outgoing chan Packet
	//Done chan bool
}

type Computes interface {
	Execute(Args)
}

//type Packet map[string]string
type Packet map[string]interface{}

func Run(computes ...Computes) {
	//done := make()
	in := make(chan Packet, 1000)
	//done := make(chan bool)
	var indx = 1
	for _, compute := range computes {
		out := make(chan Packet, 1000)
		arg := Args{Incoming: in, Outgoing: out}
		for i := 0; i < indx; i++ {
			go compute.Execute(arg)
		}
		in = out
		indx += 1
	}

	for {
		_ = <-in
	}
}

func NewPacket() Packet {
	packet := make(Packet)
	return packet
}

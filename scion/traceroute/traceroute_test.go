package traceroute

import (
	"errors"
	"testing"

	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/pkg/snet"
	snetpath "github.com/scionproto/scion/pkg/snet/path"
)

func TestHandleUDP(t *testing.T){
	replyChan := make(chan reply)
	smpc := scmpHandler{
		replies: replyChan,
	}
	pkg := snet.Packet{
			PacketInfo: snet.PacketInfo{
				Destination: snet.SCIONAddress{
					IA:   addr.MustParseIA("1-ff00:0:110"),
					Host: addr.HostSVC(addr.SvcCS),
				},
				Source: snet.SCIONAddress{
					IA:   addr.MustParseIA("1-ff00:0:112"),
					Host: addr.MustParseHost("127.0.0.1"),
				},
				Path: snetpath.OneHop{},
				Payload: snet.UDPPayload{
					SrcPort: 25,
					DstPort: 1925,
					Payload: []byte("hello packet"),
				},
			},
		}

	_, err := smpc.handle(&pkg)
	if !errors.Is(err, ErrNotSCMPPayload){
		t.Fatalf("error handling packet %v", err)
	}

}
func TestHandleParameterProblem(t *testing.T){
	replyChan := make(chan reply)
	smpc := scmpHandler{
		replies: replyChan,
	}
	pkg := snet.Packet{
			PacketInfo: snet.PacketInfo{
				Destination: snet.SCIONAddress{
					IA:   addr.MustParseIA("1-ff00:0:110"),
					Host: addr.HostSVC(addr.SvcCS),
				},
				Source: snet.SCIONAddress{
					IA:   addr.MustParseIA("1-ff00:0:112"),
					Host: addr.MustParseHost("127.0.0.1"),
				},
				Path: snetpath.OneHop{},
				Payload: snet.SCMPParameterProblem{
					// SrcPort: 25,
					// DstPort: 1925,
					Payload: []byte("hello packet"),
				},
			},
		}

	_, err := smpc.handle(&pkg)
	if !errors.Is(err, ErrParameterProblem){
		t.Fatalf("error handling packet %v", err)
	}
}

func TestHandleOtherProblem(t *testing.T){
	replyChan := make(chan reply)
	smpc := scmpHandler{
		replies: replyChan,
	}
	pkg := snet.Packet{
			PacketInfo: snet.PacketInfo{
				Destination: snet.SCIONAddress{
					IA:   addr.MustParseIA("1-ff00:0:110"),
					Host: addr.HostSVC(addr.SvcCS),
				},
				Source: snet.SCIONAddress{
					IA:   addr.MustParseIA("1-ff00:0:112"),
					Host: addr.MustParseHost("127.0.0.1"),
				},
				Path: snetpath.OneHop{},
				Payload: snet.SCMPPacketTooBig{
					Payload: []byte("hello packet"),
				},
			},
		}

	_, err := smpc.handle(&pkg)
	if !errors.Is(err, ErrWrongPayloadType){
		t.Fatalf("error handling packet %v", err)
	}
}

func TestHandleTracerouteReply(t *testing.T){
	replyChan := make(chan reply)
	smpc := scmpHandler{
		replies: replyChan,
	}
	pkg := snet.Packet{
			PacketInfo: snet.PacketInfo{
				Destination: snet.SCIONAddress{
					IA:   addr.MustParseIA("1-ff00:0:110"),
					Host: addr.HostSVC(addr.SvcCS),
				},
				Source: snet.SCIONAddress{
					IA:   addr.MustParseIA("1-ff00:0:112"),
					Host: addr.MustParseHost("127.0.0.1"),
				},
				Path: snetpath.OneHop{},
				Payload: snet.SCMPTracerouteReply{
				},
			},
		}

	_, err := smpc.handle(&pkg)
	if err != nil {
		t.Fatalf("should have handled package succesfully %v", err)
	}
}

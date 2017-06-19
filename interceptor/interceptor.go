package interceptor

import (
	"fmt"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/op/go-logging"

	"github.com/grpc-ecosystem/go-grpc-middleware"

)

/*

var headers = []string{
    "x-ot-span-context",
    "x-request-id",
    "x-b3-traceid",
    "x-b3-spanid",
    "x-b3-parentspanid",
    "x-b3-sampled",
    "x-b3-flags",
  }

var hm = map[string]string{}
*/

var logger *logging.Logger // package-level logger

func init() {
	logger = logging.MustGetLogger("interceptor/interceptor")
}



func BlockUnaryServerInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
  ) (interface{}, error) {

	// validate 'authorization' metadata
	// like headers, the value is an slice []string
fmt.Printf("\nHOLA SERVER\n")


md, ok := metadata.FromContext(ctx)
if !ok {
	fmt.Printf("Server empty, no metadata in request context. \n")
	return handler(ctx, req)
}

	GRPCRecieved(md)

	resp, err := handler(ctx, req)
		if err != nil {
			fmt.Printf("Returning from %s, error: %s", info.FullMethod, err.Error())
		} else {
			fmt.Printf("Returning from %s, response: %s", info.FullMethod, resp)
		}
		grpc.SetHeader(ctx, metadata.Pairs(GRPCMetadata()...))
		return resp, err

	//getIDs(ctx)
	//ctx = setIDs(ctx)

   // handle scopes?
   // ...
   //
}

func BlockStreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {


	fmt.Printf("\nHOLA STREAM SERVER\n")

	stream := grpc_middleware.WrapServerStream(ss)
	ctx := stream.Context()

	md, ok := metadata.FromContext(ctx)
	if !ok {
		fmt.Printf("StreamServer  empty \n")
		return handler(srv, ss)
	}

	GRPCRecieved(md)


	err :=  handler(srv, stream)
	if err != nil {
		fmt.Printf("Returning from %s, error: %s", info.FullMethod, err.Error())
	} else {
		fmt.Printf("Returning from %s, response stream", info.FullMethod)
	}
	stream.SetHeader(metadata.Pairs(GRPCMetadata()...))

	//getIDs(stream.Context())
	//ctx := setIDs(stream.Context())

	stream.WrappedContext = ctx

	return err
}

func BlockUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	fmt.Printf("\nHOLA CLIENT\n")

/*

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for i, n := range md {
						fmt.Printf("Client  %s: %s\n", i, n)
				}
		md = md.Copy()
	} else {
		fmt.Printf("Client  empty \n")
	}

  ctx = metadata.NewOutgoingContext(ctx, md)


  err := invoker(ctx, method, req, reply, cc, opts...)
  	if err != nil {
  	}
  return err

*/


var md metadata.MD
	err := invoker(metadata.NewContext(ctx, metadata.Pairs(GRPCMetadata()...)), method, req, reply, cc, append(opts, grpc.Header(&md))...)
	GRPCReturned(md)
return err

}

func BlockStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
/*

	//ctx = NewOutgoingContext(ctx)

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for i, n := range md {
						fmt.Printf("Stream Client  %s: %s\n", i, n)
				}
		md = md.Copy()
	} else {
		fmt.Printf("Stream Client  empty \n")
	}

  ctx = metadata.NewOutgoingContext(ctx, md)


/*
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for i, n := range md {
						fmt.Printf("Client  %s: %s\n", i, n)
				}
		md = md.Copy()
	} else {
		fmt.Printf("Client  empty \n")
	}

	ctx = metadata.NewOutgoingContext(ctx, md)
*/

fmt.Printf("\nHOLA STREAM CLIENT\n")

//var header metadata.MD

	clientStream, err := streamer(metadata.NewContext(ctx, metadata.Pairs(GRPCMetadata()...)), desc, cc, method, opts...) //append(opts, grpc.Header(&md))...)
		if err != nil {
			return nil, err
	}


	//header, err = clientStream.Header()
	//

	header, erro := stream.Header()
	if erro != nil {
		fmt.Printf("Failed to get header from client: %v\n", erro)
		return clientStream, err
	} else {
		GRPCReturned(header)
	}

	return clientStream, err

}

/*
// NewOutgoingContext creates a new outgoing context with metadata options.
// By default it copies all the incoming metadata from the input context.
// This should only be used in Client interceptors.
func NewOutgoingContext(ctx context.Context ) context.Context {

		md, ok := metadata.FromOutgoingContext(ctx)
		if ok {

			for i, n := range md {
			        fmt.Printf("Client  %s: %s\n", i, n)
			    }


	    for i := 0; i < len(headers); i++ {
	      if id := getID(md, headers[i]); len(id) > 0 {
					fmt.Printf( "client key %s bucket %s\n", headers[i], id)
	        hm[headers[i]] = id
	      }
	    }
				md = metadata.New(hm)
		} else {
			md = metadata.MD{}
		}
	//for _, opt := range opts {
	//	opt(md)
	//}

	return metadata.NewOutgoingContext(ctx, md)
}



// setIDs will set the trace ids on the context{
func setIDs(ctx context.Context) context.Context {

  md := metadata.New(hm)

	if existing, ok := metadata.FromContext(ctx); ok {
		md = metadata.Join(existing, md)
	}

	return metadata.NewContext(ctx, md)
}

// getIDs will return ids embededd an ahe context.
func getIDs(ctx context.Context) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {

		for i, n := range md {
		        fmt.Printf("OLD %s: %s\n", i, n)
		    }


    for i := 0; i < len(headers); i++ {
      if id := getID(md, headers[i]); len(id) > 0 {
				fmt.Printf("key %s bucket %s\n", headers[i], id)
        hm[headers[i]] = id
      }
    }
	}
}

// getID parses an id from the metadata.
func getID(md metadata.MD, name string) string {
	for _, str := range md[name] {

		if len(str) > 0 {
			return str
		}
	}
	return ""
}

*/

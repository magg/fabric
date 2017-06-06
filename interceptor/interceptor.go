package interceptor

import (
	"strconv"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/op/go-logging"

)



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
	//getIDs(ctx)

   // handle scopes?
   // ...
   return handler(ctx, req)
}

func BlockUnaryClientInterceptor(
  ctx context.Context,
  method string,
  req,
  reply interface{},
  cc *grpc.ClientConn,
  invoker grpc.UnaryInvoker,
  opts ...grpc.CallOption,
  ) (error) {

  ctx = NewOutgoingContext(ctx)

  err := invoker(ctx, method, req, reply, cc, opts...)
  	if err != nil {
  	}
  return err

}

// NewOutgoingContext creates a new outgoing context with metadata options.
// By default it copies all the incoming metadata from the input context.
// This should only be used in Client interceptors.
func NewOutgoingContext(ctx context.Context ) context.Context {
	//opts ...MetadataOption) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md = md.Copy()
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
	if md, ok := metadata.FromContext(ctx); ok {
    for i := 0; i < len(headers); i++ {
      if id := getID(md, headers[i]); id > 0 {
				logger.Errorf("Replica %s received an unknown message type %s", headers[i], strconv.Itoa(id))
        hm[headers[i]] = strconv.Itoa(id)
      }
    }
	}
}

// getID parses an id from the metadata.
func getID(md metadata.MD, name string) int {
	for _, str := range md[name] {
		id, err := strconv.Atoi(str)
		if err == nil {
			return id
		}
	}
	return 0
}

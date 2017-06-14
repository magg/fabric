package interceptor



import (
	"github.com/hyperledger/magg/client"
)


  const REQUEST_KEY = "x-request-id"

  const TRACE_KEY = "x-b3-traceid"



  // Returns a slice of strings suitable for passing to grpc/metadata.Pairs
  func GRPCMetadata() []string {
  	r := client.GetRPCMetadata()
  	return []string{REQUEST_KEY,r.RequestID, TRACE_KEY, r.TraceID}
  }


func GRPCRecieved(md map[string][]string) {
	req_strs, ok := md[REQUEST_KEY]
	if !ok || len(req_strs) < 1 {
		return
	}
	req := getIDs(md, REQUEST_KEY)

	trace_strs, ok := md[TRACE_KEY]
	if !ok || len(trace_strs) < 1 {
		return
	}
	traceID= getIDs(md, TRACE_KEY)


	client.RPCReceived(client.RPCMetadata{
		RequestID: req,
		TraceID: traceID,
	})
}

func GRPCReturned(md map[string][]string) {
  req_strs, ok := md[REQUEST_KEY]
  if !ok || len(req_strs) < 1 {
    return
  }
  req := getIDs(md, REQUEST_KEY)

  trace_strs, ok := md[TRACE_KEY]
  if !ok || len(trace_strs) < 1 {
    return
  }
  traceID= getIDs(md, TRACE_KEY)

	client.RPCReturned(client.RPCMetadata{
    RequestID: req,
		TraceID: traceID
	})
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
